package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"swift_parser/database"
	"swift_parser/handlers"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {

	// Inicjalizacja bazy danych
	dbpool, err := database.ConnectWithDatabase()
	if err != nil {
		log.Fatalf("Error connecting with database: %v", err)
	}

	// Inicjalizacja routera Gin
	r := gin.New()
	r.SetTrustedProxies([]string{"localhost"})

	// Endpointy
	r.GET("/", func(c *gin.Context) {
		fmt.Println("DB_NAME", os.Getenv("DB_NAME"))
		c.String(http.StatusOK, "hello world")
	})
	r.GET("/v1/swift-codes/:swift-code", handlers.GetSWIFTCode(dbpool))
	r.GET("/v1/swift-codes/country/:countryISO2code", handlers.GetSWIFTCodesByCountry(dbpool))
	r.POST("/v1/swift-codes", handlers.PostSwiftCode(dbpool))
	r.DELETE("/v1/swift-codes/:swift-code", handlers.DeleteSwiftCode(dbpool))

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigchan
		log.Println("Otrzymano sygnał zakończenia, zamykanie połączenia z bazą danych...")
		dbpool.Close()
		log.Println("Połączenie z bazą danych zamknięte.")
		os.Exit(0)
	}()

	// Uruchomienie serwera
	r.Run(":8080")
}
