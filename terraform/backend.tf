terraform {
  backend "s3" {
    bucket = "swhite24-cbpro-buy-state"
    key    = "cbpro-weekly.terraform.tfstate"
    region = "us-east-1"
  }
}
