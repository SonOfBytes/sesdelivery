variable "smtp_server" {
  type = "string"
}

variable "ses_recipients" {
  default = []
  type = "list"
}