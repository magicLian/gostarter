package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/magicLian/gostarter/pkg/setting"
	"github.com/magicLian/gostarter/pkg/util"

	"github.com/google/uuid"
	"github.com/magicLian/gostarter/pkg/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3ServiceImpl struct {
	config             S3Config
	client             *minio.Client
	uploadBaseFilePath string
	downloadFilePath   string
	log                log.Logger
	fileContentMap     map[string]string
}

type S3Config struct {
	loaction          string
	accessKey         string
	secretKey         string
	endpoint          string
	token             string
	secure            bool
	bucketName        string
	contentTyepMapStr string
}

func NewS3Service(cfg *setting.Cfg) (*S3ServiceImpl, error) {
	s3Svc := &S3ServiceImpl{
		log: log.New("s3"),
	}
	s3Svc.readConfig(cfg)
	s3Svc.initContentMap()
	s3Svc.uploadBaseFilePath = filepath.Join(setting.HomePath, "/data/s3/upload")
	s3Svc.downloadFilePath = filepath.Join(setting.HomePath, "/data/s3/download")
	if err := os.MkdirAll(s3Svc.uploadBaseFilePath, os.ModePerm); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(s3Svc.downloadFilePath, os.ModePerm); err != nil {
		return nil, err
	}
	if err := s3Svc.initClient(); err != nil {
		return nil, err
	}
	return s3Svc, nil
}

func (s3 *S3ServiceImpl) readConfig(cfg *setting.Cfg) {
	config := S3Config{}
	config.loaction = util.SetDefaultString(util.GetEnvOrIniValue(cfg.Raw, "s3", "location"), "us-east-1")
	config.endpoint = util.GetEnvOrIniValue(cfg.Raw, "s3", "endpoint")
	config.secure = util.SetDefaultBool(util.GetEnvOrIniValue(cfg.Raw, "s3", "secure"), false)
	config.accessKey = util.GetEnvOrIniValue(cfg.Raw, "s3", "access_key")
	config.secretKey = util.GetEnvOrIniValue(cfg.Raw, "s3", "secret_key")
	config.token = util.GetEnvOrIniValue(cfg.Raw, "s3", "token")
	config.bucketName = util.GetEnvOrIniValue(cfg.Raw, "s3", "bucket_name")
	config.contentTyepMapStr = util.GetEnvOrIniValue(cfg.Raw, "s3", "content_type_map")
	s3.config = config
}

func (s3 *S3ServiceImpl) initClient() error {
	minioClient, err := minio.New(s3.config.endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3.config.accessKey, s3.config.secretKey, s3.config.token),
		Secure: s3.config.secure,
	})
	if err != nil {
		s3.log.Errorf("Failed to init s3 client", "error", err)
		return err
	}
	exists, errBucketExists := minioClient.BucketExists(context.Background(), s3.config.bucketName)
	if errBucketExists != nil {
		s3.log.Errorf("Failed to check bucket exist", "error", errBucketExists)
		return errBucketExists
	}
	if !exists {
		err = minioClient.MakeBucket(context.Background(), s3.config.bucketName, minio.MakeBucketOptions{
			Region: s3.config.loaction,
		})
		if err != nil {
			s3.log.Errorf("Failed to create bucket when bucket not exist", "error", err)
			return err
		}
	}

	s3.log.Infof("Init s3 client successfully")
	s3.client = minioClient
	return nil
}

func (s3 *S3ServiceImpl) initContentMap() {
	s3.fileContentMap = make(map[string]string)
	contentMapStrs := strings.Split(s3.config.contentTyepMapStr, ";")
	for _, str := range contentMapStrs {
		ms := strings.Split(str, ":")
		if len(ms) != 2 {
			continue
		}
		s3.fileContentMap[ms[0]] = ms[1]
	}

	if len(s3.fileContentMap) == 0 {
		s3.fileContentMap[".tiff"] = "image/tiff"
		s3.fileContentMap[".csv"] = "application/csv"
		s3.fileContentMap[".png"] = "image/png"
	}
}

func (s3 *S3ServiceImpl) ExistObject(objectPath, fileName string) (bool, error) {
	objectName := filepath.Join(objectPath, fileName)
	objectInfo, err := s3.client.StatObject(context.Background(), s3.config.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return false, err
	}

	if objectInfo.ETag != "" && strings.Contains(objectInfo.Key, fileName) {
		return true, nil
	}
	return false, nil
}

func (s3 *S3ServiceImpl) UploadFile(objectPath string, fileName string, fileContent []byte) error {
	objectName := filepath.Join(objectPath, fileName)
	tempFloderId := uuid.New().String()
	tempFloderPath := filepath.Join(s3.uploadBaseFilePath, tempFloderId)

	if err := os.MkdirAll(tempFloderPath, os.ModePerm); err != nil {
		return err
	}
	defer os.RemoveAll(tempFloderPath)

	uploadFilePath := filepath.Join(tempFloderPath, fileName)
	file, err := os.Create(uploadFilePath)
	if err != nil {
		return fmt.Errorf("error creating file.error:%s", err.Error())
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(fileContent))
	if err != nil {
		return fmt.Errorf("error copy file.error:%s", err.Error())
	}

	contentType := s3.GetContentTypeByFileName(fileName)
	n, err := s3.client.FPutObject(context.Background(), s3.config.bucketName,
		makesureS3ObjectName(objectName), uploadFilePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		s3.log.Debugf("Faile to upload file to s3", "error", err.Error())
		return err
	}
	s3.log.Infof("file upload to s3 success", "s3file", objectName, "size", n)
	return nil
}

func (s3 *S3ServiceImpl) DownloadFile(objectPath string, fileName string) ([]byte, error) {
	objectName := filepath.Join(objectPath, fileName)
	tempFloderId := uuid.New().String()
	tempFloderPath := filepath.Join(s3.downloadFilePath, tempFloderId)

	if err := os.MkdirAll(tempFloderPath, os.ModePerm); err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempFloderPath)

	downloadFilePath := filepath.Join(tempFloderPath, fileName)
	if err := s3.client.FGetObject(context.Background(), s3.config.bucketName, makesureS3ObjectName(objectName),
		downloadFilePath, minio.GetObjectOptions{}); err != nil {
		s3.log.Debugf("Failed to download file from s3", "error", err.Error())
		return nil, err
	}
	fbytes, err := os.ReadFile(downloadFilePath)
	if err != nil {
		return nil, err
	}
	return fbytes, nil
}

func (s3 *S3ServiceImpl) DeleteS3Object(objName string) error {
	return s3.client.RemoveObject(context.Background(), s3.config.bucketName, makesureS3ObjectName(objName), minio.RemoveObjectOptions{})
}

func (s3 *S3ServiceImpl) MoveObject(sourceObjName, targetObjName string) error {
	dst := minio.CopyDestOptions{
		Bucket: s3.config.bucketName,
		Object: makesureS3ObjectName(targetObjName),
	}
	src := minio.CopySrcOptions{
		Bucket: s3.config.bucketName,
		Object: makesureS3ObjectName(sourceObjName),
	}
	_, err := s3.client.CopyObject(context.Background(), dst, src)
	if err != nil {
		s3.log.Debugf("Failed to move file from s3", "error", err.Error())
		return err
	}
	return s3.DeleteS3Object(sourceObjName)
}

func (s3 *S3ServiceImpl) GetContentTypeByFileName(fileName string) string {
	fileType := util.GetFileTypeByFileName(fileName)
	if s3.fileContentMap[fileType] != "" {
		return s3.fileContentMap[fileType]
	}
	return "image/png"
}

func makesureS3ObjectName(name string) string {
	if strings.HasPrefix(name, "/") {
		return name[1:]
	}
	return name
}
