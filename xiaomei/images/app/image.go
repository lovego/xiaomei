package app

type Image struct {
}

func (i Image) InstanceEnvName() string {
	return `GOPORT`
}

func (i Image) EnvironmentEnvName() string {
	return `GOENV`
}

func (i Image) RunEnvName() string {
	return `GODEV`
}

func (i Image) Prepare() error {
	if err := buildBinary(); err != nil {
		return err
	}
	/*
		if err :=	Assets(nil); err != nil {
			return err
		}
	*/
	return nil
}
