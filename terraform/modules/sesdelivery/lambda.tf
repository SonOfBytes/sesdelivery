data "aws_iam_policy_document" "sesdelivery_lambda_assume_role" {
  statement {
    effect = "Allow"

    actions = [
      "sts:AssumeRole",
    ]

    principals = {
      type = "Service"
      identifiers = [
        "lambda.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "sesdelivery_lambda_role" {
  statement {
    effect = "Allow"

    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = [
      "${aws_cloudwatch_log_group.lambda_logs.arn}",
    ]
  }
  statement {
    effect = "Allow"

    actions = [
      "ssm:GetParameter",
    ]

    resources = [
      "arn:aws:ssm:${var.aws_region}:${data.aws_caller_identity.current.account_id}:parameter/${aws_ssm_parameter.smtp_server.name}",
      "arn:aws:ssm:${var.aws_region}:${data.aws_caller_identity.current.account_id}:parameter/${aws_ssm_parameter.sqs_notice_queue.name}",
    ]
  }
  statement {
    effect = "Allow"
    actions = [
      "sqs:ReceiveMessage",
      "sqs:DeleteMessage"
    ]
    resources = [
      "${aws_sqs_queue.sesdelivery.arn}"
    ]
  }
  statement {
    effect = "Allow"
    actions = [
      "s3:GetObject",
      "s3:PutObject",
      "s3:DeleteObject",
    ]
    resources = [
      "${aws_s3_bucket.sesdelivery.arn}/*"]
  }
  statement {
    effect = "Allow"
    actions = [
      "kms:Decrypt",
    ]
    resources = [
      "${aws_kms_key.sesdelivery.arn}"]
  }
}

resource "aws_iam_role" "sesdelivery_lambda" {
  name = "sesdelivery_lambda"
  assume_role_policy = "${data.aws_iam_policy_document.sesdelivery_lambda_assume_role.json}"
}

resource "aws_iam_policy" "sesdelivery_lambda" {
  name = "sesdelivery_lambda_policy"
  policy = "${data.aws_iam_policy_document.sesdelivery_lambda_role.json}"
}

resource "aws_iam_role_policy_attachment" "sesdelivery_lambda_role" {
  policy_arn = "${aws_iam_policy.sesdelivery_lambda.arn}"
  role = "${aws_iam_role.sesdelivery_lambda.name}"
}

resource "aws_lambda_function" "sesdelivery_lambda" {
  filename         = "handler.zip"
  function_name    = "sesdelivery"
  role             = "${aws_iam_role.sesdelivery_lambda.arn}"
  handler          = "handler.Handle"
  source_code_hash = "${base64sha256(file("handler.zip"))}"
  runtime          = "python2.7"
  timeout = 120

  environment {
    variables = {
      LAMBDA_FUNCTION_NAME = "sesdelivery"
    }
  }
}

resource "aws_cloudwatch_log_group" "lambda_logs" {
  name = "/aws/lambda/${aws_lambda_function.sesdelivery_lambda.function_name}"
  retention_in_days = "30"
}