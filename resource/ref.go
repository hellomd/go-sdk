package resource

import "strings"

// Ref stores a reference to a resource
type Ref string

// Type returns the type part of the resource reference
func (r Ref) Type() string {
	return r.parts()[0]
}

// ID returns the ID part of the resource reference
func (r Ref) ID() string {
	return r.parts()[1]
}

func (r Ref) parts() []string {
	return strings.SplitN(string(r), ":", 2)
}
