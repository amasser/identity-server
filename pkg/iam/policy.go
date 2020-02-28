package iam

import "github.com/ory/ladon"

// PolicyURN describes an access and permission policy managed
// by IAM.
type PolicyURN string

// Policy wraps ladon.DefaultPolicy
type Policy struct {
	ladon.DefaultPolicy
}
