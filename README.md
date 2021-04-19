# cbpro-buy

`cbpro-buy` provides a simple CLI for automating purchases of BTC (or others if you so choose) from [Coinbase Pro](https://pro.coinbase.com).

## Prerequisites

- Verified Coinbase Pro account with instant deposits
- Obtain an API Key  
  https://docs.pro.coinbase.com/#authentication

## Installation

```sh
go get github.com/swhite24/cbpro-buy/cmd/cbpro-buy
```

## Features

- Purchase any product on Coinbase Pro in any currency pair (that's available to your locale / account).
- Automatically deposit funds to conduct the purchase if your USD / other currency account is unable to fulfill.
- Automatically wait for deposit to clear before attempting to issue a market buy.
- Allow for weighing purchases if current price falls below average price over configured amount of time.

## Usage

```sh
cbpro-buy purchases crypto from coinbase pro with auto deposit

Usage:
  cbpro-buy [flags]

Flags:
      --amount float               Amount of product to purchase (default 50)
      --autodeposit                Whether to auto deposit funds if current account is less than amount
      --basis-multiplier float     Scale to apply to purchase amount if current price is less than average cost (default 1.5)
      --basis-window-start float   Mumber of days in the past to for beginning of basis window (default 30)
      --coinbase                   Whether to use coinbase account for funds instead of ACH
      --currency string            Currency to deposit / purchase with (USD, EUR, etc.) (default "USD")
  -h, --help                       help for cbpro-buy
      --key string                 Coinbase Pro API key
      --passphrase string          Coinbase Pro API key passphrase
      --product string             Product to purchase from coinbase pro (BTC, ETH, etc.) (default "BTC")
      --sandbox                    Whether to use coinbase pro sandbox environment (will require different api key
      --secret string              Coinbase Pro API key secret
      --use-basis                  Whether to adjust purchase amount if current price is below average cost over time window
```

## Example

```sh
$ export CBPRO_BUY_SECRET="..."
$ export CBPRO_BUY_KEY="..."
$ export CBPRO_BUY_PASSPHRASE="..."
$ cbpro-buy --amount 10 --autodeposit --use-basis
Use basis configured. Getting average purchase price over last 30 days.
Average purchase price: 56580.94
Current book price: 54953.52
Current price less than average price.
Adjusting buy amount from 10.00 to 15.00
Fetching current funding account status.
Success.  Available balance: 0.330014
Available balance is less than requested purchase: 15.00
Initiating deposit of 15.00 USD
Waiting for deposit to be available in account.
Checking available balance: 0.33
Checking available balance: 15.33
Initiating purchase of 15.00 USD worth of BTC
Purchase successful!
```

## Lambda Usage

This repo also contains an example deployment using [AWS Lambda](https://aws.amazon.com/lambda/) executed on a schedule to provide a means to "dollar cost average" with scheduled buys regardless of price. This is handled with [terraform](https://www.terraform.io/). See [the terraform directory](terraform) for details.

```sh
terraform init
terraform apply
```

## License

See [LICENSE.txt](LICENSE.txt)
