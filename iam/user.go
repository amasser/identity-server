package iam

// UserURN uniquely identifies a paticular user/account
type UserURN string

// User represents a user account managed by IAM and authenticated by
// authn-server
type User struct {
	AccountID  int                    `json:"accountID"`
	Username   string                 `json:"username"`
	ID         UserURN                `json:"id"`
	Locked     *bool                  `json:"locked,omitempty"`
	Groups     []*string              `json:"groups,omitempty"`
	Attributes map[string]interface{} `json:"attrs,omitempty"`
}
