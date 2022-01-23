package append

import (
	"fmt"
	"github.com/godzillaframework/golog"
	"os"
)

type fileAppender struct {
	path string
}

func (fa *fileAppender) Id() string {
	return "ID"
}

func (fa *fileAppender) append(log golog.Log) {
	f, err := os.OpenFile(fa.path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err.Error())
		if log.Logger.doPanic {
			panic(err)
		}
	}

	defer f.Close()

}

func file(cnf golog.Conf) *fileAppender {
	return &fileAppender{
		path: cnf["path"],
	}
}