package log

import (
	"fmt"
	"runtime"
)

// InitLog init log
func InitLog() {
	fmt.Printf(`***********************************
******  go version %s  ******
***********************************
`, runtime.Version())
	err := initLogger()
	if err != nil {
		panic(err)
	}
}
