package util

import (
	"fmt"
	"strings"
)

func GetFileTypeByFileName(filename string) string {
	dotIndex := strings.LastIndex(filename, ".")
	return filename[dotIndex:]
}

func GetTempFileName(username, filename string) string {
	return fmt.Sprintf("%s-%s", username, filename)
}

func GetFileNameByTempFile(username, tempFileName string) string {
	if strings.HasPrefix(tempFileName, username) {
		return strings.Replace(tempFileName, username+"-", "", 1)
	}
	return tempFileName
}
