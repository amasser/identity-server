package iam

import "github.com/ory/ladon"

// PolicyURN describes an access and permission policy managed
// by IAM.
type PolicyURN string

// Policy wraps ladon.DefaultPolicy
type Policy struct {
	ID PolicyURN `json:"id"`
	ladon.DefaultPolicy
}

func (p *Policy) GetID() string {
	return string(p.ID)
}
