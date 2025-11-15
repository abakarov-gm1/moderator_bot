package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

func isAllowedTime(t time.Time) bool {
	h := t.Hour()

	inMorning := h >= 10 && h < 12
	inEvening := h >= 19 && h < 22

	return inMorning || inEvening
}

func main() {
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Bot authorized: %s", bot.Self.UserName)

	chatID := int64() // ID вашего группового чата

	ticker := time.NewTicker(30 * time.Second) // каждые 30 сек проверяем время
	defer ticker.Stop()

	var lastState bool = true // чтобы не дергать Telegram лишний раз
	// тут при запуске бота нужно менять на true or false
	// в зависимости от того allowed
	// если allowed = true то lastState = false а если allowed = false то lastState = true

	for {
		<-ticker.C
		now := time.Now()
		allowed := isAllowedTime(now)

		fmt.Println(allowed, "alowedddd", lastState)

		if allowed == lastState {
			continue
		}
		lastState = allowed

		if allowed {
			// Разрешить писать
			perms := tgbotapi.ChatPermissions{
				CanSendMessages:       true,
				CanSendMediaMessages:  true,
				CanSendPolls:          true,
				CanSendOtherMessages:  true,
				CanAddWebPagePreviews: true,
			}
			cfg := tgbotapi.SetChatPermissionsConfig{
				ChatConfig: tgbotapi.ChatConfig{
					ChatID: chatID, // обязательно с большой буквы!
				},
				Permissions: &perms,
			}
			bot.Send(cfg)

			msg := tgbotapi.NewMessage(chatID, "Чат открыт. \nВремя для сообщений: 10:00–12:00 и 19:00–21:00")
			bot.Send(msg)
		} else {
			// Запретить писать
			perms := tgbotapi.ChatPermissions{
				CanSendMessages: false,
			}

			cfg := tgbotapi.SetChatPermissionsConfig{
				ChatConfig: tgbotapi.ChatConfig{
					ChatID: chatID,
				},
				Permissions: &perms,
			}
			bot.Send(cfg)

			msg := tgbotapi.NewMessage(chatID, "Чат на паузе!\nСообщения принимаются только:\n• 10:00–12:00\n• 19:00–22:00\nОстальное время — отдых для мозга")
			bot.Send(msg)
		}
	}
}
