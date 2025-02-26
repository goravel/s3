[utils_test.go](..%2Foss%2Futils_test.go)# S3

A s3 disk driver for `facades.Storage()` of Goravel.

## Version

| goravel/s3 | goravel/framework |
| ---------- | ----------------- |
| v1.3.\*    | v1.15.\*          |
| v1.2.\*    | v1.14.\*          |
| v1.1.\*    | v1.13.\*          |
| v1.0.\*    | v1.12.\*          |

## Install

1. Add package

```
go get -u github.com/goravel/s3
```

2. Register service provider

```
// config/app.go
import "github.com/goravel/s3"

"providers": []foundation.ServiceProvider{
...
&s3.ServiceProvider{},
}
```

3. Add s3 disk to `config/filesystems.go` file

1. AWS Configuration

```
// config/filesystems.go
import (
"github.com/goravel/framework/contracts/filesystem"
s3facades "github.com/goravel/s3/facades"
)

"disks": map[string]any{
    ...
    "s3": map[string]any{
        "driver": "custom",
        "key": config.Env("S3_ACCESS_KEY_ID"),
        "secret": config.Env("S3_ACCESS_KEY_SECRET"),
        "region": config.Env("S3_REGION"),
        "bucket": config.Env("S3_BUCKET"),
        "url": config.Env("S3_URL"),
        "via": func() (filesystem.Driver, error) {
            return s3facades.S3("s3"), nil // The `s3` value is the `disks` key
        },
    },
}
```

2. DigitalOcean

```
// config/filesystems.go
import (
"github.com/goravel/framework/contracts/filesystem"
s3facades "github.com/goravel/s3/facades"
)

"disks": map[string]any{
    ...
    "s3": map[string]any{
        "driver": "custom",
        "key": config.Env("S3_ACCESS_KEY_ID"),
        "secret": config.Env("S3_ACCESS_KEY_SECRET"),
        "region": config.Env("S3_REGION"),
        "bucket": config.Env("S3_BUCKET"),
        "url": config.Env("S3_URL"),
        "endpoint": config.Env("S3_ENDPOINT", "https://sfo3.digitaloceanspaces.com"),
        "use_path_style": config.Env("S3_USE_PATH_STYLE", true),
        "via": func() (filesystem.Driver, error) {
            return s3facades.S3("s3"), nil // The `s3` value is the `disks` key
        },
    },
}
```

## Testing

Run command below to run test(fill your owner s3 configuration):

```
S3_ACCESS_KEY_ID= S3_ACCESS_KEY_SECRET= S3_DEFAULT_REGION= S3_BUCKET= S3_URL= go test ./...
```
