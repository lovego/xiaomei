package tasks

import (
	"fmt"
	"regexp"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils"
	"strings"
	"time"
)

type weatherInfo struct {
	Date        string `json:"date"`
	Week        string `json:"week"`
	Weather     string `json:"weather"`
	Wind        string `json:"wind"`
	Temperature string `json:"temperature"`
}

type baiduWeather struct {
	Error   float64
	Status  string
	Date    string
	Results interface{}
}

func GetUpdateWeather() {
	basic_url := `http://api.map.baidu.com/telematics/v3/weather?output=json&ak=FK9mkfdQsloEngodbFl4FeY3`
	data := baiduWeather{}
	utils.HttpGet(basic_url+`&location=北京`, &data)
	if data.Error != 0 {
		fmt.Println(`get baidu weather faild`, data.Status)
	}
	weather_data := data.Results.([]interface{})[0].(map[string]interface{})[`weather_data`].([]interface{})
	today, _ := time.ParseInLocation(`2006-01-02`, data.Date, config.TimeZone)
	weathers := make([]weatherInfo, len(weather_data))
	for i, v := range weather_data {
		weather := v.(map[string]interface{})
		week := weather[`date`].(string)
		if len(week) > 6 {
			week = strings.Split(week, ` `)[0]
		}
		d := weatherInfo{
			today.Format(`20060102`),
			week,
			weather[`weather`].(string),
			weather[`wind`].(string),
			processTemperature(weather[`temperature`].(string)),
		}
		weathers[i] = d
		today = today.AddDate(0, 0, 1)
	}
	success := updateWeather(weathers)
	if success {
		fmt.Println(`update weather ok`)
	}
}

func updateWeather(weathers []weatherInfo) bool {
	preIns, err := config.Mysql().Prepare(`insert into weathers(dt, week, weather, wind, temperature) values (?, ?, ?, ?, ?)
        on duplicate key update weather = ?, wind = ?, temperature = ?`)
	var msg string
	for _, weather := range weathers {
		_, err = preIns.Exec(weather.Date, weather.Week, weather.Weather, weather.Wind, weather.Temperature,
			weather.Weather, weather.Wind, weather.Temperature)
		if err != nil {
			msg += fmt.Sprintf(`update weather faild: %v, because: %s`, weather, err)
		}
	}
	if msg != `` {
		config.AlarmMail(`更新weather错误`, msg)
		fmt.Println(msg)
		return false
	}
	return true
}

// t: 33 ~ 23℃, result: 23~33℃
func processTemperature(t string) string {
	reg := regexp.MustCompile(`\d+`)
	temperature := reg.FindAllString(t, -1)
	return fmt.Sprintf(`%s~%s℃`, temperature[1], temperature[0])
}
