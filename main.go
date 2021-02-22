/*
Название бота: BrnRGS (@BryanskRGS_bot)
Комманды бота:
/start - начало работы с ботом
/help - список команд
/about - краткая информация о боте
/погода - запрос погоды по названию города на 5 дней/3 часа
/афоризмы
/phone - запрос внутреннего номера по фамилии
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	_ "os"
	"strconv"
	"strings"
	"time"
)

const telegramBaseUrl = "https://api.telegram.org/bot"
const telegramToken = "!!!telegram-token!!!"
const methodGetMe = "getMe"
const methodGetUpdates = "getUpdates"
const methodSendMessage = "sendMessage"
const botName = "@BryanskRGS_bot"

var ticker = time.NewTicker(5 * time.Second)

type GetMeT struct {
	Ok     bool        `json:"ok"`
	Result GetMeResult `json:"result"`
}
type GetMeResult struct {
	Id        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}
type SendMessageT struct {
	Ok     bool     `json:"ok"`
	Result MessageT `json:"result"`
}
type MessageT struct {
	MessageID int                          `json:"message_id"`
	From      GetUpdatesResultMessageFromT `json:"from"`
	Chat      GetUpdatesResultMessageChatT `json:"chat"`
	Date      int                          `json:"date"`
	Text      string                       `json:"text"`
}
type GetUpdatesResultMessageFromT struct {
	Id           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}
type GetUpdatesResultMessageChatT struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}
type GetUpdatesT struct {
	Ok     bool                `json:"ok"`
	Result []GetUpdatedResultT `json:"result"`
}
type GetUpdatedResultT struct {
	UpdateID int                `json:"update_id"`
	Message  GetUpdatesMessageT `json:message,omitempty"`
}
type GetUpdatesMessageT struct {
	MessageID int `json:"message_id"`
	From      struct {
		ID           int    `json:"id"`
		IsBot        bool   `json:"is_bot"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Username     string `json:"username"`
		LanguageCode string `json:"language_code"`
	} `json:"from"`
	Chat struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
		Type      string `json:"type"`
	} `json:"chat"`
	ReplyToMessage struct {
		Chat struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Type      string `json:"type"`
			Username  string `json:"username"`
		} `json:"chat"`
		Date int64 `json:"date"`
		From struct {
			FirstName string `json:"first_name"`
			ID        int64  `json:"id"`
			Username  string `json:"username"`
		} `json:"from"`
		MessageID int64  `json:"message_id"`
		Text      string `json:"text"`
	} `json:"reply_to_message"`
	Date int    `json:"date"`
	Text string `json:"text"`
}

func getUrlByMethod(methodName string, offset int) string {
	if methodName == "getUpdates" {
		return telegramBaseUrl + telegramToken + "/" + methodName + "?offset=" + strconv.Itoa(offset)
	} else {
		return telegramBaseUrl + telegramToken + "/" + methodName
	}
}
func getBodyByUrl(url string) []byte {
	// proxyUrl, err := url.Parse("92.255.252.44:4145")
	// myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	// response, err := myClient.Get(url2)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err.Error())
	}
	return body
}
func main() {
	helpText := map[string]string{
		"start":    "начать работу с ботом. Отправьте <b>/help</b> боту для вызова справки",
		"about":    "Версия бота 1.1 (go). Бот для внутреннего использования в Филиале",
		"contacts": "Author - <b>Manzulin Maksim</b> (@maksiz32)",
		"help": "краткое описание команд:\n\n" +
			"<b>/start</b> - начать работу с ботом\n" +
			"<b>/about</b> - краткая информация о боте\n" +
			"<b>/contacts</b> - информация об авторе\n" +
			"<b><u>Ожидаемые доработки функционала:</u></b>\n" +
			"\t<b>/phone</b> - запрос внутреннего номера телефона по фамилии сотрудника\n" +
			"\t!Автоматическое оповещение в конкретную группу о заявке на автомобиль!",
		"undefined": "Неизвестная команда. Отправьте <b>/help</b> боту для вызова справки",
		"noargs":    "Не введено ни одной команды. Отправьте <b>/help</b> боту для вызова справки",
	}
	updateID := -1
	getUpdates := GetUpdatesT{}
	for t := range ticker.C {
		fmt.Println("Запрос в", t)
		body := getBodyByUrl(getUrlByMethod(methodGetUpdates, updateID))
		if err := json.Unmarshal(body, &getUpdates); err != nil {
			fmt.Printf("Error in Unmarshal getUpdates: %s", err.Error())
			continue
		}
		sendMessageUrl := getUrlByMethod(methodSendMessage, updateID)
		for _, item := range getUpdates.Result {
			if len(item.Message.Text) > 0 {
				args := strings.Fields(strings.ToLower(item.Message.Text))
				workArg := args[0]
				if item.Message.Chat.ID < 0 {
					if len(args) > 1 && args[len(args)-1] == strings.ToLower("@BryanskRGS_bot") {
						workArg = args[0]
					} else {
						fmt.Println(args)
						workArg = "falseAnswer"
					}
				}
				updateID = item.UpdateID + 1
				switch workArg {
				case "/start":
					getBodyByUrl(sendMessageUrl + "?chat_id=" + strconv.Itoa(item.Message.Chat.ID) + "&text=" + url.QueryEscape(helpText["start"]) + "&parse_mode=HTML")
				case "/about":
					getBodyByUrl(sendMessageUrl + "?chat_id=" + strconv.Itoa(item.Message.Chat.ID) + "&text=" + url.QueryEscape(helpText["about"]) + "&parse_mode=HTML")
				case "/help":
					getBodyByUrl(sendMessageUrl + "?chat_id=" + strconv.Itoa(item.Message.Chat.ID) + "&text=" + url.QueryEscape(helpText["help"]) + "&parse_mode=HTML")
				case "/contacts":
					getBodyByUrl(sendMessageUrl + "?chat_id=" + strconv.Itoa(item.Message.Chat.ID) + "&text=" + url.QueryEscape(helpText["contacts"]) + "&parse_mode=HTML")
				case "falseAnswer":
					getBodyByUrl(sendMessageUrl + "?chat_id=" + strconv.Itoa(item.Message.Chat.ID) + "&text=@" + item.Message.From.Username + " введите имя бота, к которому обращаетесь в формате: /команда @имя_бота" + "&parse_mode=HTML")
				// case "/stop":
				// 	break
				default:
					getBodyByUrl(sendMessageUrl + "?chat_id=" + strconv.Itoa(item.Message.Chat.ID) + "&text=" + url.QueryEscape(helpText["undefined"]) + "&parse_mode=HTML")
				}
			}
		}
	}
}
