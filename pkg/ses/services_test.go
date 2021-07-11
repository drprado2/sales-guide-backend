package ses

import (
	"context"
	"testing"
)

type (
	TemplateData struct {
		Name           string `json:"name"`
		Favoriteanimal string `json:"favoriteanimal"`
	}
)

const (
	Sender       = "estudo.async.hb@gmail.com"
	Recipient    = "mariaaug222@gmail.com"
	CcRecipient  = "joaoaug222@gmail.com"
	CcRecipient2 = "drprado2@outlook.com"
	Subject      = "Amazon SES Test (AWS SDK for Go)"
	HtmlBody     = `
<html>
  <head>
	<meta charset="UTF-8">
  </head>
  <body>
	<h1>Amazon SES Test Email (AWS SDK for Go)</h1>
	<p>This email was sent with <a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the <a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>
  </body>
</html>
`
	TextBody = "This email was sent with Amazon SES using the AWS SDK for Go."
)

func TestSendRawEmail(t *testing.T) {
	ctx := context.Background()

	// To test with real AWS use the lines below and put valid keys, HTML is bugged in localstack
	//envs := configs.Get()
	//envs.AwsEndpoint = ""
	//envs.AwsRegion = "us-east-2"
	//envs.AwsAccessKey = "AKIA4TOL2AU6WVIAJABA"
	//envs.AwsSecretAccessKey = "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"

	if err := Setup(ctx); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	builder := NewEmailBuilder().
		WithSubject(Subject).
		ToDestinations(Recipient).
		WithCopyFor(CcRecipient, CcRecipient2).
		FromSender(Sender).
		AsRawBuilder().
		WithHtmlContent(HtmlBody).
		WithTextContent(TextBody)

	err := SendEmailSvc(ctx, builder)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
}

// Localstack has a bug with HTML, he sends as text and this cause not render HTML correctly, to test your template use an valid SES account in aws, SES sandbox need that you verify all e-mails FROM and TO
// to move out of sandbox contact support in the link https://docs.aws.amazon.com/pt_br/ses/latest/DeveloperGuide/request-production-access.html?icmpid=docs_ses_console
func TestSendTemplatedEmail(t *testing.T) {
	ctx := context.Background()
	// To test with real AWS use the lines below and put valid keys
	//envs := configs.Get()
	//envs.AwsEndpoint = ""
	//envs.AwsRegion = "us-east-2"
	//envs.AwsAccessKey = "AKIA4TOL2AU6WVIAJABA"
	//envs.AwsSecretAccessKey = "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"
	if err := Setup(ctx); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	data := TemplateData{
		Name:           "Adriano Oliveira",
		Favoriteanimal: "Cachorro",
	}

	builder := NewEmailBuilder().
		WithSubject(Subject).
		ToDestinations(Recipient).
		WithCopyFor(CcRecipient, CcRecipient2).
		FromSender(Sender).
		AsTemplatedBuilder().
		WithTemplate("TestTemplate", data)

	err := SendTemplatedEmailSvc(ctx, builder)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
}
