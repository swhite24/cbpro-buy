package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/swhite24/dca-bot/pkg/config"
	"github.com/swhite24/dca-bot/pkg/purchase"
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

	CBProBuyCmd = &cobra.Command{
		Use:   "cbpro-buy",
		Short: "cbpro-buy purchases crypto from coinbase pro with auto deposit",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.InitializeConfig(cmd.Flags())
			err := purchase.InitiatePurchase(cfg)
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
}
