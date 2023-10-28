package s3

import (
	"fmt"
	"path"
	"strings"

	"github.com/goravel/framework/contracts/filesystem"
	"github.com/goravel/framework/support/file"
)

func fullPathOfFile(filePath string, source filesystem.File, name string) (string, error) {
	extension := path.Ext(name)
	if extension == "" {
		var err error
		extension, err = file.Extension(source.File(), true)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s/%s.%s", filePath, strings.TrimSuffix(strings.TrimPrefix(path.Base(name), "/"), "/"), extension), nil
	} else {
		return fmt.Sprintf("%s/%s", filePath, strings.TrimPrefix(path.Base(name), "/")), nil
	}
}

func validPath(path string) string {
	realPath := strings.TrimPrefix(path, "./")
	realPath = strings.TrimPrefix(realPath, "/")
	realPath = strings.TrimPrefix(realPath, ".")
	if realPath != "" && !strings.HasSuffix(realPath, "/") {
		realPath += "/"
	}

	return realPath
}
