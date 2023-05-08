package main

import (
	config2 "VKtest/pkg/config"
	"VKtest/pkg/db"
	"VKtest/pkg/tools"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
)

var (
	USAGE = "<SERVICE NAME> - Название сервиса\n<LOGIN> - Логин\n<PASSWORD> - Пароль\n------------------------------\nКОМАНДЫ\n------------------------------\n/set <SERVICE NAME> <LOGIN> <PASSWORD>\n/get <SERVICE NAME> <LOGIN>\n/del <SERVICE NAME> <LOGIN>"
	ERROR = "Вы ввели неправильную команду. Возможные команды вы можете узнать по запросу /usage"
)

func main() {
	config := config2.GetConfig()
	bot, err := tgbotapi.NewBotAPI(config.ApiKey)
	fmt.Println(config)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	sql, err := db.GetDb(config)
	if err != nil {
		log.Panic(err)
	}
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		command := strings.Split(update.Message.Text, " ")
		userId := tools.Hash(tools.Hash(strconv.FormatInt(update.Message.From.ID, 10)))
		switch command[0] {
		case "/start":
			msg, err := db.RegisterUser(sql, userId)
			fmt.Println(err)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
				continue
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))

		case "/set":
			if len(command) != 4 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, ERROR))
				continue
			}
			msg, err := db.AddUserData(sql, userId, command[1], command[2], command[3])
			fmt.Println(err)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
				continue
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))

		case "/get":
			if len(command) != 3 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, ERROR))
				continue
			}
			msg, err := db.GetUserData(sql, userId, command[1], command[2])
			fmt.Println(err)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, tools.Prettify(msg)))
				continue
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, tools.Prettify(msg)))
		case "/del":
			if len(command) != 3 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, ERROR))
				continue
			}
			msg, err := db.DeleteUserData(sql, userId, command[1], command[2])
			fmt.Println(err)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
				continue
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))

		case "/usage":
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, USAGE))
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, ERROR))

		}
	}
}
