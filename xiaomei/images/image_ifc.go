package images

import (
	"path/filepath"

	"github.com/lovego/xiaomei/xiaomei/release"
)

// the 7 interfaces, image driver choose to implement.

// 1. for simple deploy
func (i Image) InstanceEnvName() string {
	if ifc, ok := i.image.(interface {
		InstanceEnvName() string
	}); ok {
		return ifc.InstanceEnvName()
	}
	return ``
}

// 2. envs for deploy and run
func (i Image) Envs() []string {
	if ifc, ok := i.image.(interface {
		Envs() []string
	}); ok {
		return ifc.Envs()
	}
	return nil
}

// 3. envs for run
func (i Image) EnvsForRun() []string {
	if ifc, ok := i.image.(interface {
		EnvsForRun() []string
	}); ok {
		return ifc.EnvsForRun()
	}
	return nil
}

// 4. prepare files for build
func (i Image) prepare() error {
	if ifc, ok := i.image.(interface {
		Prepare() error
	}); ok {
		return ifc.Prepare()
	}
	return nil
}

// 5. dir for build
func (i Image) buildDir() string {
	if ifc, ok := i.image.(interface {
		BuildDir() string
	}); ok {
		return ifc.BuildDir()
	}
	return filepath.Join(release.Root(), `img-`+i.svcName)
}

// 6. dockerfile for build
func (i Image) dockerfile() string {
	if ifc, ok := i.image.(interface {
		Dockerfile() string
	}); ok {
		return ifc.Dockerfile()
	}
	return `Dockerfile`
}
