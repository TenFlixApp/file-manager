package main

import (
	"context"
	"errors"
	"file-manager/data"
	"file-manager/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func checkEnv() {
	envs := []string{"APP_PORT", "DB_CONN_STRING"}
	for _, env := range envs {
		if value, success := os.LookupEnv(env); !success || value == "" {
			log.Fatalf("Environment variable %s is not set", env)
		}
	}
}

func main() {
	_ = godotenv.Load()
	checkEnv()

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Connect to the database
	data.ConnectToDB()
	defer data.CloseDB()

	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
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

	appPort := os.Getenv("APP_PORT")
	srv := &http.Server{
		Addr:    ":" + appPort,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	fmt.Println("Server is running on port " + appPort)

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
