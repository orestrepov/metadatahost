package app

import (
	"github.com/orestrepov/metadatahost/db"
	"github.com/sirupsen/logrus"
)

type Context struct {
	Logger        logrus.FieldLogger
	RemoteAddress string
	Database      *db.Database
	Config        *Config
}

func (ctx *Context) WithLogger(logger logrus.FieldLogger) *Context {
	ret := *ctx
	ret.Logger = logger
	return &ret
}

func (ctx *Context) WithRemoteAddress(address string) *Context {
	ret := *ctx
	ret.RemoteAddress = address
	return &ret
}
