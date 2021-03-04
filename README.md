# cbpro-buy

`cbpro-buy` provides a simple CLI for automating purchases of BTC (or others if you so choose) from [Coinbase Pro](https://pro.coinbase.com).

## Prerequisites

- Obtain an API Key  
  https://docs.pro.coinbase.com/#authentication

## Installation

```sh
go get github.com/swhite24/cbpro-buy/cmd/cbpro-buy
```

## Usage

```sh
cbpro-buy purchases crypto from coinbase pro with auto deposit

Usage:
  cbpro-buy [flags]

Flags:
      --amount float        Amount of product to purchase (default 50)
      --autodeposit         Whether to auto deposit funds if current account is less than amount
      --coinbase            Whether to use coinbase account for funds instead of ACH
      --currency string     Currency to deposit / purchase with (USD, EUR, etc.) (default "USD")
  -h, --help                help for cbpro-buy
      --key string          Coinbase Pro API key
      --passphrase string   Coinbase Pro API key passphrase
      --product string      Product to purchase from coinbase pro (BTC, ETH, etc.) (default "BTC")
      --sandbox             Whether to use coinbase pro sandbox environment (will require different api key
      --secret string       Coinbase Pro API key secret
```

## License

See [LICENSE.txt](LICENSE.txt)
