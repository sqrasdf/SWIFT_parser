package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"swift_parser/database"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupRouter(dbpool *pgxpool.Pool) *gin.Engine {
	r := gin.New()
	// r.GET("/v1/swift-codes/:swift-code", GetSWIFTCode(dbpool))

	r.GET("/v1/swift-codes/:swift-code", GetSWIFTCode(dbpool))
	r.GET("/v1/swift-codes/country/:countryISO2code", GetSWIFTCodesByCountry(dbpool))
	r.POST("/v1/swift-codes", PostSwiftCode(dbpool))
	r.DELETE("/v1/swift-codes/:swift-code", DeleteSwiftCode(dbpool))

	return r
}

func TestGetSWIFTCode_NotFound(t *testing.T) {
	// dbpool, err := setupTestDB()
	dbpool, err := database.ConnectWithDatabase("../database/schema.sql", "../data_csv/SWIFT_CODES.csv", "../.env")
	if err != nil {
		t.Fatalf("Error connection with database: %v", err)
	}
	defer dbpool.Close()

	router := setupRouter(dbpool)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/swift-codes/NOTEXISTXXX", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "not found") {
		t.Errorf("Expected no record, got: %s", w.Body.String())
	}
}

func TestPostSwiftCode_InvalidData(t *testing.T) {
	dbpool, err := database.ConnectWithDatabase("../database/schema.sql", "../data_csv/SWIFT_CODES.csv", "../.env")
	if err != nil {
		t.Fatalf("Error connectig with database: %v", err)
	}
	router := setupRouter(dbpool)

	jsonStr := `{"swiftCode":"SHORT","bankName":"Test Bank","countryISO2":"US","countryName":"UNITED STATES","address":"123 Test St","isHeadquarter":true}`
	// jsonStr := `{ "address": "nice address", "bankName": "nice_bank_name", "countryISO2": "XD", "countryName": "XDland", "isHeadquarter": true, "swiftCode": "12345678XXX" }`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/swift-codes", strings.NewReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "SWIFT code must be 11 characters long") {
		t.Errorf("Expected length error message, got %s", w.Body.String())
	}
}

func TestDeleteSwiftCode_NotFound(t *testing.T) {
	dbpool, err := database.ConnectWithDatabase("../database/schema.sql", "../data_csv/SWIFT_CODES.csv", "../.env")
	if err != nil {
		t.Fatalf("Error connecting with database: %v", err)
	}
	router := setupRouter(dbpool)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/swift-codes/NOTEXISTXXX", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found, got %d", w.Code)
	}
}
