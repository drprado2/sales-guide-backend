package ses

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type (
	BaseBuilder struct {
		toDestinations []string
		ccDestinations []string
		subject        string
		sender         string
		charSet        string
	}

	EmailRawInputBuilder struct {
		BaseBuilder
		htmlMessage string
		textMessage string
	}

	EmailTemplateInputBuilder struct {
		BaseBuilder
		templateName   string
		templateParams interface{}
	}
)

func NewEmailBuilder() *BaseBuilder {
	res := new(BaseBuilder)
	res.toDestinations = make([]string, 0)
	res.ccDestinations = make([]string, 0)
	res.charSet = DefaultCharSet
	return res
}

func (eib *BaseBuilder) WithCharset(charset string) *BaseBuilder {
	eib.charSet = charset
	return eib
}

func (eib *BaseBuilder) ToDestinations(destinations ...string) *BaseBuilder {
	eib.toDestinations = destinations
	return eib
}

func (eib *BaseBuilder) WithCopyFor(destinations ...string) *BaseBuilder {
	eib.ccDestinations = destinations
	return eib
}

func (eib *BaseBuilder) WithSubject(subject string) *BaseBuilder {
	eib.subject = subject
	return eib
}

func (eib *BaseBuilder) FromSender(sender string) *BaseBuilder {
	eib.sender = sender
	return eib
}

func (eib *BaseBuilder) AsTemplatedBuilder() *EmailTemplateInputBuilder {
	res := new(EmailTemplateInputBuilder)
	res.BaseBuilder = *eib
	return res
}

func (eib *BaseBuilder) AsRawBuilder() *EmailRawInputBuilder {
	res := new(EmailRawInputBuilder)
	res.BaseBuilder = *eib
	return res
}

func (eib *EmailRawInputBuilder) WithHtmlContent(html string) *EmailRawInputBuilder {
	eib.htmlMessage = html
	return eib
}

func (eib *EmailRawInputBuilder) WithTextContent(text string) *EmailRawInputBuilder {
	eib.textMessage = text
	return eib
}

func (eib *EmailTemplateInputBuilder) WithTemplate(templateName string, templateData interface{}) *EmailTemplateInputBuilder {
	eib.templateName = templateName
	eib.templateParams = templateData
	return eib
}

func (eib *BaseBuilder) validate() []error {
	errors := make([]error, 0)
	if eib.subject == "" {
		errors = append(errors, RequiredSubjectError)
	}
	if eib.sender == "" {
		errors = append(errors, RequiredSenderError)
	}
	if len(eib.toDestinations) == 0 {
		errors = append(errors, RequiredDestinationsError)
	}
	if eib.charSet == "" {
		errors = append(errors, RequiredCharsetError)
	}
	return errors
}

func (eib *EmailRawInputBuilder) validate() error {
	errors := eib.BaseBuilder.validate()
	if eib.htmlMessage == "" && eib.textMessage == "" {
		errors = append(errors, RequiredContentError)
	}
	if len(errors) > 0 {
		return GroupErrors(errors...)
	}
	return nil
}

func (eib *EmailTemplateInputBuilder) validate() error {
	errors := eib.BaseBuilder.validate()
	if eib.templateParams == nil || eib.templateName == "" {
		errors = append(errors, RequiredTemplateError)
	}
	if len(errors) > 0 {
		return GroupErrors(errors...)
	}
	return nil
}

func (eib *EmailRawInputBuilder) Build() (*ses.SendEmailInput, error) {
	if err := eib.validate(); err != nil {
		return nil, err
	}

	result := new(ses.SendEmailInput)
	result.Destination = new(types.Destination)
	result.Destination.ToAddresses = eib.toDestinations
	result.Destination.CcAddresses = eib.ccDestinations
	result.Message = new(types.Message)
	result.Message.Body = new(types.Body)
	if eib.htmlMessage != "" {
		result.Message.Body.Html = &types.Content{
			Charset: aws.String(eib.charSet),
			Data:    aws.String(eib.htmlMessage),
		}
	}
	if eib.textMessage != "" {
		result.Message.Body.Text = &types.Content{
			Charset: aws.String(eib.charSet),
			Data:    aws.String(eib.textMessage),
		}
	}
	result.Message.Subject = &types.Content{
		Charset: aws.String(eib.charSet),
		Data:    aws.String(eib.subject),
	}
	result.Source = aws.String(eib.sender)

	return result, nil
}

func (eib *EmailTemplateInputBuilder) Build() (*ses.SendTemplatedEmailInput, error) {
	if err := eib.validate(); err != nil {
		return nil, err
	}

	jParams, err := json.Marshal(eib.templateParams)
	if err != nil {
		return nil, err
	}
	jsonParams := string(jParams)

	result := new(ses.SendTemplatedEmailInput)
	result.Destination = new(types.Destination)
	result.Destination.ToAddresses = eib.toDestinations
	result.Destination.CcAddresses = eib.ccDestinations
	result.Template = aws.String(eib.templateName)
	result.TemplateData = aws.String(jsonParams)
	result.Source = aws.String(eib.sender)

	return result, nil
}
