package main

import (
	"database/sql"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func receiveStart(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) (err error) {
	prev, err := checkPrivilege(db, update.Message.Chat.ID)
	if err != nil {
		return err
	}
	mes := ""
	if prev != "Unauthorized" {
		mes = "Welcome to VTgoBot, "
		mes += update.Message.From.UserName
		mes += "! Write /help to see available commands"
	} else {
		mes = "What a nice weather today, isn`t it?"
	}
	_, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, mes))
	clientId, dbErr := getClientId(db, update.Message.Chat.ID)
	if dbErr != nil {
		return dbErr
	}
	if err != nil {
		dbErr = writeLog(db, clientId, update.Message.Text, "", "", mes, 1)
		if dbErr != nil {
			return dbErr
		}
		return err
	}
	err = writeLog(db, clientId, update.Message.Text, "", "", mes, 0)

	return err
}

func receiveHelp(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) (err error) {
	prev, err := checkPrivilege(db, update.Message.Chat.ID)
	if err != nil {
		return err
	}
	mes := ""
	if prev != "Unauthorized" {
		mes = "Write to me /imo or /mmsi with the desired numbers (even separated by comma), like /imo 12345678,98765432"
	} else {
		mes = "Have you seen last season of Invincible series?"
	}
	_, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, mes))
	clientId, dbErr := getClientId(db, update.Message.Chat.ID)
	if dbErr != nil {
		return dbErr
	}
	if err != nil {
		dbErr = writeLog(db, clientId, update.Message.Text, "", "", mes, 1)
		if dbErr != nil {
			return dbErr
		}
		return err
	}
	err = writeLog(db, clientId, update.Message.Text, "", "", mes, 0)

	return err
}

func receiveAuth(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) (err error) {
	prev, err := checkPrivilege(db, update.Message.Chat.ID)
	if err != nil {
		return err
	}
	mes := ""
	if prev != "Owner" {
		mes = "Fanfare For The Common Man or Make America Great Again?"
	} else {
		splits := strings.Split(update.Message.Text, " ")
		var dbErr error
		if len(splits) != 2 {
			mes = "Wrong arguments for command"
			dbErr = errors.New("wrong arguments for command")
		} else {
			var tgId int64
			tgId, dbErr = strconv.ParseInt(splits[1], 10, 64)
			if dbErr != nil {
				mes = "Failed to get tgId argument"
			} else {
				dbErr = authorizeClient(db, tgId)
				if dbErr != nil {
					mes = "Failed to authorize client"
				}
			}
		}
		if dbErr != nil {
			clientId, logErr := getClientId(db, update.Message.Chat.ID)
			if logErr != nil {
				return logErr
			}
			logErr = writeLog(db, clientId, update.Message.Text, "", "", mes, 1)
			if logErr != nil {
				return logErr
			}
		} else {
			mes = "User authorized!"
		}
	}
	clientId, logErr := getClientId(db, update.Message.Chat.ID)
	if logErr != nil {
		return logErr
	}
	logErr = writeLog(db, clientId, update.Message.Text, "", "", mes, 0)
	if logErr != nil {
		return logErr
	}
	_, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, mes))
	return err
}

func receiveDeauth(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) (err error) {
	prev, err := checkPrivilege(db, update.Message.Chat.ID)
	if err != nil {
		return err
	}
	mes := ""
	if prev != "Owner" {
		mes = "Did you know that more than a quarter of the population of Germany is of migration background?"
	} else {
		splits := strings.Split(update.Message.Text, " ")
		var dbErr error
		if len(splits) != 2 {
			mes = "Wrong arguments for command"
			dbErr = errors.New("wrong arguments for command")
		} else {
			var tgId int64
			tgId, dbErr = strconv.ParseInt(splits[1], 10, 64)
			if dbErr != nil {
				mes = "Failed to get tgId argument"
			} else {
				dbErr = deauthorizeClient(db, tgId)
				if dbErr != nil {
					mes = "Failed to deauthorize client"
				}
			}
		}
		if dbErr != nil {
			clientId, logErr := getClientId(db, update.Message.Chat.ID)
			if logErr != nil {
				return logErr
			}
			logErr = writeLog(db, clientId, update.Message.Text, "", "", mes, 1)
			if logErr != nil {
				return logErr
			}
		} else {
			mes = "User deauthorized!"
		}
	}
	clientId, logErr := getClientId(db, update.Message.Chat.ID)
	if logErr != nil {
		return logErr
	}
	logErr = writeLog(db, clientId, update.Message.Text, "", "", mes, 0)
	if logErr != nil {
		return logErr
	}
	_, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, mes))
	return err
}

func receiveImo(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) (err error) {
	prev, err := checkPrivilege(db, update.Message.Chat.ID)
	if err != nil {
		return err
	}
	mes := ""
	request := ""
	response := ""
	if prev == "Unauthorized" {
		mes = "The Emperor protects!"
	} else {
		splits := strings.Split(update.Message.Text, " ")
		var dbErr error
		if len(splits) != 2 {
			mes = "Wrong arguments for command"
			dbErr = errors.New("wrong arguments for command")
		} else {
			url := "https://api.vtexplorer.com/vessels?userkey="
			url += "WS-096EE673-456A8B"
			url += "&sat=1"
			imo := splits[1]
			url += "&imo="
			url += imo

			request = url

			req, _ := http.NewRequest("GET", url, nil)
			res, _ := http.DefaultClient.Do(req)
			defer res.Body.Close()
			body, _ := io.ReadAll(res.Body)
			response = string(body)
			mes, dbErr = processJson(body, bot, update)
			if mes == "" {
				dbErr = errors.New("no data received from server")
			}
		}
		if dbErr != nil {
			mes = "Failed to receive data"
			clientId, logErr := getClientId(db, update.Message.Chat.ID)
			if logErr != nil {
				return logErr
			}
			logErr = writeLog(db, clientId, update.Message.Text, request, response, mes, 1)
			if logErr != nil {
				return logErr
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, mes))
			return dbErr
		}
	}
	clientId, logErr := getClientId(db, update.Message.Chat.ID)
	if logErr != nil {
		return logErr
	}
	logErr = writeLog(db, clientId, update.Message.Text, request, response, mes, 0)
	if logErr != nil {
		return logErr
	}
	return err
}

func receiveMmsi(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) (err error) {
	prev, err := checkPrivilege(db, update.Message.Chat.ID)
	if err != nil {
		return err
	}
	mes := ""
	request := ""
	response := ""
	if prev == "Unauthorized" {
		mes = "All we have to decide is what to do with the time that is given us. (c) Gandalf"
	} else {
		splits := strings.Split(update.Message.Text, " ")
		var dbErr error
		if len(splits) != 2 {
			mes = "Wrong arguments for command"
			dbErr = errors.New("wrong arguments for command")
		} else {
			url := "https://api.vtexplorer.com/vessels?userkey="
			url += "WS-096EE673-456A8B"
			url += "&sat=1"
			mmsi := splits[1]
			url += "&mmsi="
			url += mmsi

			request = url

			req, _ := http.NewRequest("GET", url, nil)
			res, _ := http.DefaultClient.Do(req)
			defer res.Body.Close()
			body, _ := io.ReadAll(res.Body)
			response = string(body)
			mes, dbErr = processJson(body, bot, update)
			if mes == "" {
				dbErr = errors.New("no data received from server")
			}
		}
		if dbErr != nil {
			mes = "Failed to receive data"
			clientId, logErr := getClientId(db, update.Message.Chat.ID)
			if logErr != nil {
				return logErr
			}
			logErr = writeLog(db, clientId, update.Message.Text, request, response, mes, 1)
			if logErr != nil {
				return logErr
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, mes))
			return dbErr
		}
	}
	clientId, logErr := getClientId(db, update.Message.Chat.ID)
	if logErr != nil {
		return logErr
	}
	logErr = writeLog(db, clientId, update.Message.Text, request, response, mes, 0)
	if logErr != nil {
		return logErr
	}
	return err
}

func sendBroadcast(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) (err error) {
	prev, err := checkPrivilege(db, update.Message.Chat.ID)
	if err != nil {
		return err
	}
	mes := ""
	if prev != "Owner" {
		mes = "https://youtu.be/dQw4w9WgXcQ?si=effhCQuDyfI7v1XJ"
	} else {
		broadMes, found := strings.CutPrefix(update.Message.Text, "/broadcast ")
		if !found {
			mes = "Wrong arguments for command"
			err = errors.New("wrong arguments for command")
		} else {
			clients := make([]int64, 0, 0)
			clients, err = getAuthorizedClients(db)
			if err != nil {
				mes = "Failed to get clients from database"
			} else {
				for _, client := range clients {
					bot.Send(tgbotapi.NewMessage(client, broadMes))
				}
			}
		}
	}
	clientId, logErr := getClientId(db, update.Message.Chat.ID)
	if logErr != nil {
		return logErr
	}
	result := 0
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, mes))
		result = 1
	}
	logErr = writeLog(db, clientId, update.Message.Text, "", "", mes, result)
	if logErr != nil {
		return logErr
	}
	return err

}

func receiveDefault(update tgbotapi.Update, bot *tgbotapi.BotAPI, db *sql.DB) (err error) {
	prev, err := checkPrivilege(db, update.Message.Chat.ID)
	if err != nil {
		return err
	}
	mes := ""
	if prev != "Unauthorized" {
		mes = "Unknown command!"
	} else {
		mes = "Through discipline comes freedom! (c) Aristotle"
	}

	_, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, mes))
	clientId, dbErr := getClientId(db, update.Message.Chat.ID)
	if dbErr != nil {
		return dbErr
	}
	if err != nil {
		dbErr = writeLog(db, clientId, update.Message.Text, "", "", mes, 1)
		if dbErr != nil {
			return dbErr
		}
		return err
	}
	err = writeLog(db, clientId, update.Message.Text, "", "", mes, 0)

	return err
}
