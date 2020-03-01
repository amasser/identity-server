package iam

import (
	"strings"

	"github.com/ory/ladon"
)

// PolicyURN describes an access and permission policy managed
// by IAM.
type PolicyURN string

// IsValid returns true if the URN is a valid IAM Policy URN. False
// otherwise.
func (urn PolicyURN) IsValid() bool {
	return strings.HasPrefix(string(urn), "urn:iam::policy/")
}

// Path returns the last part of the URN. For PolicyURN, that is
// group/<group-name>
func (urn PolicyURN) Path() string {
	if !urn.IsValid() {
		return ""
	}

	parts := strings.Split(string(urn), ":")
	path := parts[3]
	return path
}

// SubType returns the subtype of the resource.
func (urn PolicyURN) SubType() string {
	path := urn.Path()
	if path == "" {
		return ""
	}

	parts := strings.Split(path, "/")
	return parts[0]
}

// PolicyName returns the name of the group.
func (urn PolicyURN) PolicyName() string {
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

// Policy wraps ladon.DefaultPolicy
type Policy struct {
	ladon.DefaultPolicy
}
