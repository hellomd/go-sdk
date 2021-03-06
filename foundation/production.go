package foundation

import (
	"github.com/hellomd/go-sdk/authentication"
	"github.com/hellomd/go-sdk/blacklist"
	"github.com/hellomd/go-sdk/config"
	"github.com/hellomd/go-sdk/contenttype"
	"github.com/hellomd/go-sdk/errors"
	"github.com/hellomd/go-sdk/logger"
	"github.com/hellomd/go-sdk/logmatic"
	"github.com/hellomd/go-sdk/recovery"
	"github.com/hellomd/go-sdk/requestid"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func newProdEnv() *Environment {
	log := logrus.New()
	log.Hooks.Add(logmatic.NewLogrusHook(logmatic.Config{APIKey: config.Get(LogmaticAPIKeyCfgKey)}))

	pipeline := negroni.New()

	pipeline.Use(blacklist.NewMiddleware())
	pipeline.Use(contenttype.NewMiddleware())
	pipeline.Use(requestid.NewMiddleware())
	pipeline.Use(logger.NewMiddleware(config.Get(AppNameCfgKey), config.Get(EnvCfgKey), log))
	pipeline.Use(errors.NewMiddleware())
	pipeline.Use(recovery.NewMiddleware(config.Get(SentryDSNCfgKey)))
	pipeline.UseFunc(authentication.NewMiddleware([]byte(config.Get(SecretCfgKey))))

	return &Environment{false, pipeline}
}
