package app

type Image struct {
}

func (i Image) EnvironmentEnvVar() string {
	return `ProENV`
}

func (i Image) PortEnvVar() string {
	return `ProPORT`
}

func (i Image) DefaultPort() uint16 {
	return 3000
}

func (i Image) OptionsForRun() []string {
	return []string{`-e=ProDEV=true`}
}

func (i Image) Prepare(goBuildFlags []string) error {
	if err := compile(true, goBuildFlags); err != nil {
		return err
	}
	// if err := Assets(nil); err != nil {
	// 	   return err
	// }
	return nil
}
