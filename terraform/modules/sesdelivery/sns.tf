resource "aws_sns_topic" "sesdelivery_notices" {
  name = "sesdelivery_notices"
}

resource "aws_sns_topic_subscription" "sesdelivery_sqs" {
  endpoint = "${aws_sqs_queue.sesdelivery.arn}"
  protocol = "sqs"
  topic_arn = "${aws_sns_topic.sesdelivery_notices.arn}"
  raw_message_delivery = true
}

resource "aws_sns_topic_subscription" "sesdelivery_lambda" {
  endpoint = "${aws_lambda_function.sesdelivery_lambda.arn}"
  protocol = "lambda"
  topic_arn = "${aws_sns_topic.sesdelivery_notices.arn}"
}