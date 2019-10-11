package access

import "fmt"

func ExampleDomainAncestor() {
	fmt.Println(DomainAncestor(`a.b.c.com`, 0))
	fmt.Println(DomainAncestor(`a.b.c.com`, 1))
	fmt.Println(DomainAncestor(`a.b.c.com`, 2))
	fmt.Println(DomainAncestor(`a.b.c.com`, 3))
	fmt.Println(DomainAncestor(`a.b.c.com`, 4))

	// Output:
	// a.b.c.com
	// b.c.com
	// c.com
	// com
	// com
}
