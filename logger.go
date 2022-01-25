package golog

import "time"

type Ctx map[string]interface{}

type Level struct {
	Value int `json:"value"`
	color string
	icon  string
	Name  string `json:"name"`
}

type Log struct {
	Time    time.Time     `json:"time"`
	Message string        `json:"message"`
	Level   Level         `json:"level"`
	Data    []interface{} `json:"data"`
	Ctx     Ctx           `json:"ctx"`
	Pid     int           `json:"pid"`
	Logger  *Logger       `json:"logger"`
}

type Logger struct {
	appenders []Appender
	disabled  bool
	Name      string `json:"name"`
	Level     Level  `json:"-"`
	DoPanic   bool   `json:"-"`
	ctx       Ctx
}

func (l *Logger) setContext(ctx Ctx) *Logger {
	l.ctx = ctx
	return l
}

func (l *Logger) addContextKey(key string, value interface{}) *Logger {
	l.ctx[key] = value
	return l
}

func (l *Logger) Copy() *Logger {
	ctxLogger := &Logger{}
	*ctxLogger = *1
	return ctxLogger
}
