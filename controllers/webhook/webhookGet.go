package webhook

import (
	"multiBot/config"

	"github.com/gofiber/fiber/v2"
)

func checkWaToken(tokens []config.WaVerifyToken, waToken string) bool {
	for _, val := range tokens {
		if val.TOKEN == waToken {
			return true
		}
	}
	return false
}

func WaWebhookGet(c *fiber.Ctx) error {

	tokens := config.GetWaTokens()

	mode := c.Query("hub.mode")

	token := c.Query("hub.verify_token")

	challenge := c.Query("hub.challenge")

	if mode != "" && token != "" {

		tokenOk := checkWaToken(tokens, token)

		if mode == "subscribe" && tokenOk {
			return c.Status(fiber.StatusOK).Send([]byte(challenge))
		} else {
			return c.SendStatus(fiber.StatusForbidden)
		}

	}

	return c.SendStatus(fiber.StatusForbidden)
}
