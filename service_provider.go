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

	app.BindWith(Binding, func(app foundation.Application, parameters map[string]any) (any, error) {
		return NewS3(context.Background(), app.MakeConfig(), parameters["disk"].(string))
	})
}

func (receiver *ServiceProvider) Boot(app foundation.Application) {
	app.Publishes("github.com/goravel/s3", map[string]string{
		"config/s3.go": app.ConfigPath(""),
	})
}
