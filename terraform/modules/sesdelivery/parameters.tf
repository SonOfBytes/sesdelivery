resource "aws_ssm_parameter" "smtp_server" {
  name = "sesdelivery.smtp_server"
  type = "String"
  value = "${var.smtp_server}"
}

resource "aws_ssm_parameter" "sqs_notice_queue" {
  name = "sesdelivery.sqs_notice_queue"
  type = "String"
  value = "${aws_sqs_queue.sesdelivery.id}"
}