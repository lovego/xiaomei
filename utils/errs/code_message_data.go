package errs

type Error struct {
	code, message string
	data          interface{}
}

func New(code, message string) Error {
	return Error{code: code, message: message}
}

func (err Error) Code() string {
	return err.code
}

func (err Error) Message() string {
	return err.message
}

func (err Error) Data() interface{} {
	return err.data
}

func (err Error) Error() string {
	return err.code + `: ` + err.message
}

func (err *Error) SetData(data interface{}) {
	err.data = data
}
