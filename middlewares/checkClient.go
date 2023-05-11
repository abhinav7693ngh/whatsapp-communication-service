package middlewares

import (
	"multiBot/config"
	"multiBot/constants"
	"multiBot/errorCodes"
	"multiBot/utils"

	"github.com/gofiber/fiber/v2"
)

func checkValidApiKey(apps []config.ClientsConfig, apiKey string) *config.ClientsConfig {
	for _, val := range apps {
		if val.API_KEY == apiKey {
			return &val
		}
	}
	return nil
}

func checkValidClientWaAcc(waAccounts []config.WhatsappAccountsConfig, waAccHeader string, apps []config.ClientsConfig, apiKey string) (*config.ClientsConfig, *config.WhatsappAccountsConfig) {
	clientInfo := checkValidApiKey(apps, apiKey)
	if clientInfo != nil {
		for _, val := range waAccounts {
			if val.WA_HEADER == waAccHeader {
				for _, client := range val.CLIENTS {
					if client == clientInfo.IDENTIFIER {
						return clientInfo, &val
					}
				}
				return nil, nil
			}
		}
	}
	return nil, nil
}

func CheckClient() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if FindStringInSlice(constants.NonMetricRoutes, c.Path()) {
			return c.Next()
		}
		if FindStringInSlice(constants.ClientIdSkipRoutes, c.Path()) {
			return c.Next()
		}
		apiKey := c.Get("x-api-key")
		waAccHeader := c.Get("x-wa-account")
		if apiKey != "" && waAccHeader != "" {
			cfg := config.GetConfig()
			clients := cfg.CLIENTS
			waAccounts := cfg.WHATSAPP_ACCOUNTS_CONFIG
			validClient, validWaAccount := checkValidClientWaAcc(waAccounts, waAccHeader, clients, apiKey)
			if validClient != nil {
				c.Locals("clientIdentifier", validClient.IDENTIFIER)
				c.Locals("waAccountIdentifier", validWaAccount.IDENTIFIER)
			} else {
				return utils.ErrorResponse(
					c,
					utils.AppError{
						DebugMessage: "",
						ErrorCode:    errorCodes.INVALID_API_KEY_OR_WA_ACCOUNT,
					},
					nil,
				)
			}
		} else {
			return utils.ErrorResponse(
				c,
				utils.AppError{
					DebugMessage: "",
					ErrorCode:    errorCodes.HEADER_NOT_FOUND,
				},
				nil,
			)
		}
		return c.Next()
	}
}
