package host

import (
	"fmt"
)

var Driver driver

type driver struct{}

func (d driver) ServiceNames() map[string]bool {
	m := make(map[string]bool)
	for svcName := range getRelease() {
		m[svcName] = true
	}
	return m
}

func (d driver) ImageNameOf(svcName string) string {
	svc := getService(svcName)
	if svc.Image == `` {
		panic(fmt.Sprintf(`release.yml: %s.image: empty.`, svcName))
	}
	return svc.Image
}
