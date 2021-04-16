data "archive_file" "zip" {
  type        = "zip"
  source_file = "../bin/${var.executable}"
  output_path = var.archive
}

resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_lambda_function" "cbpro_buy" {
  filename         = var.archive
  function_name    = var.function_name
  role             = aws_iam_role.iam_for_lambda.arn
  handler          = var.executable
  source_code_hash = data.archive_file.zip.output_base64sha256
  timeout          = 30
  runtime          = "go1.x"

  environment {
    variables = {
      CBPRO_BUY_KEY         = var.cbpro_key
      CBPRO_BUY_PASSPHRASE  = var.cbpro_passphrase
      CBPRO_BUY_SECRET      = var.cbpro_secret
      CBPRO_BUY_CURRENCY    = var.currency
      CBPRO_BUY_PRODUCT     = var.product
      CBPRO_BUY_AMOUNT      = var.amount
      CBPRO_BUY_AUTODEPOSIT = var.auto_deposit
      CBPRO_BUY_USE_BASIS   = var.use_basis
    }
  }
}

resource "aws_iam_policy" "lambda_logging" {
  name        = "lambda_logging"
  path        = "/"
  description = "IAM policy for logging from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = aws_iam_policy.lambda_logging.arn
}

resource "aws_cloudwatch_event_rule" "event_rule" {
  schedule_expression = var.lambda_schedule_expression
}

resource "aws_cloudwatch_event_target" "event_target" {
  rule = aws_cloudwatch_event_rule.event_rule.name
  arn  = aws_lambda_function.cbpro_buy.arn
}

resource "aws_lambda_permission" "cloudwatch_permission" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = var.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.event_rule.arn
}
