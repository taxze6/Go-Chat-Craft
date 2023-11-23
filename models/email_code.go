package models

type MailOptions struct {
	MailHost string
	MailPort uint
	MailUser string
	MailPass string
	MailTo   string
	Subject  string
	Body     string
}
