package images

import (
	"path/filepath"

	"github.com/lovego/xiaomei/xiaomei/release"
)

// the 7 interfaces, image driver choose to implement.

// 1. instance env variable name
func (i Image) InstanceEnvName() string {
	if img, ok := i.image.(interface {
		InstanceEnvName() string
	}); ok {
		return img.InstanceEnvName()
	}
	return ``
}

// 2. environment env variable name
func (i Image) EnvironmentEnvName() string {
	if img, ok := i.image.(interface {
		EnvironmentEnvName() string
	}); ok {
		return img.EnvironmentEnvName()
	}
	return ``
}

// 3. options for run
func (i Image) OptionsForRun() []string {
	if img, ok := i.image.(interface {
		OptionsForRun() []string
	}); ok {
		return img.OptionsForRun()
	}
	return nil
}

// 4. prepare files for build
func (i Image) prepare() error {
	if img, ok := i.image.(interface {
		Prepare() error
	}); ok {
		return img.Prepare()
	}
	return nil
}

// 5. dir for build
func (i Image) buildDir() string {
	if img, ok := i.image.(interface {
		BuildDir() string
	}); ok {
		return img.BuildDir()
	}
	return filepath.Join(release.Root(), `img-`+i.svcName)
}

// 6. dockerfile for build
func (i Image) dockerfile() string {
	if img, ok := i.image.(interface {
		Dockerfile() string
	}); ok {
		return img.Dockerfile()
	}
	return `Dockerfile`
}
