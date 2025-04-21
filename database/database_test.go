package database

import (
	"context"
	"fmt"
	"os"
	"swift_parser/models"
	"testing"
)

func TestInsertAndGetHeadquarter(t *testing.T) {
	// godotenv.Load("../.env")

	dir, _ := os.Getwd()
	fmt.Println("Working directory:", dir)

	dbpool, err := ConnectWithDatabase("../database/schema.sql", "../data_csv/SWIFT_CODES.csv", "../.env")
	if err != nil {
		t.Fatalf("Error connecting with database: %v", err)
	}
	defer dbpool.Close()

	ctx := context.Background()

	hq := &models.Headquarter{
		SwiftCode:   "AAABBBCCXXX",
		BankName:    "Bank Testow",
		CountryISO2: "PL",
		CountryName: "Poland",
		Address:     "ul. Testowa 21",
	}

	err = insertHeadquarter(ctx, dbpool, hq)
	if err != nil {
		t.Fatalf("Error inserting headquarter: %v", err)
	}

	var count int
	err = dbpool.QueryRow(ctx, "SELECT COUNT(*) FROM headquarters WHERE swift_code=$1", hq.SwiftCode).Scan(&count)
	if err != nil {
		t.Fatalf("Błąd zapytania: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 record, got %d", count)
	}

	// Cleanup
	_, err = dbpool.Exec(ctx, "DELETE FROM headquarters WHERE swift_code=$1", hq.SwiftCode)
	if err != nil {
		t.Logf("Error when cleaning after test: %v", err)
	}
}

// func TestHehe(t *testing.T) {
// 	value := 1
// 	if value != 0 {
// 		t.Fatal("\n\n tescik nie przechodzi \n\n")
// 	}
// }
