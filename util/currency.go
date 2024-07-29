package util

// constants for currency
const (
	USD = "USD"
	CAD = "CAD"
	EUR = "EUR"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	}
	return false
}

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedRole(role string) bool {
	switch role {
	case BankerRole, DepositorRole:
		return true
	}
	return false
}
