package client

import (
	"context"
	"errors"

	"github.com/tierklinik-dobersberg/identity-server/pkg/iam"
)

// GroupClient implements a HTTP client for the group management
// endpoints.
type GroupClient struct {
	*IdentityClient
}

// Get returns all groups managed by IAM.
func (gc *GroupClient) Get(ctx context.Context) ([]iam.Group, error) {
	req, err := gc.newRequest(ctx, "GET", "/v1/groups/", nil)
	if err != nil {
		return nil, err
	}

	res, err := gc.cli.Do(req)
	if err != nil {
		return nil, err
	}

	var list struct {
		Groups []iam.Group `json:"groups"`
	}

	if err := gc.parseResponse(res, &list); err != nil {
		return nil, err
	}

	return list.Groups, nil
}

// Create creates a new group.
func (gc *GroupClient) Create(ctx context.Context, groupName string, groupComment string) (iam.GroupURN, error) {
	body := struct {
		Name    string `json:"name"`
		Comment string `json:"comment"`
	}{
		Name:    groupName,
		Comment: groupComment,
	}

	req, err := gc.newRequest(ctx, "POST", "/v1/groups/", body)
	if err != nil {
		return "", err
	}

	res, err := gc.cli.Do(req)
	if err != nil {
		return "", err
	}

	var response struct {
		URN iam.GroupURN `json:"urn"`
	}
	if err := gc.parseResponse(res, &response); err != nil {
		return "", nil
	}
	return response.URN, nil
}

// Delete deletes an exsting group.
func (gc *GroupClient) Delete(ctx context.Context, urn iam.GroupURN) error {
	name := urn.GroupName()
	if name == "" {
		return errors.New("Invalid group name")
	}

	req, err := gc.newRequest(ctx, "DELETE", "/v1/groups/"+name, nil)
	if err != nil {
		return err
	}

	res, err := gc.cli.Do(req)
	if err != nil {
		return err
	}

	return gc.parseResponse(res, nil)
}

// Load loads an exsiting group by it's URN.
func (gc *GroupClient) Load(ctx context.Context, urn iam.GroupURN) (iam.Group, error) {
	var u iam.Group

	name := urn.GroupName()
	if name == "" {
		return u, errors.New("Invalid group name")
	}

	req, err := gc.newRequest(ctx, "GET", "/v1/groups/"+name, nil)
	if err != nil {
		return u, err
	}

	res, err := gc.cli.Do(req)
	if err != nil {
		return u, err
	}

	return u, gc.parseResponse(res, &u)
}

// UpdateComment updates the comment of a group.
func (gc *GroupClient) UpdateComment(ctx context.Context, urn iam.GroupURN, comment string) error {
	name := urn.GroupName()
	if name == "" {
		return errors.New("Invalid group name")
	}

	body := struct {
		Comment string `json:"comment"`
	}{comment}

	req, err := gc.newRequest(ctx, "PUT", "/v1/groups/"+name, body)
	if err != nil {
		return err
	}

	res, err := gc.cli.Do(req)
	if err != nil {
		return err
	}

	return gc.parseResponse(res, nil)
}

// AddMember adds a new member to a group.
func (gc *GroupClient) AddMember(ctx context.Context, grp iam.GroupURN, member iam.UserURN) error {
	name := grp.GroupName()
	if name == "" {
		return errors.New("Invalid group name")
	}

	user := member.AccountID()
	if user == "" {
		return errors.New("Invalid user id")
	}

	req, err := gc.newRequest(ctx, "PUT", "/v1/groups/"+name+"/members/"+user, nil)
	if err != nil {
		return err
	}

	res, err := gc.cli.Do(req)
	if err != nil {
		return err
	}

	return gc.parseResponse(res, nil)
}

// DeleteMember deletes a member from a group.
func (gc *GroupClient) DeleteMember(ctx context.Context, grp iam.GroupURN, member iam.UserURN) error {
	name := grp.GroupName()
	if name == "" {
		return errors.New("Invalid group name")
	}

	user := member.AccountID()
	if user == "" {
		return errors.New("Invalid user id")
	}

	req, err := gc.newRequest(ctx, "DELETE", "/v1/groups/"+name+"/members/"+user, nil)
	if err != nil {
		return err
	}

	res, err := gc.cli.Do(req)
	if err != nil {
		return err
	}

	return gc.parseResponse(res, nil)
}
