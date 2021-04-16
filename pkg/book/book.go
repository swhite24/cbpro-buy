package book

import (
	"strconv"

	"github.com/preichenberger/go-coinbasepro/v2"
)

// GetPrice delivers the best ask price for a particular product
func GetPrice(client *coinbasepro.Client, product string) (float64, error) {
	book, err := client.GetBook(product, 1)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(book.Asks[0].Price, 64)
}
