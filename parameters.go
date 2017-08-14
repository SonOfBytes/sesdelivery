package sesdelivery

import (
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"fmt"
)

type Parameters struct {
	session *session.Session
	smtpServer string
	sqsNoticeQueue string
}

func NewParameters() (p *Parameters, err error) {
	p = &Parameters{
		session: session.New(&aws.Config{}),
	}

	return p, nil
}

func (p *Parameters) GetSMTPServer() (smtpServer string, err error) {
	if p.smtpServer == "" {
		p.smtpServer, err = p.getParameter("SMTPServer", false)
	}
	return p.smtpServer, err
}

func (p *Parameters) GetSQSNoticeQueue() (sqsNoticeQueue string, err error) {
	if p.sqsNoticeQueue == "" {
		p.sqsNoticeQueue, err = p.getParameter("SQSNoticeQueue", false)
	}
	return p.sqsNoticeQueue, err
}

func (p *Parameters) getParameter(v string, encrypted bool) (s string, err error) {
	client := ssm.New(p.session)
	param_input := &ssm.GetParameterInput{
		Name: aws.String(v),
		WithDecryption: aws.Bool(encrypted),
	}

	param_output, err := client.GetParameter(param_input)
	if err != nil {
		return "", fmt.Errorf("getParameter(%s): %s", v, err)
	}

	return *param_output.Parameter.Value, nil
}