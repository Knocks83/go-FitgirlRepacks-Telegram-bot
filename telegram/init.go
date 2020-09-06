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
	"FitgirlBot/roles"
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Telegram bot type
type Telegram struct {
	api *tgbotapi.BotAPI
	db  *roles.Roles
}

// NewTelegramBot create a new Telegram bot instance from a token
// Returns a pointer to Telegram struct
func NewTelegramBot(token string, database *roles.Roles) (*Telegram, error) {
	// Create new variables
	bot := new(Telegram)
	var err error

	// Create bot instance from token
	bot.api, err = tgbotapi.NewBotAPI(token)

	if err != nil {
		return nil, err
	}

	// Check if input roles pointer is valid
	if database == nil {
		return nil, errors.New("roles pointer is nil, unable to configure bot")
	}

	// Assign roles to Telegram bot struct
	bot.db = database

	return bot, nil
}

// ManageUpdates starts a polling loop and handle every update with an handler function
func (tg *Telegram) ManageUpdates(handler func(tgbotapi.Update)) error {
	// Configure a new getUpdates() method
	settings := tgbotapi.NewUpdate(0)
	settings.Timeout = 60

	// Create a polling goroutine and get the channel for the communication
	updates, err := tg.api.GetUpdatesChan(settings)

	// Check for errors
	if err != nil {
		return err
	}

	// Create a new goroutine for the updates management and handle every update in a goroutine
	go func(updates tgbotapi.UpdatesChannel) {
		for update := range updates {
			go handler(update)
		}
	}(updates)

	// No errors
	return nil
}
