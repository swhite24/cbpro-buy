package purchase

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/preichenberger/go-coinbasepro/v2"
	"github.com/swhite24/cbpro-buy/pkg/config"
)

var (
	client *coinbasepro.Client
)

// InitiatePurchase conducts the necessary operations to deposit funds and purchase crypto
func InitiatePurchase(cfg *config.Config) error {
	var err error
	var fundingAccount *coinbasepro.Account
	var initialBalance float64

	client = coinbasepro.NewClient()
	client.UpdateConfig(&coinbasepro.ClientConfig{
		BaseURL:    cfg.BaseURL,
		Key:        cfg.Key,
		Passphrase: cfg.Passphrase,
		Secret:     cfg.Secret,
	})

	// Gather details on current account balance to know if funds are available
	fmt.Println("Fetching current funding account status.")
	if fundingAccount, err = getAccount(cfg.Currency); err != nil {
		return err
	}
	if initialBalance, err = strconv.ParseFloat(fundingAccount.Balance, 64); err != nil {
		return err
	}

	fmt.Printf("Success.  Available balance: %f\n", initialBalance)

	// Check if current balance is sufficient
	if initialBalance < cfg.Amount {
		fmt.Printf("Available balance is less than requested purchase: %f\n", cfg.Amount)
		if !cfg.AutoDeposit {
			return errors.New("Insufficient funds for purchase")
		}
		// initiate and wait for deposit
		fmt.Printf("Initiating deposit of %f %s\n", cfg.Amount, cfg.Currency)
		if err = deposit(cfg.Currency, cfg.Amount, cfg.UseCoinbase); err != nil {
			return err
		}

		// Wait for balance to be available
		ready := make(chan int, 1)

		go func(ch chan int, cfg *config.Config) {
			for {
				account, err := getAccount(cfg.Currency)
				if err != nil {
					continue
				}

				balance, err := strconv.ParseFloat(account.Balance, 64)
				if err != nil {
					continue
				}

				fmt.Printf("Checking available balance: %f\n", balance)

				if balance >= cfg.Amount {
					ch <- 1
				}
				time.Sleep(3 * time.Second)
			}
		}(ready, cfg)

		fmt.Println("Waiting for deposit to be available in account.")
		select {
		case <-ready:
			break
		case <-time.After(30 * time.Second):
			return errors.New("deposit did not clear in configured time window, please try again")
		}
	}

	// Make purchase
	fmt.Printf("Initiating purchase of %f %s worth of %s\n", cfg.Amount, cfg.Currency, cfg.Product)
	return purchase(cfg.Product, cfg.Currency, cfg.Amount)
}

func getAccount(typ string) (*coinbasepro.Account, error) {
	var err error
	var accounts []coinbasepro.Account
	accounts, err = client.GetAccounts()
	if err != nil {
		return nil, err
	}

	for _, a := range accounts {
		if a.Currency == typ {
			return &a, nil
		}
	}

	return nil, fmt.Errorf("Unable to find %s account", typ)
}

func purchase(product, currency string, amount float64) error {
	_, err := client.CreateOrder(&coinbasepro.Order{
		Side:      "buy",
		Type:      "market",
		ProductID: fmt.Sprintf("%s-%s", product, currency),
		Funds:     fmt.Sprintf("%.2f", amount),
	})
	return err
}

func deposit(currency string, amount float64, coinbase bool) error {
	if coinbase {
		// TODO
		return errors.New("coinbase deposit not yet implemented")
	}

	pm, err := getPaymentMethod(currency)
	if err != nil {
		return err
	}

	_, err = client.CreateDeposit(&coinbasepro.Deposit{
		Currency:        "USD",
		Amount:          fmt.Sprintf("%.2f", amount),
		PaymentMethodID: pm.ID,
	})
	return err
}

func getPaymentMethod(currency string) (*coinbasepro.PaymentMethod, error) {
	var pm coinbasepro.PaymentMethod
	pms, err := client.GetPaymentMethods()
	if err != nil {
		return nil, err
	}
	for _, p := range pms {
		if p.Currency == currency && p.Type == "ach_bank_account" {
			pm = p
			break
		}
	}

	return &pm, nil
}
