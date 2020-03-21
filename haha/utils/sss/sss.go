package sss

import (
	"bytes"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"joblessness/haha/auth/interfaces"
	"mime/multipart"
	"strconv"
	"strings"
)

var Sess = session.Must(session.NewSession(&aws.Config{
	Region: aws.String("ru-msk"),
	Credentials: credentials.NewStaticCredentials("orFNtcQG9pi8NvqcFhLAj4",
		"33CiuS769M4u1wHAk42HhdtCrCb795MGuez3biaE3CeK", ""),
	Endpoint: aws.String("https://hb.bizmrg.com"),
}))

var Svc = s3.New(Sess)

func UploadAvatar(form *multipart.Form, userID uint64) (link string, err error) {
	fileHeaders := form.File["file"]
	if len(fileHeaders) == 0 {
		return "", errors.New("no file in multipart form")
	}

	file, err := fileHeaders[0].Open()
	defer file.Close()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	io.Copy(&buf, file)
	splitName := strings.Split(fileHeaders[0].Filename, ".")
	ext := splitName[len(splitName)-1]

	link = strconv.FormatUint(userID, 10) + "-avatar." + ext

	_, err = Svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("imgs-hh"),
		Key:    aws.String(link),
		Body:   strings.NewReader(buf.String()),
		ACL:    aws.String("public-read"), // make public
	})

	if err != nil {
		return "", authInterfaces.ErrUploadAvatar
	}

	return link, nil
}