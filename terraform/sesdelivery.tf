module "sesdelivery" {
  source = "modules/sesdelivery"
  smtp_server = "${var.smtp_server}"
  sqs_notice_queue = "${var.sqs_notice_queue}"
}