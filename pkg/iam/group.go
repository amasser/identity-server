package iam

import "strings"

// GroupURN describes the unique resource name of a user/account group
type GroupURN string

// IsValid returns true if the URN is a valid IAM Group URN. False
// otherwise.
func (urn GroupURN) IsValid() bool {
	return strings.HasPrefix(string(urn), "urn:iam::group/")
}

// Path returns the last part of the URN. For GroupURN, that is
// group/<group-name>
func (urn GroupURN) Path() string {
	if !urn.IsValid() {
		return ""
	}

	parts := strings.Split(string(urn), ":")
	path := parts[3]
	return path
}

// SubType returns the subtype of the resource.
func (urn GroupURN) SubType() string {
	path := urn.Path()
	if path == "" {
		return ""
	}

	parts := strings.Split(path, "/")
	return parts[0]
}

// GroupName returns the name of the group.
func (urn GroupURN) GroupName() string {
	path := urn.Path()
	if path == "" {
		return ""
	}

	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// Group is a collection of user accounts (iam.User) and can be used for
// authorization
type Group struct {
	ID      GroupURN `json:"id"`
	Name    string   `json:"name"`
	Comment string   `json:"comment,omitempty"`
}
