package tasks

import (
	"github.com/lovego/xiaomei/utils"
	"github.com/robfig/cron"
)

func Run() {
	c := cron.New()
	if err := c.AddFunc("0 * * * * *", hello); err != nil {
		panic(err)
	}
	c.Start()
	select {}
}

func hello() {
	utils.Log(`hello`)
}
