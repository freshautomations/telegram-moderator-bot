provider "aws" {
  region = "us-east-1"
}

#terraform {
#  backend "s3" {}
#}

data "aws_caller_identity" "current" {}

resource "aws_iam_role" "tmb" {
  name        = "tmb-lambda-${var.ENVIRONMENT}"
  description = "Telegram Moderator Bot Lambda function permissions"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "tmb" {
  name        = "tmb-policy-${var.ENVIRONMENT}"
  path        = "/"
  description = "IAM policy tmb lambda logging and DB access"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "${aws_cloudwatch_log_group.tmb.arn}",
      "Effect": "Allow"
    },
    {
      "Action": "dynamodb:*",
      "Resource": "arn:aws:dynamodb:us-east-1:${data.aws_caller_identity.current.account_id}:table/tmb-${var.ENVIRONMENT}-users",
      "Effect": "Allow"
    },
    {
      "Action": "dynamodb:ListTables",
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource aws_wafregional_ipset tmb {
  name = "tmb-${var.ENVIRONMENT}"

  ip_set_descriptor {
    type  = "IPV4"
    value = "149.154.167.197/32"
  }

  ip_set_descriptor {
    type  = "IPV4"
    value = "149.154.167.198/31"
  }

  ip_set_descriptor {
    type  = "IPV4"
    value = "149.154.167.200/29"
  }

  ip_set_descriptor {
    type  = "IPV4"
    value = "149.154.167.208/28"
  }

  ip_set_descriptor {
    type  = "IPV4"
    value = "149.154.167.224/29"
  }

  ip_set_descriptor {
    type  = "IPV4"
    value = "149.154.167.232/31"
  }
}

resource aws_wafregional_rule tmb {
  name        = "tmb-${var.ENVIRONMENT}"
  metric_name = "tmb${var.ENVIRONMENT}wafrule"

  predicate {
    data_id = "${aws_wafregional_ipset.tmb.id}"
    negated = false
    type    = "IPMatch"
  }
}

resource aws_wafregional_web_acl tmb {
  name        = "tmb-${var.ENVIRONMENT}"
  metric_name = "tmb${var.ENVIRONMENT}wafacl"

  default_action {
    type = "BLOCK"
  }

  rule {
    action {
      type = "ALLOW"
    }

    priority = 1
    rule_id  = "${aws_wafregional_rule.tmb.id}"
    type     = "REGULAR"
  }
}

// Terraform 0.11.11 does not report back the stage's ARN, only the stage's execution ARN.
resource aws_wafregional_web_acl_association tmb {
  resource_arn = "arn:aws:apigateway:us-east-1::/restapis/${aws_api_gateway_rest_api.tmb.id}/stages/${var.LAMBDA_SECRET}"
  web_acl_id   = "${aws_wafregional_web_acl.tmb.id}"
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = "${aws_iam_role.tmb.name}"
  policy_arn = "${aws_iam_policy.tmb.arn}"
}

resource aws_dynamodb_table tmb {
  name = "tmb-${var.ENVIRONMENT}-users"
  hash_key = "username"
  read_capacity = 5
  write_capacity = 5
  attribute {
    name = "username"
    type = "S"
  }
  ttl {
    //TTL disabled because we can't lose moderator user IDs if they don't send messages for a while.
    //Todo: open an issue
    enabled = false
    attribute_name = "ttl"
  }
}

resource aws_lambda_function tmb {
  function_name = "tmb-${var.ENVIRONMENT}"
  filename      = "../../build/tmb.zip"
  description   = "Telegram Moderator Bot ${var.ENVIRONMENT} Lambda function"

  handler = "build/tmb"
  runtime = "go1.x"

  role    = "${aws_iam_role.tmb.arn}"
  timeout = 5

  source_code_hash = "${base64sha256(file("../../build/tmb.zip"))}"
  publish          = true

  environment {
    variables {
      "ENVIRONMENT"   = "${var.ENVIRONMENT}"
      "TIMEOUT"       = "${var.LAMBDA_TIMEOUT}"
      "AWSREGION"     = "us-east-1"
      "TELEGRAMTOKEN" = "${var.TELEGRAM_TOKEN}"
    }
  }

  tags {
    "Name"        = "tmb-${var.ENVIRONMENT}"
    "Environment" = "${var.ENVIRONMENT}"
  }
}

resource aws_cloudwatch_log_group tmb {
  name              = "/aws/lambda/${aws_lambda_function.tmb.function_name}"
  retention_in_days = 30
}

resource aws_api_gateway_rest_api tmb {
  name        = "tmb-${var.ENVIRONMENT}"
  description = "Telegram Moderator Bot ${var.ENVIRONMENT} API"

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource aws_api_gateway_method tmb_root {
  rest_api_id   = "${aws_api_gateway_rest_api.tmb.id}"
  resource_id   = "${aws_api_gateway_rest_api.tmb.root_resource_id}"
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_method_response" "tmb_root" {
  depends_on  = ["aws_api_gateway_method.tmb_root"]
  rest_api_id = "${aws_api_gateway_rest_api.tmb.id}"
  resource_id = "${aws_api_gateway_rest_api.tmb.root_resource_id}"
  http_method = "POST"
  status_code = "200"

  response_models {
    "application/json" = "Empty"
  }
}

resource aws_api_gateway_integration tmb_root {
  depends_on              = ["aws_api_gateway_method.tmb_root"]
  rest_api_id             = "${aws_api_gateway_rest_api.tmb.id}"
  resource_id             = "${aws_api_gateway_rest_api.tmb.root_resource_id}"
  http_method             = "POST"
  type                    = "AWS_PROXY"
  integration_http_method = "POST"
  content_handling        = "CONVERT_TO_TEXT"
  uri                     = "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/${aws_lambda_function.tmb.arn}/invocations"
}

resource aws_api_gateway_integration_response tmb_root {
  depends_on  = ["aws_api_gateway_integration.tmb_root"]
  rest_api_id = "${aws_api_gateway_rest_api.tmb.id}"
  resource_id = "${aws_api_gateway_rest_api.tmb.root_resource_id}"
  http_method = "POST"
  status_code = "200"

  response_templates {
    "application/json" = "Empty"
  }
}

resource "aws_lambda_permission" "tmb_root" {
  function_name = "${aws_lambda_function.tmb.function_name}"
  statement_id  = "apigateway-perm"
  action        = "lambda:InvokeFunction"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "arn:aws:execute-api:us-east-1:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.tmb.id}/${var.LAMBDA_SECRET}/POST/"
}

resource "aws_api_gateway_deployment" "tmb" {
  depends_on        = ["aws_api_gateway_integration.tmb_root"]
  rest_api_id       = "${aws_api_gateway_rest_api.tmb.id}"
  stage_name        = "${var.LAMBDA_SECRET}"
  stage_description = "${var.ENVIRONMENT} deployment with Lambda secret"
  description       = "Automated ${var.ENVIRONMENT} deployment"
}
