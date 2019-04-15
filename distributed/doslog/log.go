// 作为日志管理包
package doslog

import (
	"log"
	"os"
)

var Loger *log.Logger = log.New(os.Stderr, "", log.LstdFlags | log.Lshortfile)


func FailOnError(err error, message string) {
	if err != nil {
		Loger.Printf(message)
	}
}