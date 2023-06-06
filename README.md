# s3

A s3 disk driver for facades.Storage of Goravel.

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
    &sms.ServiceProvider{},
}
```

2. Publish configuration file

```
go run . artisan vendor:publish --package=github.com/goravel/s3
```

3. Fill your s3 configuration to `config/s3.go` file

4. Add s3 disk to `config/filesystems.go` file

## Testing

Run command below to run test(fill your owner s3 configuration):

```
AWS_ACCESS_KEY_ID= AWS_ACCESS_KEY_SECRET= AWS_DEFAULT_REGION= AWS_BUCKET= AWS_URL= go test ./...
```
