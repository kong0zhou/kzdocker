package utils

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"runtime"
// 	"strings"

// 	"github.com/sirupsen/logrus"
// )

// // Log logger对象
// var Log = logrus.New()
// var log = Log

// // InitLogger 初始化Log对象
// func InitLogger() (err error) {
// 	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
// 	if err != nil {
// 		fmt.Println(`initLogger error:`, err)
// 		return
// 	}

// 	Log.SetReportCaller(true)
// 	//设置日志级别
// 	Log.SetLevel(logrus.DebugLevel)
// 	//格式设置为文本格式
// 	Log.SetFormatter(&logrus.TextFormatter{
// 		TimestampFormat: "2006-01-02 15:04:05",
// 		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
// 			var funcName string
// 			funcIndex := strings.LastIndex(f.Function, ".")
// 			if funcIndex == -1 {
// 				funcName = f.Function
// 			} else {
// 				funcName = f.Function[funcIndex+1:]
// 			}
// 			filePath := strings.Replace(f.File, path, "", 1)
// 			filePath = filePath[1:]
// 			return fmt.Sprintf("%s()", funcName), fmt.Sprintf(" %s:%d", filePath, f.Line)
// 		},
// 	})
// 	//如果是生产环境
// 	if Config.Production {
// 		err := os.MkdirAll("logs", 0755)
// 		if err != nil {
// 			fmt.Println(`initLogger error:`, err)
// 			return err
// 		}
// 		//   设置文件为日志输出
// 		file, err := os.OpenFile("logs/log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
// 		if err != nil {
// 			fmt.Println("init Logger error:", err)
// 			return err
// 		}
// 		Log.Out = file
// 	} else {
// 		Log.Out = os.Stdout
// 	}
// 	return nil
// }
