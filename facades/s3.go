package facades

import (
	"github.com/goravel/framework/contracts/filesystem"

	"github.com/goravel/s3"
)

func S3(disk string) (filesystem.Driver, error) {
	instance, err := s3.App.MakeWith(s3.Binding, map[string]any{"disk": disk})
	if err != nil {
		return nil, err
	}

	return instance.(*s3.S3), nil
}
