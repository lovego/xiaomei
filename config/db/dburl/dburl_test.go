package dburl

import "fmt"

func ExampleParse() {
	u := Parse(
		"postgres://develop:@localhost/goods_orders_dev?sslmode=disable",
	)
	fmt.Println(u.URL)
	fmt.Println(u.MaxOpen, u.MaxIdle, u.MaxLife)
	// Output:
	// postgres://develop:@localhost/goods_orders_dev?sslmode=disable
	// 10 0 10m0s
}

func ExampleParse_withParams() {
	u := Parse(
		"postgres://develop:@localhost/goods_orders_dev" +
			"?sslmode=disable&maxIdle=1&maxOpen=20&maxLife=1h",
	)
	fmt.Println(u.URL)
	fmt.Println(u.MaxOpen, u.MaxIdle, u.MaxLife)
	// Output:
	// postgres://develop:@localhost/goods_orders_dev?sslmode=disable
	// 20 1 1h0m0s
}
