package access

import "text/template"

var funcsMap = template.FuncMap{`domainAncestor`: DomainAncestor}

// DomainAncestor return the N'th ancestor of domain
func DomainAncestor(domain string, n int) string {
	if n <= 0 {
		return domain
	}
	var index int
	for i, b := range domain {
		if b == '.' {
			index = i
			n--
			if n == 0 {
				break
			}
		}
	}
	return domain[index+1:]
}
