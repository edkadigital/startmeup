package ui

import (
	"fmt"
	"testing"

	"github.com/edkadigital/startmeup/config"
	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	path := "abc.txt"
	got := File(path)
	expected := fmt.Sprintf("/%s/%s?v=%s", config.StaticPrefix, path, cacheBuster)
	assert.Equal(t, expected, got)
}
