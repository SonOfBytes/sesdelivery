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
      "logs:*",
    ]

    resources = [
      "arn:aws:logs:*:*:*",
    ]
  }
}

resource "aws_iam_role" "sesdelivery_lambda" {
  name = "sesdelivery_lambda"
  assume_role_policy = "${data.aws_iam_policy_document.sesdelivery_lambda_assume_role.json}"
}

resource "aws_iam_policy" "sesdelivery_lambda" {
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

  environment {
    variables = {
      LAMBDA_FUNCTION_NAME = "sesdelivery"
    }
  }
}