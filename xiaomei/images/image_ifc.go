package images

import (
	"path/filepath"

	"github.com/lovego/xiaomei/xiaomei/release"
)

// the 7 interfaces, image driver choose to implement.

// 1. environment env variable name
func (i Image) EnvironmentEnvVar() string {
	if img, ok := i.image.(interface {
		EnvironmentEnvVar() string
	}); ok {
		return img.EnvironmentEnvVar()
	}
	return ``
}

// 2. instance env variable name
func (i Image) PortEnvVar() string {
	if img, ok := i.image.(interface {
		PortEnvVar() string
	}); ok {
		return img.PortEnvVar()
	}
	return ``
}

// 3. default port number
func (i Image) DefaultPort() uint16 {
	if img, ok := i.image.(interface {
		DefaultPort() uint16
	}); ok {
		return img.DefaultPort()
	}
	return 0
}

// 4. options for run
func (i Image) OptionsForRun() []string {
	if img, ok := i.image.(interface {
		OptionsForRun() []string
	}); ok {
		return img.OptionsForRun()
	}
	return nil
}

// 5. prepare files for build
func (i Image) prepare() error {
	if img, ok := i.image.(interface {
		Prepare() error
	}); ok {
		return img.Prepare()
	}
	return nil
}

// 6. dir for build
func (i Image) buildDir() string {
	if img, ok := i.image.(interface {
		BuildDir() string
	}); ok {
		return img.BuildDir()
	}
	return filepath.Join(release.Root(), `img-`+i.svcName)
}

// 7. dockerfile for build
func (i Image) dockerfile() string {
	if img, ok := i.image.(interface {
		Dockerfile() string
	}); ok {
		return img.Dockerfile()
	}
	return `Dockerfile`
}
