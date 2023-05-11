package main

import (
	"fmt"
	"multiBot/bootstrap"
	"multiBot/config"
	"multiBot/controllers"
	"multiBot/controllers/audioUrl"
	"multiBot/controllers/documentUrl"
	"multiBot/controllers/imageUrl"
	"multiBot/controllers/list"
	"multiBot/controllers/location"
	"multiBot/controllers/replyButton"
	"multiBot/controllers/stickerUrl"
	"multiBot/controllers/template"
	"multiBot/controllers/text"
	"multiBot/controllers/videoUrl"
	"multiBot/controllers/webhook"
	"multiBot/logger"
	"multiBot/middlewares"
	"multiBot/queue"
	"multiBot/service"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func main() {

	exitCtx, shutDown, err := bootstrap.Startup()

	if err != nil {
		fmt.Println("Bootup error: " + err.Error())
		os.Exit(1)
	}

	curEnv := os.Getenv("GO_ENV")

	cfg := config.GetConfig()

	var nrApp *newrelic.Application

	if curEnv == "prod" {
		nrApp, err = newrelic.NewApplication(
			newrelic.ConfigAppName("multiBot"),
			newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		)
		if err != nil {
			fmt.Println("newRelic error: " + err.Error())
			os.Exit(1)
		}
	}

	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "multiBot",
		AppName:       "multiBot v0.0.0",
	})

	app.Use(recover.New())

	app.Use(helmet.New())

	app.Use(cors.New())

	app.Use(requestid.New(requestid.Config{
		ContextKey: "requestId",
	}))

	app.Use(compress.New())

	// Request logger middlewares
	app.Use(middlewares.LogRequestMiddleware())

	if curEnv == "prod" {
		// Newrelic metric middleware
		app.Use(middlewares.New(middlewares.Config{
			NewRelicApp: nrApp,
		}))
	}

	app.Use(middlewares.CheckClient())

	bot := app.Group("/multiBot")

	api := bot.Group("/api")
	v1 := api.Group("/v1")
	// Health check
	v1.Get("/health", func(c *fiber.Ctx) error {
		res := map[string]interface{}{
			"data":   "Server is up and running",
			"errors": nil,
			"status": 200,
		}
		if err := c.JSON(res); err != nil {
			return err
		}
		return nil
	})

	// ========== Whatsapp =========== //

	whatsappApi := v1.Group("/whatsapp")
	// Retrieve message ids based on Status
	whatsappApi.Post("/waGetMessagesUsingStatus", controllers.ValiadateGetMessagesUsingStatusPayload, controllers.GetMessagesUsingStatus)
	// Retrieve count of msgs basis their status
	whatsappApi.Post("/waGetCount", controllers.ValiadateGetCountPayload, controllers.GetCount)
	// Retrieve messages status
	whatsappApi.Post("/waStatus", controllers.ValiadateStatusPayload, controllers.Status)
	// whatsapp webhook
	whatsappApi.Post("/waWebhook", webhook.ValidateWaWebhookPost, webhook.WaWebhookPost)
	whatsappApi.Get("/waWebhook", webhook.WaWebhookGet)
	// send group
	send := whatsappApi.Group("/send")
	send.Post("/text", text.ValidateTextPayload, text.Text)
	send.Post("/template", template.ValiadateTemplatePayload, template.Template)
	send.Post("/imageUrl", imageUrl.ValiadateImageUrlPayload, imageUrl.ImageUrl)
	send.Post("/audioUrl", audioUrl.ValiadateAudioUrlPayload, audioUrl.AudioUrl)
	send.Post("/documentUrl", documentUrl.ValiadateDocumentUrlPayload, documentUrl.DocumentUrl)
	send.Post("/stickerUrl", stickerUrl.ValiadateStickerUrlPayload, stickerUrl.StickerUrl)
	send.Post("/videoUrl", videoUrl.ValiadateVideoUrlPayload, videoUrl.VideoUrl)
	send.Post("/location", location.ValiadateLocationPayload, location.Location)
	send.Post("/list", list.ValiadateListPayload, list.List)
	send.Post("/replyButton", replyButton.ValiadateReplyButtonPayload, replyButton.ReplyButton)
	send.Post("/message", controllers.ValiadateMessage, controllers.Message)

	// ============================ //

	shutDownWg := &sync.WaitGroup{}
	if !cfg.WEBHOOK_FALLBACK {
		shutDownWg.Add(5)
		go bootstrap.GracefulShutDown(*exitCtx, app, shutDownWg)
		go service.ConsumerWhatsapp(*exitCtx, shutDownWg)
		go queue.Consumer(*exitCtx, shutDownWg)
		go queue.Producer(*exitCtx, shutDownWg)
		go service.UpdateStaleMessages(*exitCtx, shutDownWg)
	} else {
		shutDownWg.Add(3)
		go bootstrap.GracefulShutDown(*exitCtx, app, shutDownWg)
		go service.ConsumerWhatsapp(*exitCtx, shutDownWg)
		go service.UpdateStaleMessages(*exitCtx, shutDownWg)
	}

	err = app.Listen(":" + cfg.APP.PORT)
	if err != nil {
		logger.LogPanic(nil, "Fiber server error: "+err.Error(), nil)
		fmt.Println("Fiber server error: ", err)
		os.Exit(1)
	}

	// cleanUp
	bootstrap.CleanUp(*shutDown)
	// waiting for shutting down
	shutDownWg.Wait()
}
