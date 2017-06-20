package errs

type Error interface {
	Code() string
	Message() string
}

func New(code, message string) Error {
	return errCodeMessage{code, message}
}

type errCodeMessage struct {
	code, message string
}

func (err errCodeMessage) Code() string {
	return err.code
}

func (err errCodeMessage) Message() string {
	return err.message
}
