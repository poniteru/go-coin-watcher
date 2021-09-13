package digicoinutil

import (
	"strings"
)

func ToSingleCurrencyName(currencyPair string) string {
	if strings.Contains(currencyPair, ":") {
		return strings.Split(currencyPair, ":")[0]
	} else if strings.Contains(currencyPair, "/") {
		return strings.Split(currencyPair, "/")[0]
	}
	return currencyPair
}

func ToDefaultCurrencyPair(currencyPair string) string {
	if strings.Contains(currencyPair, ":") {
		return currencyPair
	} else if strings.Contains(currencyPair, "/") {
		return strings.ReplaceAll(currencyPair, "/", ":")
	}
	return currencyPair + ":usdt"
}

func ToHuobiWsDefaultCurrencyPair(currencyPair string) string {
	if strings.Contains(currencyPair, ":") {
		return strings.ReplaceAll(currencyPair, ":", "")
	} else if strings.Contains(currencyPair, "/") {
		return strings.ReplaceAll(currencyPair, "/", "")
	}
	return currencyPair + "usdt"
}

func ToDisplayCurrencyPair(currencyPair string) string {
	if strings.Contains(currencyPair, ":") {
		return strings.ReplaceAll(currencyPair, ":", "/")
	} else if strings.Contains(currencyPair, "/") {
		return currencyPair
	}
	return currencyPair + "/usdt"
}
