package facades

import (
	"log"

	"github.com/goravel/framework/contracts/filesystem"

	"github.com/goravel/s3"
)

func S3(disk string) filesystem.Driver {
	instance, err := s3.App.MakeWith(s3.Binding, map[string]any{"disk": disk})
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return instance.(*s3.S3)
}
