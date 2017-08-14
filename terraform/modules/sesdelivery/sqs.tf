resource "aws_sqs_queue" "sesdelivery" {
  name = "sesdelivery"
  delay_seconds = 0
  message_retention_seconds = 1209600
  receive_wait_time_seconds = 0
  policy = "${data.aws_iam_policy_document.sesdelivery_sqs.json}"
}

data "aws_iam_policy_document" "sesdelivery_sqs" {
  statement {
    effect = "Allow"

    actions = [
      "SQS:SendMessage",
    ]

    resources = [
      "arn:aws:sqs:${var.aws_region}:${data.aws_caller_identity.current.account_id}:sesdelivery",
    ]

    principals {
      identifiers = ["*"]
      type = "AWS"
    }

    condition {
      test = "ArnEquals"
      values = [
        "${aws_sns_topic.sesdelivery_notices.arn}"
      ]
      variable = "aws:SourceArn"
    }
  }
}