package golog

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

}
