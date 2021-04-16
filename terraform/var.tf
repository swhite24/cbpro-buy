variable "region" { default = "us-east-1" }
variable "function_name" { default = "cbpro-buy-weekly" }
variable "lambda_schedule_expression" { default = "cron(0 5 * * ? *)" }
variable "executable" { default = "cbpro-buy-lambda" }
variable "archive" { default = "cbpro-buy-lambda.zip" }
# set to 1 or null
variable "auto_deposit" { default = 1 }
variable "amount" { default = 50 }
variable "currency" { default = "USD" }
variable "product" { default = "BTC" }
variable "use_basis" { default = true }
variable "basis_window_start" { default = 60 }
variable "cbpro_key" {
  type      = string
  sensitive = true
}
variable "cbpro_passphrase" {
  type      = string
  sensitive = true
}
variable "cbpro_secret" {
  type      = string
  sensitive = true
}
