package main

import (
	"errors"
	"fmt"

	"github.com/lovego/xiaomei/utils/errs"
)

func main() {
	fmt.Println(testCodeMessage().Error(), "\n")

	fmt.Println(testTrace().Error())
}

func testCodeMessage() error {
	return errs.New(`no-login`, `please login first.`)
}

func testTrace() error {
	err := errors.New(`connection timeout`)
	return errs.Trace(err)
}
