package facades

import (
	"log"

	"github.com/goravel/framework/contracts/filesystem"

	"github.com/goravel/s3"
)

func S3() filesystem.Driver {
	instance, err := s3.App.Make(s3.Binding)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return instance.(*s3.S3)
}
