package logc

type Image struct {
}

func (i Image) EnvironmentEnvName() string {
	return `GOENV`
}
