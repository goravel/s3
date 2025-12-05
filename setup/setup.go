package main

import (
	"os"

	"github.com/goravel/framework/packages"
	"github.com/goravel/framework/packages/match"
	"github.com/goravel/framework/packages/modify"
	"github.com/goravel/framework/support/env"
	"github.com/goravel/framework/support/path"
)

func main() {
	config := `map[string]any{
        "driver": "custom",
        "key": config.Env("AWS_ACCESS_KEY_ID"),
        "secret": config.Env("AWS_ACCESS_KEY_SECRET"),
        "region": config.Env("AWS_REGION"),
        "bucket": config.Env("AWS_BUCKET"),
        "url": config.Env("AWS_URL"),
        "via": func() (filesystem.Driver, error) {
            return s3facades.S3("s3") // The ` + "`s3`" + ` value is the ` + "`disks`" + ` key
        },
    }`

	appConfigPath := path.Config("app.go")
	filesystemsConfigPath := path.Config("filesystems.go")
	modulePath := packages.GetModulePath()
	s3ServiceProvider := "&s3.ServiceProvider{}"
	filesystemContract := "github.com/goravel/framework/contracts/filesystem"
	s3Facades := "github.com/goravel/s3/facades"
	filesystemsDisksConfig := match.Config("filesystems.disks")
	filesystemsConfig := match.Config("filesystems")

	packages.Setup(os.Args).
		Install(
			// Add s3 service provider to app.go if not using bootstrap setup
			modify.When(func(_ map[string]any) bool {
				return !env.IsBootstrapSetup()
			}, modify.GoFile(appConfigPath).
				Find(match.Imports()).Modify(modify.AddImport(modulePath)).
				Find(match.Providers()).Modify(modify.Register(s3ServiceProvider))),

			// Add s3 service provider to providers.go if using bootstrap setup
			modify.When(func(_ map[string]any) bool {
				return env.IsBootstrapSetup()
			}, modify.AddProviderApply(modulePath, s3ServiceProvider)),

			// Add s3 disk to filesystems.go
			modify.GoFile(filesystemsConfigPath).Find(match.Imports()).Modify(
				modify.AddImport(filesystemContract),
				modify.AddImport(s3Facades, "s3facades"),
			).
				Find(filesystemsDisksConfig).Modify(modify.AddConfig("s3", config)).
				Find(filesystemsConfig).Modify(modify.AddConfig("default", `"s3"`)),
		).
		Uninstall(
			// Remove s3 disk from filesystems.go
			modify.GoFile(filesystemsConfigPath).
				Find(filesystemsConfig).Modify(modify.AddConfig("default", `"local"`)).
				Find(filesystemsDisksConfig).Modify(modify.RemoveConfig("s3")).
				Find(match.Imports()).Modify(
				modify.RemoveImport(filesystemContract),
				modify.RemoveImport(s3Facades, "s3facades"),
			),

			// Remove s3 service provider from app.go if not using bootstrap setup
			modify.When(func(_ map[string]any) bool {
				return !env.IsBootstrapSetup()
			}, modify.GoFile(appConfigPath).
				Find(match.Providers()).Modify(modify.Unregister(s3ServiceProvider)).
				Find(match.Imports()).Modify(modify.RemoveImport(modulePath))),

			// Remove s3 service provider from providers.go if using bootstrap setup
			modify.When(func(_ map[string]any) bool {
				return env.IsBootstrapSetup()
			}, modify.RemoveProviderApply(modulePath, s3ServiceProvider)),
		).
		Execute()
}
