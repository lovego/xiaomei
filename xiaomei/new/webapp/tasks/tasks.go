package main

import (
	"log"

	"github.com/fatih/color"
	"github.com/robfig/cron"
)

func main() {
	c := cron.New()
	if err := c.AddFunc("0 * * * * *", hello); err != nil {
		panic(err)
	}
	c.Start()
	log.Println(color.GreenString(`started.`))
	select {}
}

func hello() {
	log.Println(`hello`)
}
