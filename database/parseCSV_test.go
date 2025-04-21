package database

import (
	"os"
	"testing"
)

func TestParseCSV(t *testing.T) {
	// sample csv data
	csvContent := `COUNTRY ISO2 CODE,SWIFT CODE,CODE TYPE,NAME,ADDRESS,TOWN NAME,COUNTRY NAME,TIME ZONE
PL,ABCDEFGHXXX,BIC11,Bank Testowy,"ul. Testowa 1, Warszawa",Warszawa,POLAND,Europe/Warsaw
PL,ABCDEFGH001,BIC11,Bank Testowy Oddział,"ul. Testowa 2, Kraków",Kraków,POLAND,Europe/Warsaw
`
	tmpFile, err := os.CreateTemp("", "test_swift_*.csv")
	if err != nil {
		t.Fatalf("Cannot open temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(csvContent)
	if err != nil {
		t.Fatalf("Cannot save to the file: %v", err)
	}
	tmpFile.Close()

	hqMap, branches := parseCSV(tmpFile.Name())

	if len(hqMap) != 1 {
		t.Errorf("Expected 1 headquarter, got %d", len(hqMap))
	}

	if len(branches) != 1 {
		t.Errorf("Expected 1 branch, got %d", len(branches))
	}

	hq, exists := hqMap["ABCDEFGHXXX"]
	if !exists {
		t.Errorf("Headquarter with code ABCDEFGHXXX not found")
	}
	if hq.BankName != "Bank Testowy" {
		t.Errorf("Incorrect bank name: %s", hq.BankName)
	}

	if branches[0].HQSwiftCode != "ABCDEFGHXXX" {
		t.Errorf("Branch does not match headquarter")
	}
}
