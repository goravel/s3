package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("s3", map[string]any{
		"key":    config.Env("AWS_ACCESS_KEY_ID"),
		"secret": config.Env("AWS_ACCESS_KEY_SECRET"),
		"region": config.Env("AWS_REGION"),
		"bucket": config.Env("AWS_BUCKET"),
		"url":    config.Env("AWS_URL"),
	})
}
