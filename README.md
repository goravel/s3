[utils_test.go](..%2Foss%2Futils_test.go)# S3

A s3 disk driver for `facades.Storage()` of Goravel.

## Version

| goravel/s3 | goravel/framework |
|------------|-------------------|
| v1.4.*     | v1.16.*           |
| v1.3.*     | v1.15.*           |
| v1.2.*     | v1.14.*           |
| v1.1.*     | v1.13.*           |
| v1.0.*     | v1.12.*           |

## Install

Run the command below in your project to install the package automatically:

```
./artisan package:install github.com/goravel/s3
```

Or check [the setup file](./setup/setup.go) to install the package manually.

## Testing

Run command below to run test:

```
AWS_ACCESS_KEY_ID= AWS_ACCESS_KEY_SECRET= AWS_REGION= AWS_BUCKET= AWS_URL= go test ./...
```
