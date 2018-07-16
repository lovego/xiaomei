package registry

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/olekukonko/tablewriter"
)

func ListTimeTags(svcName, env string) {
	if svcName == `` {
		listSvcsTimeTags(env)
	} else {
		envTagsMap := getEnvTimeTagsMap(svcName, env)
		fmt.Println(strings.Join(envTagsMap[env], "\n"))
	}
}

func listSvcsTimeTags(env string) {
	svcs := conf.ServiceNames(env)
	svcsTags := make(map[string]map[string]bool) // map[svcName][env]tag
	for _, svcName := range svcs {
		svcsTags[svcName] = getTimeTagsMap(svcName, env)
	}
	tags := getUniqTimeTags(svcsTags)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorders(tablewriter.Border{false, false, false, false})
	table.SetColumnSeparator(``)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.Append(append([]string{``}, svcs...))
	for _, tag := range tags {
		table.Append(getTagRow(tag, svcs, svcsTags))
	}
	table.Render()
}

func getTagRow(tag string, svcs []string, svcsTags map[string]map[string]bool) []string {
	result := []string{tag}
	for _, svc := range svcs {
		if svcsTags[svc][tag] {
			result = append(result, `*`)
		} else {
			result = append(result, ``)
		}
	}
	return result
}

func getUniqTimeTags(svcsTags map[string]map[string]bool) []string {
	m := make(map[string]bool)
	for _, tags := range svcsTags {
		for tag := range tags {
			m[tag] = true
		}
	}
	slice := make([]string, 0, len(m))
	for key := range m {
		slice = append(slice, key)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(slice)))
	return slice
}

func getTimeTagsMap(svcName, env string) map[string]bool {
	tags := make(map[string]bool)
	for _, tag := range Tags(conf.GetService(svcName, env).ImageName()) {
		if strings.HasPrefix(tag, env) && tag != env {
			tagEnv, timeTag := splitTag(tag)
			if timeTag != `` && tagEnv == env {
				tags[timeTag] = true
			}
		}
	}
	return tags
}

func getEnvTimeTagsMap(svcName, env string) map[string][]string {
	envTagsMap := make(map[string][]string)
	for _, tag := range Tags(conf.GetService(svcName, env).ImageName()) {
		tagEnv, timeTag := splitTag(tag)
		if timeTag == `` || tagEnv != env {
			continue
		}
		envTags := envTagsMap[tagEnv]
		if envTags == nil {
			envTags = []string{timeTag}
		} else {
			envTags = append(envTags, timeTag)
		}
		envTagsMap[tagEnv] = envTags
	}
	for _, envTags := range envTagsMap {
		sort.Sort(sort.Reverse(sort.StringSlice(envTags)))
	}
	return envTagsMap
}

func splitTag(tag string) (string, string) {
	tagLen := len(tag)
	if tagLen < 14 {
		return tag, ``
	}
	var tagEnv, timeTag string
	if strings.Index(tag, `-`) < len(tag)-7 {
		tagEnv, timeTag = tag[:len(tag)-14], tag[len(tag)-13:]
	} else {
		tagEnv, timeTag = tag[:len(tag)-13], tag[len(tag)-13:]
	}
	return tagEnv, timeTag
}
