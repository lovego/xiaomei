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
		tags := getTimeTagsSlice(svcName, env)
		fmt.Println(strings.Join(tags, "\n"))
	}
}

func listSvcsTimeTags(env string) {
	svcs := conf.ServiceNames(env)
	svcsTags := make(map[string]map[string]bool)
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
	for key, _ := range m {
		slice = append(slice, key)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(slice)))
	return slice
}

func getTimeTagsMap(svcName, env string) map[string]bool {
	tags := make(map[string]bool)
	for _, tag := range Tags(conf.GetService(svcName, env).ImageName()) {
		if strings.HasPrefix(tag, env) && tag != env {
			tags[tag[len(env):]] = true
		}
	}
	return tags
}

func getTimeTagsSlice(svcName, env string) (tags []string) {
	for _, tag := range Tags(conf.GetService(svcName, env).ImageName()) {
		if strings.HasPrefix(tag, env) && tag != env {
			tags = append(tags, tag[len(env):])
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(tags)))
	return tags
}
