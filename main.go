package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type food map[string]float64

var db = map[int64]food{}

var count int = 0

func infoOutput(count int) string {
	if count == 1 {
		return "Добавьте второе"
	}
	if count == 2 {
		return "Добавьте апитайзер"
	}
	if count == 3 {
		return "Добавьте напиток"
	}
	return ""
}

func main() {
	bot, err := tgbotapi.NewBotAPI("5435714187:AAGpfH4chr8fiXOpRnZTuClRNrZMgRR1ZxM")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		command := strings.Split(update.Message.Text, " ")

		switch command[0] {

		case "/start":
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "ADD *Название блюда* *Цена* - Добавление первого, второго, апитайзер(салат, холдные/горячие закуски), напиток в меню\nDELETE *Название блюда* - Удаление чего-то из меню\nSHOW - Показать меню"))
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Добавьте первое"))
		case "ADD":
			count++
			if len(command) != 3 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверная команда"))
				continue
			}

			amount, err := strconv.ParseFloat(command[2], 64)

			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
			}

			if _, ok := db[update.Message.Chat.ID]; !ok {
				db[update.Message.Chat.ID] = food{}
			}

			db[update.Message.Chat.ID][command[1]] += amount

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, command[1]+" успешно добавлен в меню"))

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, infoOutput(count)))

		case "DELETE":
			if len(command) != 2 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверная команда"))
			}

			delete(db[update.Message.Chat.ID], command[1])

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, command[1]+" успешно удален из меню"))

			//command[1] = "Неизвестно"
			//db[update.Message.Chat.ID][command[1]] = 0

		case "SHOW":
			msg := ""
			for key, value := range db[update.Message.Chat.ID] {
				msg += fmt.Sprintf("%s: %s\n", key, strconv.FormatFloat(value, 'f', -1, 64)+" тг")
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда"))
		}

	}
}
