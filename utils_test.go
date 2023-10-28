package s3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullPathOfFile(t *testing.T) {
	path, err := fullPathOfFile("a/b", &File{path: "logo.png"}, "c")
	assert.Equal(t, "a/b/c.png", path)
	assert.Nil(t, err)

	path, err = fullPathOfFile("a/b", &File{path: "logo.png"}, "c.jpg")
	assert.Equal(t, "a/b/c.jpg", path)
	assert.Nil(t, err)
}

func TestValidPath(t *testing.T) {
	assert.Equal(t, "a/", validPath("./a"))
	assert.Equal(t, "a/", validPath(".a"))
	assert.Equal(t, "a/", validPath("/a"))
	assert.Equal(t, "a/", validPath("./a/"))
	assert.Equal(t, "a/", validPath(".a/"))
	assert.Equal(t, "a/", validPath("/a/"))
	assert.Equal(t, "", validPath("./"))
	assert.Equal(t, "", validPath("."))
	assert.Equal(t, "", validPath("/"))
}
