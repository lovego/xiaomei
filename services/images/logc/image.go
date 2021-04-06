package logc

type Image struct {
}

func (i Image) EnvironmentEnvVar() string {
	return `GoEnv`
}
