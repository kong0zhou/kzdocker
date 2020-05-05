package base

import uuid "github.com/satori/go.uuid"

// AppUUID 每次启动程序生成一个uuid
var AppUUID = uuid.NewV4().String()

// AppVersion 程序的版本
var AppVersion = `0.01`

// InitBase base
func InitBase() {
	err := initBasePath()
	if err != nil {
		panic(err)
	}
	err = initConf()
	if err != nil {
		panic(err)
	}
}
