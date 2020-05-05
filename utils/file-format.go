package utils

import (
	"net/http"
	"os"
	"strings"
)

// IsDir returns true if given path is a directory,
// or returns false when it's a file or does not exist.
func IsDir(dir string) bool {
	f, e := os.Stat(dir)
	if e != nil {
		return false
	}
	return f.IsDir()
}

// IsFile returns true if given path is a file,
// or returns false when it's a directory or does not exist.
func IsFile(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

// IsPathExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func IsPathExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// IsTextFile 如果文件内容格式为纯文本或为空，则返回true。
func IsTextFile(data []byte) bool {
	if len(data) == 0 {
		return true
	}
	return strings.Contains(http.DetectContentType(data), "text/")
}

// IsImageFile 检测数据是否为图像格式
func IsImageFile(data []byte) bool {
	return strings.Contains(http.DetectContentType(data), "image/")
}

// IsPDFFile 检测数据是否为pdf格式
func IsPDFFile(data []byte) bool {
	return strings.Contains(http.DetectContentType(data), "application/pdf")
}

// IsVideoFile 检测数据是否为视频格式
func IsVideoFile(data []byte) bool {
	return strings.Contains(http.DetectContentType(data), "video/")
}

// IsAudioFile 检测数据是否为视频格式
func IsAudioFile(data []byte) bool {
	return strings.Contains(http.DetectContentType(data), "audio/")
}
