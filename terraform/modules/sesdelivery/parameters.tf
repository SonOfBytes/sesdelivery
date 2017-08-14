resource "aws_ssm_parameter" "smtp_server" {
  name  = "SMTPServer"
  type = "String"
  value = "${var.smtp_server}"
}

resource "aws_ssm_parameter" "sqs_notice_queue" {
  name  = "SQSNoticeQueue"
  type = "String"
  value = "${aws_sqs_queue.sesdelivery.id}"
}