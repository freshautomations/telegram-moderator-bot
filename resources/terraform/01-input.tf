variable ENVIRONMENT {
  type = "string"
  default = "staging"
  description = "Prefix for the AWS environment for better identification. Do not reuse it unless with terragrunt."
}

variable LAMBDA_SECRET {
  type = "string"
  default = "staging"
  description = "Token that identifies Telegram for the Lambda function"
}

variable LAMBDA_TIMEOUT {
  type = "string"
  default = 15
  description = "Lambda function timeout value"
}

variable "TELEGRAM_TOKEN" {
  type = "string"
  description = "Telegram Token received from @BotFather"
}
