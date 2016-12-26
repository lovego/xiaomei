package config

var Godoc godocVar

type godocVar struct {
	conf godocConf
}

type godocConf struct {
	Port          string `yaml:"port"`
	Domain        string `yaml:"domain"`
	IndexInterval string `yaml:"indexInterval"`
}

func (g *godocVar) Port() string {
	Load()
	return g.conf.Port
}

func (g *godocVar) Domain() string {
	Load()
	return g.conf.Domain
}

func (g *godocVar) IndexInterval() string {
	Load()
	return g.conf.IndexInterval
}
