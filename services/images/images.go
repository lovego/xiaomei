package images

import (
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/lovego/cmd"
	"github.com/lovego/config/config"
	"github.com/lovego/xiaomei/release"
)

type Build struct {
	Env, Tag         string
	PrepareFlags     []string
	DockerBuildFlags []string
}

func (b Build) Run(svcName string) error {
	return imagesDo(svcName, b.Env, func(img Image) error {
		svcDir := release.ImageDir(b.Env, img.svcName)
		if err := img.prepare(b.Env, svcDir, b.PrepareFlags); err != nil {
			return err
		}

		log.Println(color.GreenString(`building ` + img.svcName + ` image.`))

		_, err := cmd.Run(cmd.O{Dir: svcDir, Print: true}, `docker`, b.args(img)...)
		return err
	})
}

func (b Build) args(img Image) []string {
	var result = []string{
		`build`,
		`--tag=` + release.GetService(b.Env, img.svcName).ImageName(b.Tag),
	}
	environment := config.NewEnv(b.Env)
	for _, envVar := range environment.Vars() {
		result = append(result, `--build-arg`, envVar)
	}
	if len(b.DockerBuildFlags) == 0 {
		result = append(result, `--pull`)
	} else {
		result = append(result, b.DockerBuildFlags...)
	}

	return append(result, `.`)
}

type Push struct {
	Env, Tag string
	DockerLogin
}

func (p Push) Run(svcName string) error {
	return imagesDo(svcName, p.Env, func(img Image) error {
		log.Println(color.GreenString(`pushing ` + img.svcName + ` image.`))
		imgName := release.GetService(p.Env, img.svcName).ImageName(p.Tag)
		if err := p.DockerLogin.Run(imgName); err != nil {
			return err
		}
		_, err := cmd.Run(cmd.O{Print: true}, `docker`, `push`, imgName)
		return err
	})
}

type DockerLogin struct {
	User              string
	Password          string
	loginedRegistries map[string]bool
}

func (dl *DockerLogin) Run(imgName string) error {
	if dl.empty() {
		return nil
	}
	registry := DockerRegistry(imgName)
	if dl.loginedRegistries[registry] {
		return nil
	}

	command := dl.Command(registry, false)
	if len(command) == 0 {
		return nil
	}
	_, err := cmd.Run(cmd.O{}, command[0], command[1:]...)
	if err == nil {
		if dl.loginedRegistries == nil {
			dl.loginedRegistries = map[string]bool{}
		}
		dl.loginedRegistries[registry] = true
	}
	return err
}

func (dl *DockerLogin) BashCommand(env string, svcs []string) string {
	if dl.empty() {
		return ""
	}
	var commands []string
	var registries = map[string]bool{}
	for _, svc := range svcs {
		registry := DockerRegistry(release.GetService(env, svc).Image)
		if !registries[registry] {
			commands = append(commands, strings.Join(dl.Command(registry, true), " "))
		}
	}
	return strings.Join(commands, "\n")
}

func (dl *DockerLogin) Command(registry string, bashQuote bool) []string {
	user, password := dl.User, dl.Password
	if bashQuote {
		user, password = release.BashQuote(user), release.BashQuote(password)
	}
	return []string{`docker`, `login`, `-u=` + user, `-p=` + password, registry}
}

func (dl *DockerLogin) empty() bool {
	return dl.User == "" && dl.Password == ""
}

func DockerRegistry(imgName string) string {
	i := strings.IndexByte(imgName, '/')
	if i <= 0 {
		panic("unknown registry for image: " + imgName)
	}
	return imgName[:i]
}

func List(svcName, env string) error {
	return imagesDo(svcName, env, func(img Image) error {
		if svcName == `` {
			color.Green(img.svcName + `:`)
		}
		_, err := cmd.Run(cmd.O{}, `docker`, `images`,
			`-f`, `reference=`+release.GetService(env, img.svcName).ImageName(``))
		return err
	})
}

func imagesDo(svcName, env string, work func(Image) error) error {
	if svcName == `` {
		for _, svcName := range release.ServiceNames(env) {
			if err := work(Get(svcName)); err != nil {
				return err
			}
		}
		return nil
	} else {
		return work(Get(svcName))
	}
}
