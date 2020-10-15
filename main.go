package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var allUpdates []tgbotapi.Update
var chatIds []int

func main() {
	fmt.Printf("%v", 22)
	bot, err := tgbotapi.NewBotAPI("1221394093:AAGGvqYt7RvqFdFzlZ4h98AzIRLEpQvJZO0")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.ChannelPost != nil {
			//if update.ChannelPost.Document != nil {
			allUpdates = append(allUpdates, update)
			//}
		}

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		fmt.Println(allUpdates)

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "delete":
				count, err := deleteMessages(bot)
				if count == 0 {
					msg.Text = "No files to delete"
					break
				}
				if err == nil {
					msg.Text = "I have deleted all the files"
					break
				}
				msg.Text = err.Error()
			case "dfilecaption":
				count, err := deleteMessagesWithCaption(bot)
				if count == 0 {
					msg.Text = "No files to delete"
					break
				}
				if err == nil {
					msg.Text = "I have deleted all the files"
					break
				}
				msg.Text = err.Error()

			default:
				msg.Text = "I don't know that command"
			}
			bot.Send(msg)
		}
	}
}

type Config struct {
	Bot              string   `json:"bot"`
	ChatID           []int64  `json:"ChatId"`
	Str              []string `json:"str"`
	SendTimeoutInSec int      `json:"sendTimeoutInSec"`
}

func deleteMessages(tgapi *tgbotapi.BotAPI) (int, error) {
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened config.json")
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)
	var config Config

	json.Unmarshal(byteValue, &config)
	fmt.Println(config)

	count := 0
	for _, key := range config.Str {
		if len(allUpdates) != 0 {

			for _, files := range allUpdates {
				if files.ChannelPost != nil {
					if files.ChannelPost.Document != nil {
						if strings.HasPrefix(files.ChannelPost.Document.FileName, key) {
							cfg := tgbotapi.DeleteMessageConfig{
								ChatID:    files.ChannelPost.Chat.ID,
								MessageID: files.ChannelPost.MessageID,
							}
							resp, err := tgapi.DeleteMessage(cfg)
							if err != nil {
								return 0, err
							}
							fmt.Println(resp.Description)
							chatIds = append(chatIds, files.ChannelPost.MessageID)
							count++
						}
					}

				}
			}
			for _, id := range chatIds {
				index := pos(allUpdates, id)
				fmt.Println(index)
				if index != -1 {
					allUpdates = RemoveIndex(allUpdates, index)
				}
			}
			chatIds = nil
		}
	}
	return count, nil
}

func deleteMessagesWithCaption(tgapi *tgbotapi.BotAPI) (int, error) {
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened config.json")
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)
	var config Config

	json.Unmarshal(byteValue, &config)
	fmt.Println(config)

	count := 0
	for _, key := range config.Str {
		if len(allUpdates) != 0 {

			for _, files := range allUpdates {
				if files.ChannelPost != nil {
					if files.ChannelPost.Caption != "" {
						if strings.HasPrefix(files.ChannelPost.Caption, key) {
							cfg := tgbotapi.DeleteMessageConfig{
								ChatID:    files.ChannelPost.Chat.ID,
								MessageID: files.ChannelPost.MessageID,
							}
							resp, err := tgapi.DeleteMessage(cfg)
							if err != nil {
								return 0, err
							}
							fmt.Println(resp.Description)
							chatIds = append(chatIds, files.ChannelPost.MessageID)
							count++
						}
					}
				}
			}
			for _, id := range chatIds {
				index := pos(allUpdates, id)
				if index != -1 {
					allUpdates = RemoveIndex(allUpdates, index)
				}
			}
			chatIds = nil
		}
	}
	return count, nil
}

func pos(slice []tgbotapi.Update, value int) int {
	for p, v := range slice {
		if v.ChannelPost != nil {
			if v.ChannelPost.MessageID == value {
				return p
			}
		}
	}
	return -1
}

func RemoveIndex(s []tgbotapi.Update, index int) []tgbotapi.Update {
	return append(s[:index], s[index+1:]...)
}
