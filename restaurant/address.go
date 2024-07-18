package restaurant

type Address struct {
	AddressType           string  `json:"addressType,omitempty"`
	AddressLine1          string  `json:"addressLine1,omitempty"`
	AddressLine2          string  `json:"addressLine2,omitempty"`
	Locality              string  `json:"locality,omitempty"`
	AdministrativeArea    string  `json:"administrativeArea,omitempty"`
	PostalCode            string  `json:"postalCode,omitempty"`
	SubAdministrativeArea string  `json:"subAdministrativeArea,omitempty"`
	CountryCode           string  `json:"countryCode,omitempty"`
	Latitude              float64 `json:"latitude,omitempty"`
	Longitude             float64 `json:"longitude,omitempty"`
	AccuracyDetermination string  `json:"accuracyDetermination,omitempty"`
}
