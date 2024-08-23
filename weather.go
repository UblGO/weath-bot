package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

//complete weather struct and queries

type weatherStruct struct {
	Location struct {
		Country   string `json:"country"`
		City      string `json:"name"`
		Localtime string `json:"localtime"`
	}
	Current struct {
		Last_updated string  `json:"last_updated"`
		Temp         float64 `json:"temp_c"`
		Condition    struct {
			Text string `json:"text"`
		}
		Wind      float64 `json:"wind_kph"`
		Wind_dir  string  `json:"wind_dir"`
		Cloud     int     `json:"cloud"`
		Feelslike float64 `json:"feelslike_c"`
		Gust      float64 `json:"gust_kph"`
	}
	Forecast struct {
		Forecastday []struct {
			Date string `json:"date,omitempty"`
			Day  struct {
				Maxtemp    float64 `json:"maxtemp_c,omitempty"`
				Mintemp    float64 `json:"mintemp_c,omitempty"`
				Avgtemp    float64 `json:"avgtemp_c,omitempty"`
				Maxwind    float64 `json:"maxwind_kph,omitempty"`
				Rainchance int     `json:"daily_chance_of_rain,omitempty"`
				Snowchance int     `json:"daily_chance_of_snow,omitempty"`
				Condition  struct {
					Text string `json:"text,omitempty"`
				}
			}
		}
	}
}

var weather_token = "c26b8244007a46dfafb152118241504"

const (
	currWeatherForamt     = "Регион: %s, %s\nМестное время: %s\nПоследнее обновление: %s\nТемпература: %.1f°; Ощущается как: %.1f°\nПогода: %s\nВетер: %.1f; направление: %s\nОблачность: %d%%\nПорывы ветра(км.ч): %.1f\n"
	forecastWeatherFormat = "Регион: %s, %s\n"
	dayForecastFormat     = "Дата: %s\nПогода: %s\nМинимальная темература: %v°\nМаксимальная температура: %v°\nСредняя температура: %v°\nМаксимальная скорость ветра: %v км/ч\nШанс дождя: %d%%\nШанс снега: %d%%\n\n"
)

func WeatherForecast(code string, lang string, location string) (string, error) {
	var query string
	switch code {
	case "now":
		query = fmt.Sprintf(
			"https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no&lang=%s",
			weather_token, location, lang)
	case "forecast":
		query = fmt.Sprintf(
			"https://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&days=3&aqi=no&lang=%s",
			weather_token, location, lang)
	}
	return getQuery(query, code)
}

func getQuery(query string, code string) (string, error) {
	var weather weatherStruct
	var result string

	response, err := http.Get(query)
	if err != nil {
		log.Println(err)
		return "", err
	}
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		log.Println(err)
		return "", err
	}
	err = json.Unmarshal(body, &weather)
	if err != nil {
		log.Println("Unmarshaling error: ", err)
	}
	if code == "now" {
		result = fmt.Sprintf(currWeatherForamt, weather.Location.Country, weather.Location.City, weather.Location.Localtime, weather.Current.Last_updated, weather.Current.Temp, weather.Current.Feelslike, weather.Current.Condition.Text,
			weather.Current.Wind, weather.Current.Wind_dir, weather.Current.Cloud, weather.Current.Gust)
	} else {
		result = fmt.Sprintf(forecastWeatherFormat, weather.Location.Country, weather.Location.City)
		for _, v := range weather.Forecast.Forecastday {
			result += fmt.Sprintf(dayForecastFormat, v.Date, v.Day.Condition.Text, v.Day.Mintemp, v.Day.Maxtemp, v.Day.Avgtemp, v.Day.Maxwind, v.Day.Rainchance, v.Day.Snowchance)
		}
	}
	return result, nil
}
