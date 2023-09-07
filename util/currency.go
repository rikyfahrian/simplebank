package util

const (
	RP  = "RP"
	EUR = "EUR"
	USD = "USD"
)

func IsSupportedCurrency(currency string) bool {

	switch currency {
	case RP, EUR, USD:
		return true
	}

	return false

}
