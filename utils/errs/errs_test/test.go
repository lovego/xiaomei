package main

import (
	"errors"
	"fmt"

	"github.com/lovego/xiaomei/utils/errs"
)

func main() {
	fmt.Println(testCodeMessage().Error(), "\n")

	fmt.Println(testStack().Error())
}

func testCodeMessage() error {
	return errs.New(`no-login`, `please login first.`)
}

func testStack() error {
	err := errors.New(`connection timeout`)
	return errs.Stack(err)
}
