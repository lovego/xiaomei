package config

type Conf struct {
	App     *AppConf
	Db      *DbConf
	Deploy  *DeployConf
	Godoc   *GodocConf
	Servers *ServerConf
}

func Data() Conf {
	return Conf{&App, &DB, &Deploy, &Godoc, &Servers}
}
