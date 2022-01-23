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
}
