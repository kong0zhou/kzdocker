package base

import (
	"github.com/spf13/viper"
)

// https://mholt.github.io/json-to-go/
type config struct {
	ZapLog struct {
		Console struct {
			Enable bool `mapstructure:"enable"`
		} `mapstructure:"console"`
		Postgres struct {
			Enable    bool   `mapstructure:"enable"`
			Host      string `mapstructure:"host"`
			Port      string `mapstructure:"port"`
			User      string `mapstructure:"user"`
			Password  string `mapstructure:"password"`
			Sslmode   string `mapstructure:"sslmode"`
			Dbname    string `mapstructure:"dbname"`
			Schemas   string `mapstructure:"schemas"`
			Tablename string `mapstructure:"tablename"`
		} `mapstructure:"postgres"`
		File struct {
			Enable  bool   `mapstructure:"enable"`
			RelPath string `mapstructure:"relPath"`
		} `mapstructure:"file"`
	} `mapstructure:"zapLog"`
}

// Config 存储所有的配置
var Config config

// initConf 初始化配置文件
func initConf() (err error) {
	v := viper.New()

	// 从指定的文件中读取配置文件
	v.SetConfigFile(BasePath + `/config.json`)

	err = v.ReadInConfig()
	if err != nil {
		return err
	}

	err = v.Unmarshal(&Config)
	if err != nil {
		return err
	}
	return nil
}
