package models

type Headquarter struct {
	SwiftCode   string `json:"swiftCode"`
	BankName    string `json:"bankName"`
	CountryISO2 string `json:"countryISO2"`
	CountryName string `json:"countryName"`
	Address     string `json:"address"`
}

type Branch struct {
	SwiftCode   string `json:"swiftCode"`
	HQSwiftCode string `json:"hqSwiftCode"`
	BankName    string `json:"bankName"`
	CountryISO2 string `json:"countryISO2"`
	CountryName string `json:"countryName"`
	Address     string `json:"address"`
}

// type SwiftCodeRequest struct {
// 	SwiftCode     string `json:"swiftCode" binding:"required,len=11"`
// 	BankName      string `json:"bankName" binding:"required"`
// 	CountryISO2   string `json:"countryISO2" binding:"required,len=2"`
// 	CountryName   string `json:"countryName" binding:"required"`
// 	Address       string `json:"address" binding:"required"`
// 	IsHeadquarter bool   `json:"isHeadquarter"`
// }

type SwiftCodeRequest struct {
	SwiftCode     string `json:"swiftCode"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName"`
	Address       string `json:"address"`
	IsHeadquarter bool   `json:"isHeadquarter"`
}

type HeadquarterResponse struct {
	SwiftCode     string        `json:"swiftCode"`
	BankName      string        `json:"bankName"`
	CountryISO2   string        `json:"countryISO2"`
	CountryName   string        `json:"countryName"`
	Address       string        `json:"address"`
	IsHeadquarter bool          `json:"isHeadquarter"`
	Branches      []BranchShort `json:"branches"`
}

type BranchShort struct {
	SwiftCode     string `json:"swiftCode"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	Address       string `json:"address"`
	IsHeadquarter bool   `json:"isHeadquarter"`
}

type BranchResponse struct {
	SwiftCode     string `json:"swiftCode"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName"`
	Address       string `json:"address"`
	IsHeadquarter bool   `json:"isHeadquarter"`
}
