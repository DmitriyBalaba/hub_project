package smtp

import (
	"bitbucket.org/jazzserve/webapi/utils"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/smtp"
	"strconv"
	"strings"
)

const contentTypeTextPlain = "text/plain"

type Email struct {
	*Config
	Templates *template.Template
}

type EmailBodyPart struct {
	ContentType             string
	ContentTransferEncoding string
	ContentDisposition      string
	Body                    string
	ContentID               string
}

func NewEmail(conf *Config) *Email {
	if conf == nil {
		panic("can't create new smtp email sender with nil config")
	}
	return &Email{Config: conf}
}

type attachment struct {
	body     string
	fileName string
}

func (e *Email) AddAttachment(att, fileName string) *attachment {
	return &attachment{
		body:     att,
		fileName: fileName,
	}
}

func (e *Email) SendHTMLMessage(toList []string, subject string, message string, attachment *attachment) error {
	parts := []EmailBodyPart{e.GetHTMLPart(message)}
	if attachment != nil {
		parts = append(parts, e.GetAttachmentPart(attachment))
	}
	return e.SendMultipartEmail(toList, subject, parts)
}

func (e *Email) HTMLTemplateToString(tmplName string, tmplData interface{}) (tmpl string, err error) {
	if e.Templates == nil {
		return "", errors.New("no templates were loaded during application configuration")
	}
	t := e.Templates.Lookup(tmplName)
	if t == nil {
		return "", fmt.Errorf("template '%s' not found", tmplName)
	}
	var b bytes.Buffer
	if err := t.Execute(&b, tmplData); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (e *Email) GetHTMLPart(body string) (part EmailBodyPart) {
	part.ContentType = `text/html; charset="UTF-8"`
	part.ContentTransferEncoding = "base64"
	part.Body = base64.StdEncoding.EncodeToString([]byte(body))
	return part
}

func getAttachmentPart(attachment *attachment, contentType string, encodeToBase64 bool) (part EmailBodyPart) {
	if attachment == nil {
		return EmailBodyPart{}
	}
	part.ContentType = contentType + `; name="` + attachment.fileName + `"`
	part.ContentTransferEncoding = "base64"
	part.ContentDisposition = `attachment; filename="` + attachment.fileName + `"`
	if encodeToBase64 == true {
		part.Body = base64.StdEncoding.EncodeToString([]byte(attachment.body))
	} else {
		part.Body = attachment.body
	}
	var err error
	part.ContentID, err = utils.GenerateRandomString(10)
	if err != nil {
		log.Error().Msgf("failed to generate content id for attachment %s", attachment.fileName)
	}
	return
}

func (e *Email) GetAttachmentPart(attachment *attachment) (part EmailBodyPart) {
	return getAttachmentPart(attachment, contentTypeTextPlain, true)
}

func (e *Email) GetAttachmentPartRaw(attachment *attachment) (part EmailBodyPart) {
	return getAttachmentPart(attachment, contentTypeTextPlain, false)
}

func (e *Email) GetAttachmentPartCT(attachment *attachment, contentType string) (part EmailBodyPart) {
	return getAttachmentPart(attachment, contentType, true)
}

func (e *Email) GetAttachmentPartRawCT(attachment *attachment, contentType string) (part EmailBodyPart) {
	return getAttachmentPart(attachment, contentType, false)
}

func (e *Email) SendMultipartEmail(to []string, subject string, emailBodyParts []EmailBodyPart) (err error) {
	msgId, err := utils.GenerateRandomString(32)
	if err != nil {
		return
	}

	header := make(map[string]string)
	header["MIME-Version"] = "1.0"
	header["From"] = e.SenderName + "<" + e.From + ">"
	header["To"] = strings.Join(to, ",")
	header["Reply-to"] = e.From
	header["Subject"] = subject
	header["Content-Type"] = "multipart/alternative;\r\n  boundary=" + `"` + msgId + `"`

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += "\r\n"

	for _, part := range emailBodyParts {
		message += "--" + msgId
		message += "\r\n"
		if part.ContentType != "" {
			message += "Content-Type: " + part.ContentType
			message += "\r\n"
		}
		if part.ContentTransferEncoding != "" {
			message += "Content-Transfer-Encoding: " + part.ContentTransferEncoding
			message += "\r\n"
		}
		if part.ContentDisposition != "" {
			message += "Content-Disposition: " + part.ContentDisposition
			message += "\r\n"
		}

		message += "\r\n"
		message += part.Body
		message += "\r\n"
	}

	message += "--" + msgId + "--"

	return e.sendMail(to, []byte(message))

}

type SendTo struct {
	To  []string
	Cc  []string
	Bcc []string
}

func (e *Email) SendMultipartEmailEx(sendTo SendTo, subject string, emailBodyParts []EmailBodyPart) (err error) {
	if len(sendTo.To) == 0 {
		return errors.New("empty send to address")
	}
	msgId, err := utils.GenerateRandomString(32)
	if err != nil {
		return
	}

	header := make(map[string]string)
	header["MIME-Version"] = "1.0"
	header["From"] = e.SenderName + "<" + e.From + ">"
	header["To"] = strings.Join(sendTo.To, ",")
	if len(sendTo.Cc) != 0 {
		header["Cc"] = strings.Join(sendTo.Cc, ",")
	}
	if len(sendTo.Bcc) != 0 {
		header["Bcc"] = strings.Join(sendTo.Bcc, ",")
	}
	header["Reply-to"] = e.From
	header["Subject"] = subject
	header["Content-Type"] = "multipart/mixed;\r\n  boundary=" + `"` + msgId + `"`

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += "\r\n"

	for _, part := range emailBodyParts {
		message += "--" + msgId
		message += "\r\n"
		if part.ContentType != "" {
			message += "Content-Type: " + part.ContentType
			message += "\r\n"
		}
		if part.ContentTransferEncoding != "" {
			message += "Content-Transfer-Encoding: " + part.ContentTransferEncoding
			message += "\r\n"
		}
		if part.ContentDisposition != "" {
			message += "Content-Disposition: " + part.ContentDisposition
			message += "\r\n"
		}
		if part.ContentID != "" {
			message += "Content-ID: <" + part.ContentID + ">"
			message += "\r\n"
		}

		message += "\r\n"
		message += part.Body
		message += "\r\n"
	}

	message += "--" + msgId + "--"

	return e.sendMailEx(sendTo, []byte(message))
}

func (e *Email) sendMailEx(sendTo SendTo, msg []byte) error {
	allRecipients := append(sendTo.To, sendTo.Cc...)
	allRecipients = append(allRecipients, sendTo.Bcc...)
	return e.sendMail(allRecipients, msg)
}

func (e *Email) sendMail(to []string, msg []byte) error {
	auth := smtp.PlainAuth("", e.UserName, e.Password, e.Host)
	if e.Port == 465 {
		return sslEmail(e.Host, e.Port, auth, e.From, to, msg)
	}
	return smtp.SendMail(e.Host+":"+strconv.Itoa(e.Port), auth, e.From, to, msg)
}

func sslEmail(host string, port int, a smtp.Auth, from string, to []string, msg []byte) error {
	addr := host + ":" + strconv.Itoa(port)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	if err = c.Auth(a); err != nil {
		return err
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addrTo := range to {
		if err = c.Rcpt(addrTo); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()
	return err

}
