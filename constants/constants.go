package constants

var NonMetricRoutes = []string{
	"/multiBot/metrics",
	"/multiBot/app_metrics",
	"/multiBot/api/v1/health",
}

var ClientIdSkipRoutes = []string{
	"/multiBot/api/v1/whatsapp/waWebhook",
	"/multiBot/api/v1/whatsapp/waStatus",
	"/multiBot/api/v1/whatsapp/waGetCount",
	"/multiBot/api/v1/whatsapp/waGetMessagesUsingStatus",
}

const SystemClientIdentifier = "9c9b6d8388dede5d2f3f20852601119db9ad056849c7f3b297a64606772f56ce"

const GptBotWhatsappAccountIdentifier = "5eebd2df394fa1749028680f38feeb0826e1be13d46589bae3c319ed088d01e3"
