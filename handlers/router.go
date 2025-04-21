package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"swift_parser/database"

	"github.com/gin-gonic/gin"
)

func InitializeRouter() {
	dbpool, err := database.ConnectWithDatabase("database/schema.sql", "data_csv/SWIFT_CODES.csv", ".env")
	if err != nil {
		log.Fatalf("Error connecting with database: %v", err)
	}

	r := gin.New()
	r.SetTrustedProxies([]string{"localhost"})

	r.GET("/", func(c *gin.Context) {
		fmt.Println("DB_NAME", os.Getenv("DB_NAME"))
		c.String(http.StatusOK, "hello world")
	})
	r.GET("/v1/swift-codes/:swift-code", GetSWIFTCode(dbpool))
	r.GET("/v1/swift-codes/country/:countryISO2code", GetSWIFTCodesByCountry(dbpool))
	r.POST("/v1/swift-codes", PostSwiftCode(dbpool))
	r.DELETE("/v1/swift-codes/:swift-code", DeleteSwiftCode(dbpool))

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigchan
		log.Println("Closing connection with database...")
		dbpool.Close()
		log.Println("Connection with database has been closed")
		os.Exit(0)
	}()

	r.Run(":8080")
}
