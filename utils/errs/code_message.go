package errs

type CodeMessageErr interface {
	Code() string
	Message() string
	Error() string
}

func New(code, message string) CodeMessageErr {
	return codeMessage{code, message}
}

type codeMessage struct {
	code, message string
}

func (err codeMessage) Code() string {
	return err.code
}

func (err codeMessage) Message() string {
	return err.message
}

func (err codeMessage) Error() string {
	return err.code + `: ` + err.message
}
