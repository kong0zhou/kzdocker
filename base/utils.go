package base

import "os"

// isFile returns true if given path is a file,
// or returns false when it's a directory or does not exist.
func isFile(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

// isPathExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func isPathExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
