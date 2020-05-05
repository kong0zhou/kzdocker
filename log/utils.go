package log

import "os"

// isPathExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func isPathExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
