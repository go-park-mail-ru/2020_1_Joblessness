package sss

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

var sess = session.Must(session.NewSession(&aws.Config{
	Region: aws.String("ru-msk"),
	Credentials: credentials.NewStaticCredentials(os.Getenv("HOTBOX_ID"),
		os.Getenv("HOTBOX_SECRETE"), os.Getenv("HOTBOX_TOKEN")),
	Endpoint: aws.String("https://hb.bizmrg.com"),
}))

var svc = s3.New(sess)

func UploadAvatar(form *multipart.Form, userID uint64) (link string, err error) {
	fileHeaders := form.File["file"]
	if len(fileHeaders) == 0 {
		return "", errors.New("no file in multipart form")
	}

	file, err := fileHeaders[0].Open()
	if err != nil {
		return "", err
	}
	defer func() {
		err = file.Close()
	}()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	splitName := strings.Split(fileHeaders[0].Filename, ".")
	ext := splitName[len(splitName)-1]
	
	t := time.Now()
	link = fmt.Sprintf("%d%d%d%d%d%d-%d", t.Year(), 
			   t.Month(), t.Day(), t.Hour(), 
			   t.Minute(), t.Second(), userID) + "-avatar." + ext

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("imgs-hh"),
		Key:    aws.String(link),
		Body:   strings.NewReader(buf.String()),
		ACL:    aws.String("public-read"), // make public
	})

	if err != nil {
		return "", err
	}

	link = "https://hb.bizmrg.com/imgs-hh/" + link

	return link, nil
}
