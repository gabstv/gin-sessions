package sessions

type ErrorLogger interface {
	Errorln(args ...interface{})
}
