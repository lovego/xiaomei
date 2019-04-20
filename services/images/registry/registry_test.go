// +build registry

package registry

import (
	"fmt"
	"testing"
)

func TestDigest(t *testing.T) {
	fmt.Println(Digest(`registry.cn-beijing.aliyuncs.com/lovego/ubuntu`, `latest`))
}

func TestTags(t *testing.T) {
	fmt.Println(Tags(`registry.cn-beijing.aliyuncs.com/lovego/ubuntu`))
}

func TestRemove(t *testing.T) {
	Remove(
		`registry.cn-beijing.aliyuncs.com/lovego/ubuntu`,
		`sha256:5a0163196bc9cd3e3c50579b309dd122633c7d29fa979c077c8e7a12ad50ff24`,
	)
}
