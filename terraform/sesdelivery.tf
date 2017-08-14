module "sesdelivery" {
  source = "modules/sesdelivery"
  smtp_server = "${var.smtp_server}"
  ses_recipients = "${var.ses_recipients}"
}