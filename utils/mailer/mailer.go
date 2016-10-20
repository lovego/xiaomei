package mailer

import (
	"bytes"
	"encoding/base64"
	"errors"
	// "fmt"
	"github.com/bughou-go/xm"
	"mime"
	"net/smtp"
	"path/filepath"
	"strings"
	"time"
)

type Mailer struct {
	smtp.Auth
	Addr, Sender, From string
}

type mailMessage struct {
	headers       string
	receivers, to []string
	title, body   string
}

// sender: name<address@example.com>
func New(host, port, sender, password string) *Mailer {
	if sender == `` {
		panic(`empty sender`)
	}

	mailer := Mailer{Addr: host + `:` + port}
	mailer.Sender, mailer.From = parseNameAddr(sender)
	mailer.Auth = smtp.PlainAuth(``, mailer.From, password, host)

	return &mailer
}

// receivers: [name1<address1@example.com>, name2<address2@example.com>]
func (m *Mailer) Send(receivers []string, title, body string) (err error) {
	if m == nil {
		return nil
	}
	return m.SendWithAttaches(receivers, title, body, nil)
}

// receivers: [ name1<address1@example.com>, ... ],
// attaches: { `filename.ext`: `content`, ... }
func (m *Mailer) SendWithAttaches(
	receivers []string, title, body string, attaches map[string]string,
) (err error) {
	if m == nil {
		return nil
	}
	if len(receivers) == 0 {
		return errors.New(`empty receivers`)
	}
	encRcvs, to := parseReceivers(receivers)
	message := makeMessage(m.Sender, encRcvs, title, body, attaches)

	xm.Protect(func() {
		err = smtp.SendMail(m.Addr, m.Auth, m.From, to, message)
	})
	return
}

const boundary = `f46d043c813270fc6b04c2d223da`

func makeMessage(sender, receivers, title, body string, attaches map[string]string) []byte {
	var buf bytes.Buffer
	multipart := len(attaches) > 0
	buf.WriteString(makeHeaders(sender, receivers, title, multipart))
	if multipart {
		buf.WriteString(makeMultipartBody(body, attaches))
	} else {
		buf.WriteString("\r\n" + body + "\r\n")
	}
	return buf.Bytes()
}

func makeHeaders(sender, receivers, title string, multipart bool) string {
	headers := "From: " + sender + "\r\n" +
		"To: " + receivers + "\r\n" +
		"Date: " + time.Now().Format(time.RFC1123Z) + "\r\n" +
		"Subject: " + mime.BEncoding.Encode(`utf-8`, title) + "\r\n" +
		"MIME-Version: 1.0\r\n"

	if multipart {
		headers += "Content-Type: multipart/mixed; boundary=\"" + boundary + "\"\r\n"
	} else {
		headers += "Content-Type: text/plain; charset=UTF-8\r\n"
	}
	return headers
}

func makeMultipartBody(body string, attaches map[string]string) string {
	result := "\r\n--" + boundary + "\r\n" +
		"Content-Type: text/plain; charset=utf-8\r\n" +
		"\r\n" + body + "\r\n"

	for name, content := range attaches {
		result += makeAttach(name, content)
	}
	result += "\r\n--" + boundary + "--\r\n"
	return result
}

func makeAttach(name, content string) string {
	return "\r\n--" + boundary + "\r\n" +
		"Content-Type: " + mime.TypeByExtension(filepath.Ext(name)) + "\r\n" +
		"Content-Disposition: attachment; filename=\"" + name + "\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64Encode(content) + "\r\n"
}

func base64Encode(content string) string {
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(content)))
	base64.StdEncoding.Encode(buf, []byte(content))

	var result bytes.Buffer
	// devide base64 content in lines of 78 chars
	for i, l := 0, len(buf); i < l; i++ {
		result.WriteByte(buf[i])
		if (i+1)%78 == 0 {
			result.WriteString("\r\n")
		}
	}
	return result.String()
}

func parseReceivers(receivers []string) (encoded string, addrs []string) {
	var nameAddrs []string
	for _, receiver := range receivers {
		nameAddr, addr := parseNameAddr(receiver)
		nameAddrs = append(nameAddrs, nameAddr)
		addrs = append(addrs, addr)
	}
	encoded = strings.Join(nameAddrs, `,`)
	return
}

func parseNameAddr(nameAddr string) (encoded, addr string) {
	arr := strings.Split(nameAddr, `<`)
	encoded = mime.BEncoding.Encode(`utf-8`, arr[0]) + `<` + arr[1]
	addr = strings.TrimRight(arr[1], `>`)
	return
}
