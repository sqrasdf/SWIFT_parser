package handlers

import (
	"context"
	// "encoding/csv"
	"errors"
	"net/http"
	"strings"
	"swift_parser/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetSWIFTCode(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		swiftCode := c.Param("swift-code")

		if len(swiftCode) != 11 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SWIFT code length"})
			return
		}

		isHQ := strings.HasSuffix(swiftCode, "XXX")

		if isHQ {
			var hq models.HeadquarterResponse
			err := dbpool.QueryRow(context.Background(), `
			SELECT swift_code, bank_name, country_iso2, country_name, address
			FROM headquarters
			WHERE swift_code = $1`, swiftCode).
				Scan(&hq.SwiftCode, &hq.BankName, &hq.CountryISO2, &hq.CountryName, &hq.Address)

			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "SWIFT code not found"})
				return
			}

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				return
			}

			hq.IsHeadquarter = true

			rows, _ := dbpool.Query(context.Background(), `
			SELECT swift_code, bank_name, country_iso2, address
			FROM branches
			WHERE hq_swift_code = $1`, swiftCode)
			defer rows.Close()

			var branches []models.BranchShort
			for rows.Next() {
				var br models.BranchShort
				rows.Scan(&br.SwiftCode, &br.BankName, &br.CountryISO2, &br.Address)
				br.IsHeadquarter = false
				branches = append(branches, br)
			}

			hq.Branches = branches
			c.JSON(http.StatusOK, hq)
		} else {
			var br models.BranchResponse
			err := dbpool.QueryRow(context.Background(), `
			SELECT swift_code, bank_name, country_iso2, country_name, address
			FROM branches
			WHERE swift_code = $1`, swiftCode).
				Scan(&br.SwiftCode, &br.BankName, &br.CountryISO2, &br.CountryName, &br.Address)

			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "SWIFT code not found"})
				return
			}

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				return
			}

			br.IsHeadquarter = false
			c.JSON(http.StatusOK, br)
		}
	}
}

func GetSWIFTCodesByCountry(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		countryISO2 := strings.ToUpper(c.Param("countryISO2code"))

		hqRows, _ := dbpool.Query(context.Background(), `
		SELECT swift_code, bank_name, country_iso2, country_name, address
		FROM headquarters
		WHERE country_iso2 = $1`, countryISO2)
		defer hqRows.Close()

		var responseData []gin.H
		var countryName string

		for hqRows.Next() {
			var swiftCode, bankName, iso2, country, address string
			hqRows.Scan(&swiftCode, &bankName, &iso2, &country, &address)

			if countryName == "" {
				countryName = country
			}

			responseData = append(responseData, gin.H{
				"swiftCode":     swiftCode,
				"bankName":      bankName,
				"countryISO2":   iso2,
				"isHeadquarter": true,
				"address":       address,
			})
		}

		branchRows, _ := dbpool.Query(context.Background(), `
		SELECT swift_code, bank_name, country_iso2, country_name, address
		FROM branches
		WHERE country_iso2 = $1`, countryISO2)
		defer branchRows.Close()

		for branchRows.Next() {
			var swiftCode, bankName, iso2, country, address string
			branchRows.Scan(&swiftCode, &bankName, &iso2, &country, &address)

			if countryName == "" {
				countryName = country
			}

			responseData = append(responseData, gin.H{
				"swiftCode":     swiftCode,
				"bankName":      bankName,
				"countryISO2":   iso2,
				"isHeadquarter": false,
				"address":       address,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"countryISO2": countryISO2,
			"countryName": countryName,
			"swiftCodes":  responseData,
		})
	}
}

func PostSwiftCode(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.SwiftCodeRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data: " + err.Error()})
			return
		}

		// Data normalization
		req.SwiftCode = strings.ToUpper(req.SwiftCode)
		req.CountryISO2 = strings.ToUpper(req.CountryISO2)
		req.CountryName = strings.ToUpper(req.CountryName)

		if len(req.SwiftCode) != 11 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "SWIFT code must be 11 characters long"})
			return
		}

		if !isValidSwiftCodeFormat(req.SwiftCode) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "SWIFT code contains invalid characters"})
			return
		}

		isHQ := strings.HasSuffix(req.SwiftCode, "XXX")
		if req.IsHeadquarter != isHQ {
			c.JSON(http.StatusBadRequest, gin.H{"error": "SWIFT code of headquarter must end with 'XXX'"})
			return
		}

		if strings.TrimSpace(req.BankName) == "" || strings.TrimSpace(req.Address) == "" || len(req.CountryISO2) != 2 || strings.TrimSpace(req.CountryName) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields have to be valid"})
			return
		}

		ctx := context.Background()

		if req.IsHeadquarter {
			_, err := dbpool.Exec(ctx, `
			INSERT INTO headquarters (swift_code, bank_name, country_iso2, country_name, address)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (swift_code) DO UPDATE SET
			bank_name = EXCLUDED.bank_name,
			country_iso2 = EXCLUDED.country_iso2,
			country_name = EXCLUDED.country_name,
			address = EXCLUDED.address`,
				req.SwiftCode, req.BankName, req.CountryISO2, req.CountryName, req.Address)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error: " + err.Error()})
				return
			}
		} else {
			hqSwiftCode := req.SwiftCode[:8] + "XXX"
			exists, err := checkHeadquarterExists(ctx, dbpool, hqSwiftCode)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error: " + err.Error()})
				return
			}

			if !exists {
				c.JSON(http.StatusBadRequest, gin.H{"error": "No headquarter with code " + hqSwiftCode})
				return
			}

			_, err = dbpool.Exec(ctx, `
			INSERT INTO branches (swift_code, hq_swift_code, bank_name, country_iso2, country_name, address)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (swift_code) DO UPDATE SET
			bank_name = EXCLUDED.bank_name,
			country_iso2 = EXCLUDED.country_iso2,
			country_name = EXCLUDED.country_name,
			address = EXCLUDED.address`,
				req.SwiftCode, hqSwiftCode, req.BankName, req.CountryISO2, req.CountryName, req.Address)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error: " + err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "SWIFT code saved successfully"})
	}
}

func DeleteSwiftCode(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		swiftCode := strings.ToUpper(c.Param("swift-code"))

		if len(swiftCode) != 11 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "SWIFT code must be 11 characters long"})
			return
		}

		ctx := context.Background()

		isHQ := strings.HasSuffix(swiftCode, "XXX")

		var cmdTag pgconn.CommandTag
		var err error

		if isHQ {
			cmdTag, err = dbpool.Exec(ctx, `DELETE FROM headquarters WHERE swift_code=$1`, swiftCode)
		} else {
			cmdTag, err = dbpool.Exec(ctx, `DELETE FROM branches WHERE swift_code=$1`, swiftCode)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
			return
		}

		if cmdTag.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "SWIFT code not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "SWIFT code deleted successfully"})
	}
}

func isValidSwiftCodeFormat(code string) bool {
	for _, r := range code {
		if !(r >= 'A' && r <= 'Z') && !(r >= '0' && r <= '9') {
			return false
		}
	}
	return true
}

func checkHeadquarterExists(ctx context.Context, db *pgxpool.Pool, swiftCode string) (bool, error) {
	var exists bool
	err := db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM headquarters WHERE swift_code=$1)", swiftCode).Scan(&exists)
	return exists, err
}
