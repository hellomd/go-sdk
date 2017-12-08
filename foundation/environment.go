package foundation

import (
	"github.com/urfave/negroni"
)

const (
	prodEnv = "production"
	stagEnv = "staging"
)

// Environment -
type Environment struct {
	IsDevelopment bool
	Pipeline      *negroni.Negroni
}

// NewEnv -
func NewEnv() *Environment {
	return newProdEnv()
}
