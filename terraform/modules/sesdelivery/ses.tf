resource "aws_ses_receipt_rule" "sesdelivery" {
  name = "sesdelivery"
  rule_set_name = "default-rule-set"
  recipients    = "${sort(var.ses_recipients)}"
  enabled       = true
  scan_enabled  = true

  add_header_action {
    header_name  = "X-AWS-PROCESSOR"
    header_value = "Processed by SES Delivery"
    position = 0
  }

  s3_action {
    bucket_name = "${aws_s3_bucket.sesdelivery.bucket}"
    kms_key_arn = "${aws_kms_alias.sesdelivery_alias.arn}"
    topic_arn = "${aws_sns_topic.sesdelivery_notices.arn}"
    position = 1
  }
}