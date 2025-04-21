package main

import (
	// "fmt"
	// "log"
	// "net/http"
	// "os"
	// "os/signal"
	// "swift_parser/database"
	// "swift_parser/handlers"
	// "syscall"

	// "github.com/gin-gonic/gin"

	"swift_parser/handlers"
)

func main() {

	// Initialize database
	// dbpool, err := database.ConnectWithDatabase("database/schema.sql", "data_csv/SWIFT_CODES.csv", ".env")
	// if err != nil {
	// 	log.Fatalf("Error connecting with database: %v", err)
	// }

	// // Initialize router
	// r := gin.New()
	// r.SetTrustedProxies([]string{"localhost"})

	// r.GET("/", func(c *gin.Context) {
	// 	fmt.Println("DB_NAME", os.Getenv("DB_NAME"))
	// 	c.String(http.StatusOK, "hello world")
	// })
	// r.GET("/v1/swift-codes/:swift-code", handlers.GetSWIFTCode(dbpool))
	// r.GET("/v1/swift-codes/country/:countryISO2code", handlers.GetSWIFTCodesByCountry(dbpool))
	// r.POST("/v1/swift-codes", handlers.PostSwiftCode(dbpool))
	// r.DELETE("/v1/swift-codes/:swift-code", handlers.DeleteSwiftCode(dbpool))

	// // Signal closing connection with database
	// sigchan := make(chan os.Signal, 1)
	// signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// go func() {
	// 	<-sigchan
	// 	log.Println("Closing connection with database...")
	// 	dbpool.Close()
	// 	log.Println("Connection with database has been closed")
	// 	os.Exit(0)
	// }()

	// // Run server
	// r.Run(":8080")

	handlers.InitializeRouter()
}
