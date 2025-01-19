package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/ARF-DEV/image-processing-api/configs"
	"github.com/ARF-DEV/image-processing-api/handlers"
	"github.com/ARF-DEV/image-processing-api/handlers/imagehand"
	"github.com/ARF-DEV/image-processing-api/handlers/userhand"
	"github.com/ARF-DEV/image-processing-api/repos/googlecloudstorage"
	"github.com/ARF-DEV/image-processing-api/repos/imagerepo"
	"github.com/ARF-DEV/image-processing-api/repos/userrepo"
	"github.com/ARF-DEV/image-processing-api/services/imageserv"
	"github.com/ARF-DEV/image-processing-api/services/userserv"
)

func main() {
	configs.LoadConfig()
	cfg := configs.GetConfig()
	db, err := configs.SetupDB(cfg.DB_MASTER)
	if err != nil {
		panic(err)
	}
	userRepo := userrepo.New(db)

	gcsRepo := googlecloudstorage.New(context.Background(), cfg)
	defer gcsRepo.Close()

	imageRepo := imagerepo.New(db)

	userServ := userserv.New(userRepo)
	imageServ := imageserv.New(gcsRepo, imageRepo)

	imageHand := imagehand.New(imageServ)
	userHand := userhand.New(userServ)

	h := handlers.CreateHandlers(userHand, imageHand)

	server := http.Server{
		Addr:    ":8080",
		Handler: h,
	}
	go func() {
		log.Println("server is now listening at port :8080")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("HTTP server ListenAndServe: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}

}
