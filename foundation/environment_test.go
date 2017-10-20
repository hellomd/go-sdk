package foundation

import (
	"net/http/httptest"
	"testing"

	"github.com/hellomd/go-sdk/config"
)

func TestDevEnvironment(t *testing.T) {

	config.Set(EnvCfgKey, "")

	env := NewEnv()

	if !env.IsDevelopment {
		t.Error("IsDevelopment expected to be true, got false")
	}

	response := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	env.Pipeline.ServeHTTP(response, req)

	if len(env.Pipeline.Handlers()) != 7 {
		t.Errorf("Unexpected number of middlewares. Expcted %v, got %v", 7, len(env.Pipeline.Handlers()))
	}
}

func TestProdEnvironment(t *testing.T) {

	config.Set(EnvCfgKey, prodEnv)

	env := NewEnv()

	if env.IsDevelopment {
		t.Error("IsDevelopment expected to be false, got true")
	}

	response := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	env.Pipeline.ServeHTTP(response, req)

	if len(env.Pipeline.Handlers()) != 7 {
		t.Errorf("Unexpected number of middlewares. Expcted %v, got %v", 7, len(env.Pipeline.Handlers()))
	}
}

func TestStagingEnvironment(t *testing.T) {

	config.Set(EnvCfgKey, stagEnv)

	env := NewEnv()

	if env.IsDevelopment {
		t.Error("IsDevelopment expected to be false, got true")
	}

	response := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	env.Pipeline.ServeHTTP(response, req)

	if len(env.Pipeline.Handlers()) != 7 {
		t.Errorf("Unexpected number of middlewares. Expcted %v, got %v", 7, len(env.Pipeline.Handlers()))
	}
}
