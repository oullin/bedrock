package testing

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"

	"github.com/bedrock/packages/httpx"
)

// FakeFile creates a fake UploadedFile with random content for testing.
func FakeFile(name string, sizeKB int) *httpx.UploadedFile {
	content := make([]byte, sizeKB*1024)
	rand.Read(content)

	encoded := base64.StdEncoding.EncodeToString(content)
	file, _ := httpx.CreateFromBase64(encoded, name)

	return file
}

// FakeImage creates a fake PNG image UploadedFile with the given dimensions.
func FakeImage(name string, width, height int) *httpx.UploadedFile {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 0, G: 100, B: 200, A: 255})
		}
	}

	var buf bytes.Buffer
	png.Encode(&buf, img)

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	file, _ := httpx.CreateFromBase64(encoded, name)

	return file
}

// FakeFileWithContent creates a fake UploadedFile with specific content.
func FakeFileWithContent(name string, content []byte) *httpx.UploadedFile {
	encoded := base64.StdEncoding.EncodeToString(content)
	file, _ := httpx.CreateFromBase64(encoded, name)

	return file
}
