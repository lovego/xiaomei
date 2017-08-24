package web

type Image struct {
}

func (i Image) InstanceEnvName() string {
	return `NGINXPORT`
}

func (i Image) EnvsForRun() []string {
	return []string{`SendfileOff=true`}
}
