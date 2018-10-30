// +build registry

package registry

import (
	"fmt"
	"testing"
)

func TestDigest(t *testing.T) {
	fmt.Println(Digest(`hub.c.163.com/lovego/xiaomei/test`, `latest`))
}

func TestTags(t *testing.T) {
	fmt.Println(Tags(`hub.c.163.com/lovego/xiaomei/test`))
}

func TestRemove(t *testing.T) {
	Remove(
		`hub.c.163.com/lovego/xiaomei/test`,
		`sha256:5a0163196bc9cd3e3c50579b309dd122633c7d29fa979c077c8e7a12ad50ff24`,
	)
}
