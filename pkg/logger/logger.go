package logger

type Logger interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Printf(format string, args ...interface{})
}

type Dummy struct {
}

func (d Dummy) Errorf(string, ...interface{}) {}
func (d Dummy) Fatalf(string, ...interface{}) {}
func (d Dummy) Infof(string, ...interface{})  {}
func (d Dummy) Warnf(string, ...interface{})  {}
func (d Dummy) Debugf(string, ...interface{}) {}
func (d Dummy) Printf(string, ...interface{}) {}
