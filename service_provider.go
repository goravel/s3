package s3

import (
	"context"

	"github.com/goravel/framework/contracts/foundation"
)

const Binding = "goravel.s3"

var App foundation.Application

type ServiceProvider struct {
}

func (receiver *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.Bind(Binding, func(app foundation.Application) (any, error) {
		return NewS3(context.Background(), app.MakeConfig())
	})
}

func (receiver *ServiceProvider) Boot(app foundation.Application) {
	app.Publishes("github.com/goravel/s3", map[string]string{
		"config/s3.go": app.ConfigPath(""),
	})
}
