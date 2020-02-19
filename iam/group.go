package iam

// GroupURN describes the unique resource name of a user/account group
type GroupURN string

// Group is a collection of user accounts (iam.User) and can be used for
// authorization
type Group struct {
	ID      GroupURN  `json:"id"`
	Name    string    `json:"name"`
	Comment string    `json:"comment,omitempty"`
	Members []UserURN `json:"members,omitempty"`
}
