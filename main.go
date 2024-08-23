package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/NicoNex/echotron/v3"
)

// Add setLocationHandler

// Telegram-bot token from env variable
var tg_token = "6519797623:AAHc4Ap3Dadc3vFKcw_Ke9K_refnbGkMdw8"

// Inline keyboard markup
var inlineKeyboard = echotron.InlineKeyboardMarkup{
	InlineKeyboard: [][]echotron.InlineKeyboardButton{
		[]echotron.InlineKeyboardButton{
			echotron.InlineKeyboardButton{
				Text:         "Weather now",
				CallbackData: "now"},
			echotron.InlineKeyboardButton{
				Text:         "Weather forecast",
				CallbackData: "forecast"}}},
}

// Recursive type definition of the bot state function.
type stateFn func(*echotron.Update) stateFn

type bot struct {
	chatID            int64
	state             stateFn
	name              string
	lastUserLocatrion string
	echotron.API
}

// New bot instance
func newBot(chatID int64) echotron.Bot {
	bot := &bot{
		chatID: chatID,
		name:   "WeatherBot",
		API:    echotron.NewAPI(tg_token),
	}
	bot.state = bot.handleMessage
	return bot
}

// Execute the current state and set the next one.
func (b *bot) Update(update *echotron.Update) {
	b.state = b.state(update)
}

func (b *bot) handleMessage(update *echotron.Update) stateFn {
	var lang string
	if update.Message != nil {
		if strings.HasPrefix(update.Message.Text, "/start") {
			b.SendMessage("Hi, i'm weather forecast bot.\nType your city:", b.chatID, nil)
			return b.handleLocation
		} else if strings.HasPrefix(update.Message.Text, "/setlocation") {
			b.SendMessage("Type your city:", b.chatID, nil)
			return b.handleLocation
		} else {
			b.SendMessage("Choose an option:", b.chatID, &echotron.MessageOptions{ReplyMarkup: inlineKeyboard})
			return b.handleMessage
		}
	}
	lang = update.CallbackQuery.From.LanguageCode
	weather, err := WeatherForecast(update.CallbackQuery.Data, lang, b.lastUserLocatrion)
	if err != nil {
		b.SendMessage("Error getting weather for you (((", b.chatID, nil)
		return b.handleMessage
	}
	b.SendMessage(weather, b.chatID, nil)
	return b.handleMessage
}

func (b *bot) handleLocation(update *echotron.Update) stateFn {
	loc := update.Message.Text
	if _, err := WeatherForecast("now", "eng", loc); err != nil {
		b.SendMessage("Invalid location.\nType your city:", b.chatID, nil)
		return b.handleLocation
	}
	b.lastUserLocatrion = loc
	b.SendMessage(fmt.Sprintf("Location %s added successfully", loc), b.chatID, nil)
	return b.handleMessage
}

// Reverse geocoding
func parseLocation(latitude float64, longitude float64) string {
	var location struct {
		City string `json:"city"`
	}
	query := fmt.Sprintf("https://api.bigdatacloud.net/data/reverse-geocode-client?latitude=%v&longitude=%v&localityLanguage=eng", latitude, longitude)
	resp, err := http.Get(query)
	if err != nil {
		log.Println(err)
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return ""
	}
	err = json.Unmarshal(body, &location)
	if err != nil {
		log.Println(err)
		return ""
	}
	return location.City
}

func main() {
	dsp := echotron.NewDispatcher(tg_token, newBot)
	log.Println(dsp.Poll())
}
