package main

import (
	"fmt"
	"harago/cmd"
	"harago/config"
	"harago/gservice"
	"harago/gservice/gchat"
	"harago/handler"
	"harago/repository"
	"harago/stream"
	"log"
	"net/http"

	harborModel "github.com/dukhyungkim/harbor-client/model"
	"github.com/gofiber/fiber/v2"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	log.Println("Harago Starting.")

	opts, err := config.ParseFlags()
	if err != nil {
		log.Panicln(err)
	}

	cfg, err := config.NewConfig(opts)
	if err != nil {
		log.Panicln(err)
	}

	db, err := repository.NewPostgres(&cfg.DB)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("connect to postgres ... success")

	etcdClient, err := repository.NewEtcd(&cfg.Etcd)
	if err != nil {
		log.Panicln(err)
	}
	defer etcdClient.Close()
	log.Println("connect to etcd ... success")

	go etcdClient.WatchConfigList()

	gService, err := gservice.NewGService(opts.Credential)
	if err != nil {
		log.Println(err)
	}

	streamClient, err := stream.NewClient(&cfg.Nats)
	if err != nil {
		log.Panicln(err)
	}
	defer streamClient.Close()
	log.Println("connect to nats ... success")

	executor := cmd.NewExecutor()
	if err = executor.LoadCommands(cfg, streamClient, etcdClient); err != nil {
		log.Panicln(err)
	}

	dmHandler := handler.NewDMHandler(executor, db)
	roomHandler := handler.NewRoomHandler(executor, db)
	gChat, err := gchat.NewGChat(gService, dmHandler, roomHandler)
	if err != nil {
		log.Panicln(err)
	}

	respHandler := handler.NewResponseHandler(gChat, db)
	if err = streamClient.ClamResponse(respHandler.NotifyResponse); err != nil {
		log.Panicln(err)
	}

	harborEventHandle := handler.NewHarborEventHandler(streamClient, etcdClient)

	app := setup(gChat, harborEventHandle)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server startup ... %s\n", addr)
	log.Panicln(app.Listen(addr))
}

func setup(gChat *gchat.GChat, harborEventHandler *handler.HarborEventHandler) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Post("/message", func(ctx *fiber.Ctx) error {
		var event gchat.ChatEvent
		if err := ctx.BodyParser(&event); err != nil {
			log.Println(err)
		}
		return ctx.JSON(gChat.HandleMessage(&event))
	})

	app.Post("/harbor_notify", func(ctx *fiber.Ctx) error {
		var event harborModel.WebhookEvent
		if err := ctx.BodyParser(&event); err != nil {
			log.Println(err)
		}
		go harborEventHandler.HandleHarborEvent(&event)
		return ctx.SendStatus(http.StatusOK)
	})

	return app
}
