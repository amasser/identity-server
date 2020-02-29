package iam

import "strings"

// UserURN uniquely identifies a paticular user/account
type UserURN string

// IsValid returns true if the URN is a valid IAM User URN. False
// otherwise.
func (urn UserURN) IsValid() bool {
	return strings.HasPrefix(string(urn), "urn:iam::user/")
}

// Path returns the last part of the URN. For UserURN, that is
// user/<account-id>
func (urn UserURN) Path() string {
	if !urn.IsValid() {
		return ""
	}

	parts := strings.Split(string(urn), ":")
	path := parts[3]
	return path
}

// SubType returns the subtype of the resource
func (urn UserURN) SubType() string {
	path := urn.Path()
	if path == "" {
		return ""
	}

	parts := strings.Split(path, "/")
	return parts[0]
}

// AccountID returns the account ID for the user.
func (urn UserURN) AccountID() string {
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

// User represents a user account managed by IAM and authenticated by
// authn-server
type User struct {
	AccountID  int                    `json:"accountID"`
	Username   string                 `json:"username"`
	ID         UserURN                `json:"id"`
	Locked     *bool                  `json:"locked,omitempty"`
	Attributes map[string]interface{} `json:"attrs,omitempty"`
}
