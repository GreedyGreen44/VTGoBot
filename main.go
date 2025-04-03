package main

import (
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

func main() {
	dbPath := os.Getenv("VTGOBOT_DB_PATH")
	db, err := createConnection(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer closeConnection(db)

	if db.Ping() != nil {
		log.Fatal(err)
	}

	botToken := os.Getenv("VTGOBOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to connect to telegram bot: %v", err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		go func(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) {
			if update.Message != nil {
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

				err = addNewClient(db, update.Message.Chat.ID, update.Message.From.UserName, update.Message.From.FirstName, update.Message.From.LastName)
				if err != nil {
					if err != existingClientError {
						log.Printf("Error adding new client: %v", err)
						return
					} else {
						err = setClientNames(db, update.Message.Chat.ID, update.Message.From.UserName, update.Message.From.FirstName, update.Message.From.LastName)
						if err != nil {
							log.Printf("Error setting client names: %v", err)
						}
					}
				}

				switch update.Message.Command() {
				case "start":
					err = receiveStart(update, bot, db)
					if err != nil {
						log.Printf("Error processing /start command: %v", err)
					}
				case "help":
					err = receiveHelp(update, bot, db)
					if err != nil {
						log.Printf("Error processing /help command: %v", err)
					}
				case "auth":
					err = receiveAuth(update, bot, db)
					if err != nil {
						log.Printf("Error processing /auth command: %v", err)
					}
				case "deauth":
					err = receiveDeauth(update, bot, db)
					if err != nil {
						log.Printf("Error processing /deauth command: %v", err)
					}
				case "imo":
					err = receiveImo(update, bot, db)
					if err != nil {
						log.Printf("Error processing /imo command: %v", err)
					}
				case "mmsi":
					err = receiveMmsi(update, bot, db)
					if err != nil {
						log.Printf("Error processing /mmsi command: %v", err)
					}
				case "broadcast":
					err = sendBroadcast(update, bot, db)
					if err != nil {
						log.Printf("Error processing /broadcast commsnd: %v", err)
					}
				default:
					err = receiveDefault(update, bot, db)
					if err != nil {
						log.Printf("Error processing unknown command: %v", err)
					}
				}
			}
		}(update, bot, db)
	}
}
