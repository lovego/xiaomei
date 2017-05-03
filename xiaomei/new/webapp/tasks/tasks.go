package tasks

import (
	"github.com/fatih/color"
	"github.com/lovego/xiaomei/utils"
	"github.com/robfig/cron"
)

func Run() {
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
