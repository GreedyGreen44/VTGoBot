package main

import (
	"bytes"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func processJson(body []byte, bot *tgbotapi.BotAPI, update tgbotapi.Update) (combinedAnswer string, err error) {
	body = bytes.Replace(body, []byte("[{\"AIS\":{"), []byte("{\"AIS\":[{"), -1)
	body = bytes.Replace(body, []byte("}},{\"AIS\":{"), []byte("},{"), -1)
	body = bytes.Replace(body, []byte("}}]"), []byte("}]}"), -1)

	var resp VtResp
	var answer string
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", err
	}
	for _, ves := range resp.AisDataset {
		answer = "\nTime: " + ves.Timestamp
		answer += "\nName: " + ves.Name
		answer += "\nMMSI: " + strconv.FormatInt(ves.Mmsi, 10)
		answer += "\nIMO: " + strconv.FormatInt(ves.Imo, 10)
		answer += "\nCallsign: " + ves.Callsign
		answer += "\nLat: " + strconv.FormatFloat(float64(ves.Latitude), 'f', -1, 32)
		answer += "\nLon: " + strconv.FormatFloat(float64(ves.Longitude), 'f', -1, 32)
		answer += "\nSpeed: " + strconv.FormatFloat(float64(ves.Speed), 'f', -1, 32)
		answer += "\nHeading: " + strconv.FormatInt(int64(ves.Heading), 10)
		answer += "\nShip type: " + decodeShipType(ves.Type)
		answer += "\nNavigation Status: " + decodeNavSat(ves.Navstat)
		answer += "\nA: " + strconv.FormatInt(int64(ves.A), 10)
		answer += "\nB: " + strconv.FormatInt(int64(ves.B), 10)
		answer += "\nC: " + strconv.FormatInt(int64(ves.C), 10)
		answer += "\nD: " + strconv.FormatInt(int64(ves.D), 10)
		answer += "\nDrought: " + strconv.FormatFloat(float64(ves.Draught), 'f', -1, 32)
		answer += "\nDestination: " + ves.Destination
		answer += "\nLOCODE: " + ves.Locode
		answer += "\nETA_AIS: " + ves.ETA_Ais
		answer += "\nETA: " + ves.ETA
		answer += "\nSrc: " + ves.Src
		answer += "\nZone: " + ves.Zone
		answer += "\nECA: " + strconv.FormatBool(ves.Eca)
		answer += "\nDistance remaining: " + strconv.FormatInt(int64(ves.Distance_remaining), 10)
		answer += "\nETA_Predicted: " + ves.ETA_predicted
		answer += "\n----------------------------------------------------------"

		combinedAnswer += answer

		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, answer))
		bot.Send(tgbotapi.NewLocation(update.Message.Chat.ID, float64(ves.Latitude), float64(ves.Longitude)))
	}

	return combinedAnswer, nil
}
