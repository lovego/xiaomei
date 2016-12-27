package config

var Godoc GodocConf

type GodocConf struct {
	conf godocConf
}

type godocConf struct {
	Port          string `yaml:"port"`
	Domain        string `yaml:"domain"`
	IndexInterval string `yaml:"indexInterval"`
}

func (g *GodocConf) Port() string {
	Load()
	return g.conf.Port
}

func (g *GodocConf) Domain() string {
	Load()
	return g.conf.Domain
}

func (g *GodocConf) IndexInterval() string {
	Load()
	return g.conf.IndexInterval
}
