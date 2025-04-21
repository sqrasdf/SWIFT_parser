package database

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"swift_parser/models"
)

func ConnectWithDatabase(sql_path string, csv_path string, env_path string) (*pgxpool.Pool, error) {
	// Load env variables
	godotenv.Load(env_path)
	dbUser := os.Getenv("DB_USER")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	fmt.Println("dbURL" + dbURL + "\n\n\n")
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Could not connect with database: %v", err)
		return nil, err
	}

	// Test connection
	err = dbpool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Error when pinging: %v", err)
		return nil, err
	}

	// Create database with sql file
	// err = executeSQLFile(context.Background(), dbpool, "schema.sql")
	dir, _ := os.Getwd()
	fmt.Println("Working directory:", dir)
	err = executeSQLFile(context.Background(), dbpool, sql_path)
	if err == nil {
		fmt.Println("SQL file executed correctly")
	}
	if err != nil {
		log.Fatalf("Error when executing SQL file: %v", err)
	}
	fmt.Println("Connected with database PostgreSQL")

	// Read data from csv file
	// hqMap, branches := parseCSV("SWIFT_CODES.csv")
	hqMap, branches := parseCSV(csv_path)

	// Insert headquarters to database
	licznik_bledow := 0
	for _, hq := range hqMap {
		err := insertHeadquarter(context.Background(), dbpool, &hq)
		if err != nil {
			licznik_bledow++
			log.Fatalf("Error - inserting headquarter: %v", err)
		}
	}

	fmt.Println("headquarters bledy:", licznik_bledow)

	// Insert branches to database
	for _, br := range branches {
		if br.HQSwiftCode == "" {
			continue
		}
		err := insertBranch(context.Background(), dbpool, &br)
		if err != nil {
			licznik_bledow++
			// log.Printf("Błąd zapisu oddziału %s: %v", br.SwiftCode, err)
			// log.Printf("Error - inserting branch %s", br.SwiftCode)
		}
	}

	fmt.Println("branches bledy:", licznik_bledow)

	return dbpool, nil
}

func executeSQLFile(ctx context.Context, db *pgxpool.Pool, filepath string) error {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, string(content))
	return err
}

func parseCSV(filename string) (map[string]models.Headquarter, []models.Branch) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	hqMap := make(map[string]models.Headquarter)
	var branches []models.Branch

	// Skipping row names
	for _, record := range records[1:] {
		if len(record) < 8 {
			continue
		}

		// Data normalization
		countryISO2 := strings.ToUpper(record[0])
		swiftCode := record[1]
		bankName := record[3]
		address := strings.TrimSpace(record[4])
		countryName := strings.ToUpper(record[6])

		// Checking if it is headquarter
		if len(swiftCode) == 11 && strings.HasSuffix(swiftCode, "XXX") {
			hq := models.Headquarter{
				SwiftCode:   swiftCode,
				BankName:    bankName,
				CountryISO2: countryISO2,
				CountryName: countryName,
				Address:     address,
			}
			hqMap[swiftCode] = hq
		} else {
			// Find headquarter for branch
			hqSwift := swiftCode[:8] + "XXX"
			branches = append(branches, models.Branch{
				SwiftCode:   swiftCode,
				HQSwiftCode: hqSwift,
				BankName:    bankName,
				CountryISO2: countryISO2,
				CountryName: countryName,
				Address:     address,
			})
		}
	}

	return hqMap, branches
}

func insertHeadquarter(ctx context.Context, db *pgxpool.Pool, hq *models.Headquarter) error {
	_, err := db.Exec(ctx, `
        INSERT INTO headquarters (swift_code, bank_name, country_iso2, country_name, address)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (swift_code) DO UPDATE SET
            bank_name = EXCLUDED.bank_name,
            country_iso2 = EXCLUDED.country_iso2,
            country_name = EXCLUDED.country_name,
            address = EXCLUDED.address`,
		hq.SwiftCode, hq.BankName, hq.CountryISO2, hq.CountryName, hq.Address)
	return err
}

func insertBranch(ctx context.Context, db *pgxpool.Pool, br *models.Branch) error {
	var exists bool
	err := db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM headquarters WHERE swift_code=$1)", br.HQSwiftCode).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error when checking headquarters: %w", err)
	}
	if !exists {
		return fmt.Errorf("headquarters with code %s doesn't exist", br.HQSwiftCode)
	}

	_, err = db.Exec(ctx, `
        INSERT INTO branches (swift_code, hq_swift_code, bank_name, country_iso2, country_name, address)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (swift_code) DO UPDATE SET
            bank_name = EXCLUDED.bank_name,
            country_iso2 = EXCLUDED.country_iso2,
            country_name = EXCLUDED.country_name,
            address = EXCLUDED.address`,
		br.SwiftCode, br.HQSwiftCode, br.BankName, br.CountryISO2, br.CountryName, br.Address)
	return err
}
