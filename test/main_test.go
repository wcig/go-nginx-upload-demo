package test

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/go-resty/resty/v2"
)

const (
	url         = "http://localhost:8001/upload"
	file        = "/tmp/test.mp4"
	fileName    = "test.mp4"
	segmentSize = int64(1048576) // 单个片段上传的字节数
)

// 简单上传
func TestSimpleUpload(t *testing.T) {
	client := resty.New()
	resp, err := client.R().SetFile("file", file).Post(url)
	fmt.Println(err, resp)
}

// 分片上传
func TestMultipartUpload(t *testing.T) {
	fileMd5, _ := FileMD5(file)
	fileSize, _ := FileSize(file)
	headers := map[string]string{
		"Session-ID":          fileMd5,
		"Content-Type":        "application/octet-stream",
		"Content-Disposition": "attachment; filename=" + fileName,
	}

	filePos := int64(0)
	size := int64(0)
	file, _ := os.Open(file)

	for {
		if filePos+segmentSize >= fileSize {
			size = fileSize - filePos
			upload(filePos, size, fileSize, file, headers)
			break
		} else {
			size = segmentSize
			upload(filePos, size, fileSize, file, headers)
			filePos += segmentSize
		}
	}
}

func upload(filePos, size, fileSize int64, file *os.File, headers map[string]string) {
	client := resty.New()
	headers["X-Content-Range"] = fmt.Sprintf("bytes %d-%d/%d", filePos, filePos+size-1, fileSize)
	headers["Content-Length"] = strconv.FormatInt(size, 10)

	file.Seek(filePos, 0)
	partBuffer := make([]byte, size)
	file.Read(partBuffer)
	resp, err := client.R().SetHeaders(headers).SetFileReader("file", fileName, bytes.NewReader(partBuffer)).Post(url)
	fmt.Println(err, resp)
}

func FileMD5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func FileSize(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}
