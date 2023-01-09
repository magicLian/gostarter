package util

import (
	"io"

	"github.com/magicLian/gostarter/pkg/models"
)

func GetFormFile(c *models.ReqContext, key string) ([]byte, error) {
	file, err := c.FormFile(key)
	if err != nil {
		return nil, err
	}

	f, err := file.Open()
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
