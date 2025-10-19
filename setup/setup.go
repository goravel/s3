package main

import (
	"os"

	"github.com/goravel/framework/packages"
	"github.com/goravel/framework/packages/match"
	"github.com/goravel/framework/packages/modify"
	"github.com/goravel/framework/support/path"
)

var config = `map[string]any{
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

func main() {
	packages.Setup(os.Args).
		Install(
			modify.GoFile(path.Config("app.go")).
				Find(match.Imports()).Modify(modify.AddImport(packages.GetModulePath())).
				Find(match.Providers()).Modify(modify.Register("&s3.ServiceProvider{}")),
			modify.GoFile(path.Config("filesystems.go")).
				Find(match.Imports()).Modify(modify.AddImport("github.com/goravel/framework/contracts/filesystem"), modify.AddImport("github.com/goravel/s3/facades", "s3facades")).
				Find(match.Config("filesystems.disks")).Modify(modify.AddConfig("s3", config)).
				Find(match.Config("filesystems")).Modify(modify.AddConfig("default", `"s3"`)),
		).
		Uninstall(
			modify.GoFile(path.Config("app.go")).
				Find(match.Providers()).Modify(modify.Unregister("&s3.ServiceProvider{}")).
				Find(match.Imports()).Modify(modify.RemoveImport(packages.GetModulePath())),
			modify.GoFile(path.Config("filesystems.go")).
				Find(match.Config("filesystems.disks")).Modify(modify.RemoveConfig("s3")).
				Find(match.Imports()).Modify(modify.RemoveImport("github.com/goravel/framework/contracts/filesystem"), modify.RemoveImport("github.com/goravel/s3/facades", "s3facades")).
				Find(match.Config("filesystems")).Modify(modify.AddConfig("default", `"local"`)),
		).
		Execute()
}
