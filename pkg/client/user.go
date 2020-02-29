package client

import (
	"context"
	"errors"

	"github.com/tierklinik-dobersberg/identity-server/pkg/iam"
)

// UserClient provides access to user management endpoints.
type UserClient struct {
	*IdentityClient
}

// CreateUser creates a new user.
func (uc *UserClient) CreateUser(ctx context.Context, username, password string, attrs map[string]interface{}) (iam.UserURN, error) {
	body := struct {
		Username   string                 `json:"username"`
		Password   string                 `json:"password"`
		Attributes map[string]interface{} `json:"attrs"`
	}{
		Username:   username,
		Password:   password,
		Attributes: attrs,
	}

	req, err := uc.newRequest(ctx, "POST", "/v1/users/", body)
	if err != nil {
		return "", err
	}

	res, err := uc.cli.Do(req)
	if err != nil {
		return "", err
	}

	response := struct {
		URN iam.UserURN `json:"urn"`
	}{}
	if err := uc.parseResponse(res, &response); err != nil {
		return "", err
	}

	return response.URN, nil
}

// LoadUser loads the user identified by URN.
func (uc *UserClient) LoadUser(ctx context.Context, urn iam.UserURN) (iam.User, error) {
	id := urn.AccountID()
	if id == "" {
		return iam.User{}, errors.New("Invalid UserURN")
	}

	req, err := uc.newRequest(ctx, "GET", "/v1/users/"+id, nil)
	if err != nil {
		return iam.User{}, nil
	}

	res, err := uc.cli.Do(req)
	if err != nil {
		return iam.User{}, err
	}

	var u iam.User
	if err := uc.parseResponse(res, &u); err != nil {
		return iam.User{}, err
	}

	return u, nil
}

// DeleteUser deletes the user identified by URN.
func (uc *UserClient) DeleteUser(ctx context.Context, urn iam.UserURN) error {
	id := urn.AccountID()
	if id == "" {
		return errors.New("Invalid UserURN")
	}

	req, err := uc.newRequest(ctx, "DELETE", "/v1/users/"+id, nil)
	if err != nil {
		return err
	}

	res, err := uc.cli.Do(req)
	if err != nil {
		return err
	}

	return uc.parseResponse(res, nil)
}

// LockUser locks or unlocks the user identified by URN.
func (uc *UserClient) LockUser(ctx context.Context, urn iam.UserURN, locked bool) error {
	id := urn.AccountID()
	if id == "" {
		return errors.New("Invalid UserURN")
	}

	method := "PUT"
	if !locked {
		method = "DELETE"
	}

	req, err := uc.newRequest(ctx, method, "/v1/users/"+id+"/locked", nil)
	if err != nil {
		return err
	}

	res, err := uc.cli.Do(req)
	if err != nil {
		return err
	}

	return uc.parseResponse(res, nil)
}

// Users returns a list of users stored and managed by IAM.
func (uc *UserClient) Users(ctx context.Context) ([]iam.User, error) {
	req, err := uc.newRequest(ctx, "GET", "/v1/users/", nil)
	if err != nil {
		return nil, err
	}

	res, err := uc.cli.Do(req)
	if err != nil {
		return nil, err
	}

	var response struct {
		Users []iam.User `json:"users"`
	}

	if err := uc.parseResponse(res, &response); err != nil {
		return nil, err
	}

	return response.Users, nil
}

// UpdateAttrs updates a users attributes.
func (uc *UserClient) UpdateAttrs(ctx context.Context, urn iam.UserURN, attrs map[string]interface{}) error {
	id := urn.AccountID()
	if id == "" {
		return errors.New("Invalid UserURN")
	}

	req, err := uc.newRequest(ctx, "PUT", "/v1/users/"+id+"/attrs/", attrs)
	if err != nil {
		return err
	}

	res, err := uc.cli.Do(req)
	if err != nil {
		return err
	}

	return uc.parseResponse(res, nil)
}

// SetAttr updates a single user attribute.
func (uc *UserClient) SetAttr(ctx context.Context, urn iam.UserURN, key string, value interface{}) error {
	id := urn.AccountID()
	if id == "" {
		return errors.New("Invalid UserURN")
	}

	req, err := uc.newRequest(ctx, "PUT", "/v1/users/"+id+"/attrs/"+key, value)
	if err != nil {
		return err
	}

	res, err := uc.cli.Do(req)
	if err != nil {
		return err
	}

	return uc.parseResponse(res, nil)
}

// DeleteAttr updates a single user attribute.
func (uc *UserClient) DeleteAttr(ctx context.Context, urn iam.UserURN, key string) error {
	id := urn.AccountID()
	if id == "" {
		return errors.New("Invalid UserURN")
	}

	req, err := uc.newRequest(ctx, "DELETE", "/v1/users/"+id+"/attrs/"+key, nil)
	if err != nil {
		return err
	}

	res, err := uc.cli.Do(req)
	if err != nil {
		return err
	}

	return uc.parseResponse(res, nil)
}
