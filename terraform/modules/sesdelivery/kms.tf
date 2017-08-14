resource "aws_kms_key" "sesdelivery" {
  description = "sesdelivery KMS Key for symmetric encryption and decryption"
  deletion_window_in_days = 10
  policy = "${data.aws_iam_policy_document.sesdelivery_kms.json}"
}

resource "aws_kms_alias" "sesdelivery_alias" {
  name = "alias/sesdelivery"
  target_key_id = "${aws_kms_key.sesdelivery.key_id}"
}

data "aws_iam_policy_document" "sesdelivery_kms" {
  statement {
    effect = "Allow"
    principals {
      identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"]
      type = "AWS"
    }

    actions = [
      "kms:*",
    ]
    resources = ["*"]
  }
  statement {
    effect = "Allow"
    principals {
      identifiers = [
        "ses.amazonaws.com",
      ]
      type = "Service"
    }
    actions = [
      "kms:Encrypt",
      "kms:GenerateDataKey*",
    ]
    resources = ["*"]
    condition {
      test = "Null"
      values = [
        "false"
      ]
      variable = "kms:EncryptionContext:aws:ses:message-id"
    }
    condition {
      test = "Null"
      values = [
       "false"
      ]
      variable = "kms:EncryptionContext:aws:ses:rule-name"
    }
    condition {
      test = "StringEquals"
      values = [
        "${data.aws_caller_identity.current.account_id}",
      ]
      variable = "kms:EncryptionContext:aws:ses:source-account"
    }
  }
}