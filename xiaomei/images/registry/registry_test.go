package registry

import (
	"fmt"
	"testing"
)

func TestDigest(t *testing.T) {
	fmt.Println(Digest(`192.168.202.12:5000/example/app`, `latest`))
}

func TestTags(t *testing.T) {
	fmt.Println(Tags(`192.168.202.12:5000/example/app`))
}

func TestRemove(t *testing.T) {
	Remove(
		`192.168.202.12:5000/example/app`,
		`sha256:5a0163196bc9cd3e3c50579b309dd122633c7d29fa979c077c8e7a12ad50ff24`,
	)
}
