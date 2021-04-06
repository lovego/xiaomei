package app

type Image struct {
}

func (i Image) EnvironmentEnvVar() string {
	return `GoEnv`
}

func (i Image) PortEnvVar() string {
	return `GOPORT`
}

func (i Image) DefaultPort() uint16 {
	return 3000
}

func (i Image) OptionsForRun() []string {
	return []string{`-e=GODEV=true`}
}

func (i Image) Prepare() error {
	if err := compile(true); err != nil {
		return err
	}
	/*
		if err :=	Assets(nil); err != nil {
			return err
		}
	*/
	return nil
}
