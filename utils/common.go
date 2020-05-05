package utils

import (
	"bytes"
	"fmt"
	"iron/log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// ExecCommand 执行命令，执行过程中函数阻塞
func ExecCommand(commandName string, params []string) (result string, err error) {
	if commandName == `` {
		err = fmt.Errorf("commandName is null")
		log.Error(err.Error())
		return ``, err
	}
	if params == nil {
		err = fmt.Errorf(`params is nil`)
		log.Error(err.Error())
		return ``, err
	}
	cmd := exec.Command(commandName, params...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf(err.Error() + `:` + stderr.String())
		return ``, err
	}
	result = out.String()
	return
}

// Underline2hump 名字转化，将下划线转驼峰
func Underline2hump(s string) (r string) {
	if s == `` {
		return ``
	}
	if s == `id` {
		return `ID`
	}
	ss := strings.Split(s, `_`)
	for _, v := range ss {
		if v == `` {
			continue
		}
		if v == `id` {
			r = r + `ID`
			continue
		}
		r = r + strings.Title(v)
	}
	return r
}

// FindAllFile 在某个目录下查询所有的文件(包含路径 eg:src/file.txt)
func FindAllFile(filePath string) (allFiles []string, err error) {
	log.Info(filePath)
	if filePath == `` {
		err = fmt.Errorf("filePath is empty")
		log.Error(err.Error())
		return nil, err
	}
	allFiles = make([]string, 0)
	err = filepath.Walk(filePath,
		func(path string, f os.FileInfo, err error) error {
			if err != nil {
				log.Error(err.Error())
				return err
			}
			if f == nil {
				log.Error(err.Error())
				return err
			}
			if path == `` {
				err = fmt.Errorf("path is null")
				log.Error(err.Error())
				return err
			}
			//判断是否是文件夹，如果是文件夹，直接返回，不读取
			if f.IsDir() {
				return nil
			}
			allFiles = append(allFiles, path)
			return nil
		})
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return allFiles, nil
}

// RemoveBlankLine 删除文本中的空行
func RemoveBlankLine(text string) (result string, err error) {
	if text == `` {
		return ``, nil
	}
	reg, err := regexp.Compile(`\n(\s*)\n`)
	if err != nil {
		log.Error(err.Error())
		return ``, err
	}
	// log.Info(reg.String())
	result = reg.ReplaceAllString(text, "\n")
	return result, nil
}

// GenUniqueFileName 生成唯一文件名,通过uuid包实现
// suffixName:后缀名
func GenUniqueFileName(suffixName string) string {
	u2 := uuid.NewV4()
	log.Info(u2.String())
	return u2.String() + `.` + suffixName
}
