variable "discord_webhook_url" {
  type = string
  description = "Discord webhook URL"
  sensitive = true
}

variable "region" {
  type = string
  description = "AWS region"
  default = "ap-northeast-1"
}
