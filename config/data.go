package config

type Conf struct {
	Fmwk    *FmwkConf
	App     *AppConf
	Db      *DbConf
	Godoc   *GodocConf
	Cluster *ClusterConf
}

func Data() Conf {
	return Conf{&Fmwk, &App, &DB, &Godoc, &Cluster}
}
