package iam

// UserURN uniquely identifies a paticular user
type UserURN string

// User represents a user account managed by IAM and authenticated by
// authn-server
type User struct {
	AccountID  int                    `json:"accountID"`
	Username   string                 `json:"username"`
	ID         UserURN                `json:"id"`
	Locked     *bool                  `json:"locked"`
	Groups     []*string              `json:"groups"`
	Attributes map[string]interface{} `json:"attrs,omitempty"`
}
