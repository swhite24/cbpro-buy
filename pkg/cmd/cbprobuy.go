package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/preichenberger/go-coinbasepro/v2"
	"github.com/spf13/cobra"
	"github.com/swhite24/cbpro-buy/pkg/book"
	"github.com/swhite24/cbpro-buy/pkg/config"
	"github.com/swhite24/cbpro-buy/pkg/purchase"

	basisconfig "github.com/swhite24/cbpro-cost-basis/pkg/config"
	"github.com/swhite24/cbpro-cost-basis/pkg/costbasis"
)

var (
	// CBProBuyCmd purchases crypto from coinbase pro with auto deposit
	CBProBuyCmd *cobra.Command
)

func init() {
	var key, passphrase, secret string
	var currency, product string
	var useCoinbase, useSandbox, autoDeposit bool
	var amount float64
	var useBasis bool
	var basisWindowStart, basisMultiplier float64

	CBProBuyCmd = &cobra.Command{
		Use:   "cbpro-buy",
		Short: "cbpro-buy purchases crypto from coinbase pro with auto deposit",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.InitializeConfig(cmd.Flags())
			client := initializeClient(cfg)
			fmt.Println(cfg)

			// Determine basis if configured
			if cfg.UseBasis {
				fmt.Printf("use-basis configured; getting average purchase price over last %f days.\n", cfg.BasisWindowStart)
				// Calculate average cost
				c := &basisconfig.Config{
					Key:        cfg.Key,
					Passphrase: cfg.Passphrase,
					Secret:     cfg.Secret,
					Product:    fmt.Sprintf("%s-%s", cfg.Product, cfg.Currency),
					StartDate:  time.Now().Add(time.Duration(cfg.BasisWindowStart*-24) * time.Hour),
					EndDate:    time.Now(),
				}
				info, err := costbasis.Calculate(client, c)
				if err != nil {
					fmt.Println("failed to calculate basis")
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Printf("Average purchase price: %s\n", info.AverageCost)

				// Get current price
				average, _ := strconv.ParseFloat(info.AverageCost, 64)
				price, err := book.GetPrice(client, fmt.Sprintf("%s-%s", cfg.Product, cfg.Currency))
				if err != nil {
					fmt.Println("failed to calculate current price")
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Printf("Current book price: %.2f\n", price)

				// Update purchase amount if current price is less than average cost
				if price < average {
					fmt.Println("Current price less than average price.")
					fmt.Printf("Adjusting buy amount from %.2f to %.2f\n", cfg.Amount, cfg.Amount*cfg.BasisMultiplier)
					cfg.Amount = cfg.Amount * cfg.BasisMultiplier
				}
			}

			err := purchase.InitiatePurchase(client, cfg)
			if err != nil {
				fmt.Println("failed to purchase")
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("Purchase successful!")
		},
	}

	CBProBuyCmd.Flags().StringVar(&key, "key", "", "Coinbase Pro API key")
	CBProBuyCmd.Flags().StringVar(&passphrase, "passphrase", "", "Coinbase Pro API key passphrase")
	CBProBuyCmd.Flags().StringVar(&secret, "secret", "", "Coinbase Pro API key secret")

	CBProBuyCmd.Flags().StringVar(&currency, "currency", "USD", "Currency to deposit / purchase with (USD, EUR, etc.)")
	CBProBuyCmd.Flags().StringVar(&product, "product", "BTC", "Product to purchase from coinbase pro (BTC, ETH, etc.)")
	CBProBuyCmd.Flags().BoolVar(&useCoinbase, "coinbase", false, "Whether to use coinbase account for funds instead of ACH")
	CBProBuyCmd.Flags().BoolVar(&useSandbox, "sandbox", false, "Whether to use coinbase pro sandbox environment (will require different api key")
	CBProBuyCmd.Flags().BoolVar(&autoDeposit, "autodeposit", false, "Whether to auto deposit funds if current account is less than amount")
	CBProBuyCmd.Flags().Float64Var(&amount, "amount", 50, "Amount of product to purchase")
	CBProBuyCmd.Flags().BoolVar(&useBasis, "use-basis", false, "Whether to adjust purchase amount if current price is below average cost over time window")
	CBProBuyCmd.Flags().Float64Var(&basisWindowStart, "basis-window-start", 30, "Mumber of days in the past to for beginning of basis window")
	CBProBuyCmd.Flags().Float64Var(&basisMultiplier, "basis-multiplier", 1.5, "Scale to apply to purchase amount if current price is less than average cost")
}

func initializeClient(cfg *config.Config) *coinbasepro.Client {
	client := coinbasepro.NewClient()
	client.UpdateConfig(&coinbasepro.ClientConfig{
		BaseURL:    cfg.BaseURL,
		Key:        cfg.Key,
		Passphrase: cfg.Passphrase,
		Secret:     cfg.Secret,
	})
	return client
}
