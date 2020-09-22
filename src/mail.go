package postman

import (
	"net/smtp"

	"github.com/pelletier/go-toml"
)

/**
 * Mail context to target receiver
 */
func Mail(t *toml.Tree, to string, context string) error {
	auth := smtpAuth(t)
	user := t.Get("mail.auth.user").(string)
	addr := t.Get("mail.msg.addr").(string)
	err := smtp.SendMail(
		addr, auth, user,
		[]string{to}, genMsg(t, to, context),
	)

	if err != nil {
		return err
	}

	return nil
}

/**
 * # Generate mail content
 * ```
 * From: Postman <darwinia-postman@yandex.com>
 * To: reciver@example.com
 * Subject: Subject: Darwinia Crashed
 *
 * 8444aa5b-76ce-48de-b9d9-5108b5b39a13
 * ```
 */
func genMsg(t *toml.Tree, to string, body string) []byte {
	cc := t.Get("mail.msg.subject").(string)
	from := t.Get("mail.msg.from").(string)

	plain := "" +
		"From: " + from + "\n" +
		"To: " + to + "\n" +
		cc + " \n\n" +
		body

	return []byte(plain)
}

/**
 * Auth stmp server
 */
func smtpAuth(t *toml.Tree) smtp.Auth {
	ident := t.Get("mail.auth.ident").(string)
	user := t.Get("mail.auth.user").(string)
	pass := t.Get("mail.auth.pass").(string)
	host := t.Get("mail.auth.host").(string)

	return smtp.PlainAuth(ident, user, pass, host)
}
