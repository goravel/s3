[utils_test.go](..%2Foss%2Futils_test.go)# S3

A s3 disk driver for `facades.Storage()` of Goravel.

## Version

| goravel/s3 | goravel/framework |
|------------|-------------------|
| v1.3.*     | v1.15.*           |
| v1.2.*     | v1.14.*           |
| v1.1.*     | v1.13.*           |
| v1.0.*     | v1.12.*           |

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

AWS Configuration

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
        "key": config.Env("AWS_ACCESS_KEY_ID"),
        "secret": config.Env("AWS_ACCESS_KEY_SECRET"),
        "region": config.Env("AWS_REGION"),
        "bucket": config.Env("AWS_BUCKET"),
        "url": config.Env("S3_URL"),
        "via": func() (filesystem.Driver, error) {
            return s3facades.S3("s3") // The `s3` value is the `disks` key
        },
    },
}
```

DigitalOcean

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
        "key": config.Env("SPACES_ACCESS_KEY_ID"),
        "secret": config.Env("SPACES_ACCESS_KEY_SECRET"),
        "region": config.Env("SPACES_REGION", "us-east-1"),
        "bucket": config.Env("SPACES_BUCKET"),
        "url": config.Env("SPACES_URL"),
        "endpoint": config.Env("SPACES_ENDPOINT"),
        "use_path_style": config.Env("SPACES_USE_PATH_STYLE", true),
        "cdn": config.Env("SPACES_CDN"),
        "object_canned_acl: config.Env("SPACES_OBJECT_CANNED_ACL"),
        "via": func() (filesystem.Driver, error) {
            return s3facades.S3("s3"), nil // The `s3` value is the `disks` key
        },
    },
}
```

## Testing

Run command below to run test(fill your owner s3 configuration):

```
AWS_ACCESS_KEY_ID= AWS_ACCESS_KEY_SECRET= AWS_DEFAULT_REGION= AWS_BUCKET= AWS_URL= go test ./...
```
