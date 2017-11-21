package main

import (
	"github.com/fatih/color"
	"github.com/lovego/utils"
	"github.com/robfig/cron"
)

func main() {
	c := cron.New()
	if err := c.AddFunc("0 * * * * *", hello); err != nil {
		panic(err)
	}
	c.Start()
	utils.Log(color.GreenString(`started.`))
	select {}
}

func hello() {
	utils.Log(`hello`)
}
