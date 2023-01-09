package s3

import (
	"github.com/magicLian/gostarter/pkg/setting"
)

type S3Service interface {
	ExistObject(objectPath, fileName string) (bool, error)
	UploadFile(objectPath string, fileName string, filContent []byte) error
	DownloadFile(objectPath string, fileName string) ([]byte, error)
	DeleteS3Object(objName string) error
	MoveObject(sourceObjName, targetObjName string) error
	GetContentTypeByFileName(fileName string) string
}

func ProvideS3Service(cfg *setting.Cfg) (S3Service, error) {
	return NewS3Service(cfg)
}
