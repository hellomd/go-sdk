package foundation

import (
	"github.com/hellomd/go-sdk/config"
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
	switch config.Get(EnvCfgKey) {
	case prodEnv, stagEnv:
		return newProdEnv()
	default:
		return newDevEnv()
	}
}
