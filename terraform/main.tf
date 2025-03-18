terraform {
  required_version = "= 1.11.1"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "= 5.91.0"
    }    
  }
}

provider "aws" {
  region = "ap-northeast-1"
}

# Lambda用IAMロールの作成
resource "aws_iam_role" "lambda_role" {
  name = "weather_lambda_role"
  assume_role_policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [{
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }]
  })
}

# CloudWatch Logsへのアクセス権限付与
resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Lambda関数の定義
resource "aws_lambda_function" "weather_forecast" {
  function_name = "weather_forecast"
  role          = aws_iam_role.lambda_role.arn
  handler          = "bootstrap"
  runtime          = "provided.al2"
  filename      = "../lambda/function.zip"
  source_code_hash = filebase64sha256("../lambda/function.zip")
  
  environment {
    variables = {
      DISCORD_WEBHOOK_URL = var.discord_webhook_url
    }
  }
}

# EventBridgeルール（毎日7時にトリガー）
resource "aws_cloudwatch_event_rule" "daily_trigger" {
  name                = "daily_weather_trigger"
  schedule_expression = "cron(0 22 * * ? *)"  # 22:00 UTC = 7:00 JST
}

# EventBridgeルールのターゲットとしてLambdaを設定
resource "aws_cloudwatch_event_target" "lambda_target" {
  rule      = aws_cloudwatch_event_rule.daily_trigger.name
  target_id = "weather_forecast_target"
  arn       = aws_lambda_function.weather_forecast.arn
}

# Lambda関数にEventBridgeからの呼び出しを許可
resource "aws_lambda_permission" "allow_eventbridge" {
  statement_id  = "AllowExecutionFromEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.weather_forecast.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.daily_trigger.arn
}
