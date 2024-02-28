package util

// Constants for all supported currencies
const (
	USD = "USD"
	INR = "INR"
	CAD = "CAD"
	EUR = "EUR"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, INR, CAD, EUR:
		return true
	default:
		return false
	}
}
