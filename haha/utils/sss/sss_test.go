package sss

import (
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/textproto"
	"testing"
)

func TestSetAvatarNoFile(t *testing.T) {
	form := multipart.Form{
		File: map[string][]*multipart.FileHeader{},
	}

	fileHeader := &multipart.FileHeader{
		Filename: "File",
		Header:   textproto.MIMEHeader{},

		Size: 0,
	}
	form.File["file"] = append(form.File["file"], fileHeader)

	fileHeader.Header.Add("Content-Type", "text/plain")

	_, err := UploadAvatar(&form, uint64(1))
	assert.Error(t, err)
}
