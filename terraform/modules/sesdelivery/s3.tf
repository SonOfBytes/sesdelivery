resource "aws_s3_bucket" "sesdelivery" {
  bucket = "sesdelivery-${data.aws_caller_identity.current.account_id}"
  acl = "private"
  policy = "${data.aws_iam_policy_document.sesdelivery_s3.json}"
}

data "aws_iam_policy_document" "sesdelivery_s3" {
  statement {
    effect = "Allow"

    principals {
      identifiers = [
        "ses.amazonaws.com",
      ]
      type = "Service"
    }

    actions = [
      "s3:PutObject",
    ]

    resources = [
      "arn:aws:s3:::sesdelivery-${data.aws_caller_identity.current.account_id}/*"
    ]

    condition {
      test = "StringEquals"
      values = [
        "${data.aws_caller_identity.current.account_id}",
      ]
      variable = "aws:Referer"
    }
  }
}