package app

type Image struct {
}

func (i Image) InstanceEnvName() string {
	return `GOPORT`
}

func (i Image) EnvironmentEnvName() string {
	return `GOENV`
}

func (i Image) OptionsForRun() []string {
	return []string{`-e=GODEV=true`}
}

func (i Image) Prepare() error {
	if err := compile(); err != nil {
		return err
	}
	/*
		if err :=	Assets(nil); err != nil {
			return err
		}
	*/
	return nil
}
