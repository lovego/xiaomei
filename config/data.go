package config

type Conf struct {
	Fmwk    *FmwkConf
	App     *AppConf
	Db      *DbConf
	Deploy  *DeployConf
	Godoc   *GodocConf
	Servers *ServerConf
}

func Data() Conf {
	return Conf{&Fmwk, &App, &DB, &Deploy, &Godoc, &Servers}
}
