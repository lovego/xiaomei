package release

var getEnvs func() []string

func RegisterEnvsGetter(getter func() []string) {
	getEnvs = getter
}
