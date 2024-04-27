package main

import (
	"context"
	"errors"
	"file-manager/data"
	"file-manager/routes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func checkEnv() {
	envs := []string{"DB_CONN_STRING"}
	for _, env := range envs {
		if value, success := os.LookupEnv(env); !success || value == "" {
			log.Fatalf("Environment variable %s is not set", env)
		}
	}
}

func main() {
	_ = godotenv.Load()
	checkEnv()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	data.ConnectToDB()
	defer data.CloseDB()

	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	router.GET("/ping", routes.PingRoute)

	// Media routes
	router.POST("/upload", routes.UploadMediaRoute)
	router.GET("/files/:id", routes.FileInfoRoute)
	router.GET("/files/random", routes.RandomFileRoute)
	router.DELETE("/files/:id", routes.DestroyRoute)
	router.GET("/files/:id/stream", routes.StreamVideo)
	router.GET("/files/:id/cover", routes.StreamCover)

	// Generic storage routes
	router.GET("/storage/:id", routes.StreamGeneric)
	router.POST("/storage", routes.UploadGenericRoute)

	srv := &http.Server{
		Addr:    ":8888",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	fmt.Printf("Server is running on port %s \n", srv.Addr)

	<-ctx.Done()

	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
