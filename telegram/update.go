/*
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package telegram

import (
	"fmt"
	"go-FitgirlRepacks-Telegram-bot/fitgirl"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (tg *Telegram) HandleUpdate(update tgbotapi.Update) {
	// If the update is from a Callback, search the GameID in the CallbackData
	if update.CallbackQuery != nil {
		var keyboardRows = make([][]tgbotapi.InlineKeyboardButton, 0)
		var gameTitle string
		var magnet string = ""
		var message string = "Download links for "

		// Convert the Game ID from string to int
		gameID, _ := strconv.Atoi(update.CallbackQuery.Data)
		// Find the download link for that Game ID
		downloadLinks := fitgirl.FindDownloadLinks(uint16(gameID))
		// Find the game title of that Game ID
		gameTitle = fitgirl.FindGameTitle(uint16(gameID))
		// Add the game title to the message text
		message += "`" + gameTitle + "`:"

		for i := range downloadLinks {
			// Don't add the magnet to the keyboard
			if downloadLinks[i].Name == "magnet" {
				magnet = downloadLinks[i].Link
			} else {
				// Add the links to the inline buttons array
				keyboardRows = append(keyboardRows,
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonURL(downloadLinks[i].Name, downloadLinks[i].Link),
					))
			}

		}

		// Create the keyboard using the inline button array
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			keyboardRows...,
		)

		// If there was a magnet link, add it to the message text
		if magnet != "" {
			message += "\nMagnet: `" + magnet + "`"
		}

		// Generate the message
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, message)
		msg.ReplyMarkup = keyboard
		msg.ParseMode = tgbotapi.ModeMarkdownV2
		_, err := tg.api.Send(msg)
		if err != nil {
			fmt.Println(err)
		}

		tg.api.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
	}

	// Skip if there isn't a real update
	if update.Message == nil && update.EditedMessage == nil {
		return
	}

	// If message is edited, set it as message to handle
	if update.EditedMessage != nil {
		update.Message = update.EditedMessage
	}

	// Skip messages from channels
	if update.Message.Chat.Type == "channel" {
		return
	}

	// Check for ban and inform the user only on private chat to avoid flood
	if tg.db.FindBan(int64(update.Message.From.ID)) >= 0 {
		if update.Message.Chat.Type == "private" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🚫 You have been banned from this bot!")
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = tg.api.Send(msg)
		}

		return
	}

	// Commands
	if len(update.Message.Text) >= 5 && strings.ToLower(update.Message.Text[0:5]) == "/help" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "*HELP*\n\n*/help* \\- to use this command, view available commands\\.\n*/search <uid\\>* \\- Search a game on Fitgirl-Repacks website\\! 🦝")
		msg.ParseMode = tgbotapi.ModeMarkdownV2
		_, _ = tg.api.Send(msg)
		return
	}

	if len(update.Message.Text) >= 7 && strings.ToLower(update.Message.Text[0:7]) == "/search" {
		message := strings.SplitN(update.Message.Text, " ", 2)
		var searchQuery string
		var keyboardRows = make([][]tgbotapi.InlineKeyboardButton, 0)

		if len(message) == 2 {
			searchQuery = message[1]
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: /search <query>")
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = tg.api.Send(msg)
			return
		}

		searchRes := fitgirl.Search(searchQuery)

		for i := range searchRes {
			keyboardRows = append(keyboardRows,
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(searchRes[i].Name, strconv.Itoa(int(searchRes[i].ID))),
				))
		}

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			keyboardRows...,
		)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Search results for the query `"+searchQuery+"`")
		msg.ParseMode = tgbotapi.ModeMarkdownV2
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ReplyMarkup = keyboard
		_, _ = tg.api.Send(msg)

	}
}
