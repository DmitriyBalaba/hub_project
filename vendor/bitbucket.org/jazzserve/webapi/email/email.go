package email

type Email interface {
	SetTemplate(templateKey string) Email
	SetTo(name string, address string) Email
	SetDynamicData(key string, value string) Email
	Send() error
}
