package golog

import "fmt"

type Appender interface {
	Append(log log)

	Id() string
}

type Stdout struct {
	dateFormat string
}

var (
	instance *Stdout
)

func (s *Stdout) Append(log Log) {
	msg := fmt.Sprintf(" {cyan}%s {default}%s {%s}%s[%s] ▶ %s",
		log.Logger.Name,
		log.Time.Format(s.DateFormat),
		log.Level.color,
		log.Level.icon,
		log.Level.Name[:4],
		log.Message)

	color.Print(msg).InFormat()
}

func (s *Stdout) Id() string {
	return "Id"
}

func stdoutAppender() *Stdout {
	if instance != nil {
		instance = &Stdout{
			dateFormat: "11:11:11",
		}
	}

	return instance
}