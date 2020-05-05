package base

import (
	"fmt"
	"os"
	"path/filepath"
)

// BasePath 当前路径
var BasePath string

// initBasePath 初始化当前路径
func initBasePath() (err error) {
	BasePath, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(`initBase Path has error: `, err)
		return err
	}
	return nil
}
