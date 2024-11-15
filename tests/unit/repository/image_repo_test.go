package repository_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestSaveImage_Success(t *testing.T) {
	images := SampleImagePNG()
	imageData, err := base64.StdEncoding.DecodeString(images)
	if err != nil {
		t.Fatal(err)
	}
	mimeType := http.DetectContentType(imageData)
	mimeTypePrefix := strings.Split(mimeType, "/")[1]

	id := "PR-"
	ms := time.Now().UnixNano() / int64(time.Millisecond)
	imageName := fmt.Sprintf("%s-%d.%s", id, ms, mimeTypePrefix)

	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		t.Fatal(err)
	}
	compressedImage, err := compressImage(img, format)
	if err != nil {
		t.Fatal(err)
	}

	err = imageRepository.SaveImage(compressedImage, imageName)
	if err != nil {
		t.Fatal(err)
	}
}
