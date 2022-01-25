package golog

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	DEBUG = Level{
		Value: 10,
		color: "blue",
		icon:  "★",
		Name:  "DEBUG",
	}

	INFO = Level{
		Value: 20,
		color: "green",
		icon:  "♥",
		Name:  "INFO",
	}

	WARN = Level{
		Value: 30,
		color: "yellow",
		icon:  "\u26A0",
		Name:  "WARN",
	}

	ERROR = Level{
		Value: 40,
		color: "red",
		Name:  "ERROR",
		icon:  "✖",
	}

	PANIC = Level{
		Value: 50,
		color: "black",
		icon:  "☹",
		Name:  "PANIC",
	}

	maxnamelen = 20
	curnamelen = 7

	separators []byte = []byte{'/', '.', '-'}
)

type Ctx map[string]interface{}

type Level struct {
	Value int `json:"value"`

	color string

	icon string

	Name string `json:"name"`
}

type Log struct {
	Time time.Time `json:"time"`

	Message string `json:"message"`

	Level Level `json:"level"`

	Data []interface{} `json:"data"`

	Ctx Ctx `json:"ctx"`

	Pid int `json:"pid"`

	Logger *Logger `json:"logger"`
}

type Logger struct {
	appenders []Appender
	disabled  bool

	Name string `json:"name"`

	Level Level `json:"-"`

	DoPanic bool `json:"-"`

	ctx Ctx
}

func (l *Logger) shouldAppend(lvl Level) bool {
	if l.disabled || lvl.Value < l.Level.Value {
		return false
	}

	return true
}

func (l *Logger) makeLog(msg interface{}, lvl Level, data []interface{}) {
	log := Log{
		Time:    time.Now().UTC(),
		Message: l.toString(msg),
		Level:   lvl,
		Data:    data,
		Logger:  l,
		Pid:     os.Getpid(),
		Ctx:     l.ctx,
	}

	for _, appender := range l.appenders {
		appender.Append(log)
	}
}

func (l *Logger) toString(object interface{}) string {
	return fmt.Sprintf("%v", object)
}

func (l *Logger) normalizeName() {
	length := len(l.Name)

	if length == maxnamelen || length == curnamelen {
		return
	}

	if length < curnamelen {
		l.normalizeNameLen()
		return
	}

	var (
		normalized string
		parts      []string
		separator  byte
	)

	for _, sep := range separators {
		parts = strings.Split(l.Name, string(sep))
		if len(parts) > 1 {
			separator = sep
			break
		}
	}

	if len(parts) > 1 {
		appendSeparator := true

		for i, str := range parts {
			switch len(str) {
			case 0:
				appendSeparator = false
				break
			case 1:
				normalized += str[:1]
				break
			case 2:
				normalized += str[:2]
				break
			default:
				normalized += str[:3]
				break
			}

			if appendSeparator && (i != (len(parts) - 1)) {
				normalized += string(separator)
			}
		}

		if len(normalized) > maxnamelen {
			normalized = normalized[:maxnamelen]
		}
	} else {
		length := len(l.Name)
		if length > maxnamelen {
			normalized = l.Name[:maxnamelen]
		} else {
			normalized = l.Name[0:length]
		}
	}

	l.Name = normalized
	if len(normalized) >= curnamelen {
		curnamelen = len(normalized)
	} else {
		l.normalizeNameLen()
	}
}

func (l *Logger) normalizeNameLen() {
	length := len(l.Name)
	missing := curnamelen - length
	for i := 0; i < missing; i++ {
		l.Name += " "
	}
}

func (l *Logger) Debug(msg interface{}, data ...interface{}) {
	if l.shouldAppend(DEBUG) {
		l.makeLog(msg, DEBUG, data)
	}
}

func (l *Logger) Info(msg interface{}, data ...interface{}) {
	if l.shouldAppend(INFO) {
		l.makeLog(msg, INFO, data)
	}
}

func (l *Logger) Warn(msg interface{}, data ...interface{}) {
	if l.shouldAppend(WARN) {
		l.makeLog(msg, WARN, data)
	}
}

func (l *Logger) Error(msg interface{}, data ...interface{}) {
	if l.shouldAppend(ERROR) {
		l.makeLog(msg, ERROR, data)
	}
}

func (l *Logger) Panic(msg interface{}, data ...interface{}) {
	if l.shouldAppend(PANIC) {
		l.makeLog(msg, PANIC, data)
		panic(msg)
	}
}

func (l *Logger) Debugf(msg string, params ...interface{}) {
	if l.shouldAppend(DEBUG) {
		l.makeLog(fmt.Sprintf(msg, params...), DEBUG, nil)
	}
}

func (l *Logger) Infof(msg string, params ...interface{}) {
	if l.shouldAppend(INFO) {
		l.makeLog(fmt.Sprintf(msg, params...), INFO, nil)
	}
}

func (l *Logger) Warnf(msg string, params ...interface{}) {
	if l.shouldAppend(WARN) {
		l.makeLog(fmt.Sprintf(msg, params...), WARN, nil)
	}
}

func (l *Logger) Errorf(msg string, params ...interface{}) {
	if l.shouldAppend(ERROR) {
		l.makeLog(fmt.Sprintf(msg, params...), ERROR, nil)
	}
}

func (l *Logger) Panicf(msg string, params ...interface{}) {
	if l.shouldAppend(PANIC) {
		l.makeLog(fmt.Sprintf(msg, params...), PANIC, nil)
		panic(msg)
	}
}

func (l *Logger) Enable(appender Appender) {
	l.appenders = append(l.appenders, appender)
}

func (l *Logger) Disable(target interface{}) {
	var id string
	var appender Appender

	switch object := target.(type) {
	case string:
		id = object
	case Appender:
		appender = object
	default:
		l.Warn("Error while disabling logger. Cannot cast to target type.")
		return
	}

	for i, app := range l.appenders {
		if (appender != nil && (app == appender || appender.Id() == app.Id())) || id == app.Id() {
			var toAppend []Appender

			if len(l.appenders) >= i+1 {
				toAppend = l.appenders[i+1:]
			}

			l.appenders = append(l.appenders[:i], toAppend...)
			return
		}
	}
}

func (l *Logger) SetContext(ctx Ctx) *Logger {
	l.ctx = ctx
	return l
}

func (l *Logger) AddContextKey(key string, value interface{}) *Logger {
	l.ctx[key] = value
	return l
}

func (l *Logger) Copy() *Logger {
	ctxLogger := &Logger{}
	*ctxLogger = *l

	return ctxLogger
}
