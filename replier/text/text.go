package text

import (
	"multiBot/clients/openai"
	"multiBot/config"
	"multiBot/constants"

	"github.com/gofiber/fiber/v2"
)

// This method replies to text messages from user
func ReplyToText(c *fiber.Ctx, msgText string, from string, waAccount *config.WhatsappAccountsConfig) error {
	if waAccount.IDENTIFIER == constants.GptBotWhatsappAccountIdentifier { // Checking if this account is GPT_BOT Account
		return openai.SendReply(c, msgText, from, waAccount)
	}
	return c.SendStatus(fiber.StatusOK)
}
