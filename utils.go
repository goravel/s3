package s3

import (
	"math/rand"
	"path"
	"strings"

	"github.com/goravel/framework/contracts/filesystem"
	"github.com/goravel/framework/support/file"
	supporttime "github.com/goravel/framework/support/time"
)

func fullPathOfFile(filePath string, source filesystem.File, name string) (string, error) {
	extension := path.Ext(name)
	if extension == "" {
		var err error
		extension, err = file.Extension(source.File(), true)
		if err != nil {
			return "", err
		}

		return strings.TrimSuffix(filePath, "/") + "/" + strings.TrimSuffix(strings.TrimPrefix(path.Base(name), "/"), "/") + "." + extension, nil
	} else {
		return strings.TrimSuffix(filePath, "/") + "/" + strings.TrimPrefix(path.Base(name), "/"), nil
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

func random(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	rand.Seed(supporttime.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}

	return string(result)
}
