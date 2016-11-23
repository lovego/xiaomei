package mailer

import (
	"bytes"
	"encoding/base64"
	"mime"
	"net/smtp"
	"path/filepath"
	"strings"
	"time"
)

type Mailer struct {
	smtp.Auth
	HostPort, SenderNameAddr, SenderAddr string
}

type Message struct {
	Receivers   []string // [name1<address1@example.com>, name2<address2@example.com>]
	ContentType string
	Title, Body string
	Attaches    map[string]string // { `filename.ext`: `content`, ... }
}

// sender: name<address@example.com>
func New(host, port, sender, password string) *Mailer {
	if sender == `` {
		panic(`empty sender`)
	}

	mailer := Mailer{HostPort: host + `:` + port}
	mailer.SenderNameAddr, mailer.SenderAddr = parseNameAddr(sender)
	mailer.Auth = smtp.PlainAuth(``, mailer.SenderAddr, password, host)

	return &mailer
}

func (m *Mailer) Send(msg *Message) error {
	if m == nil || msg == nil || len(msg.Receivers) == 0 {
		return nil
	}
	if msg.ContentType == `` {
		msg.ContentType = `text/plain; charset=UTF-8`
	}
	rcvrNameAddrs, rcvrAddrs := parseNameAddrSlice(msg.Receivers)
	message := m.makeMessage(rcvrNameAddrs, msg)

	return smtp.SendMail(m.HostPort, m.Auth, m.SenderAddr, rcvrAddrs, message)
}

const boundary = `f46d043c813270fc6b04c2d223da`

func (m *Mailer) makeMessage(rcvrNameAddrs string, msg *Message) []byte {
	var buf bytes.Buffer
	multipart := len(msg.Attaches) > 0
	buf.WriteString(m.makeHeaders(rcvrNameAddrs, msg, multipart))
	buf.WriteString(m.makeBody(msg, multipart))
	return buf.Bytes()
}

func (m *Mailer) makeHeaders(rcvrNameAddrs string, msg *Message, multipart bool) string {
	headers := "From: " + m.SenderNameAddr + "\r\n" +
		"To: " + rcvrNameAddrs + "\r\n" +
		"Date: " + time.Now().Format(time.RFC1123Z) + "\r\n" +
		"Subject: " + mime.BEncoding.Encode(`utf-8`, msg.Title) + "\r\n" +
		"MIME-Version: 1.0\r\n"

	if multipart {
		headers += "Content-Type: multipart/mixed; boundary=\"" + boundary + "\"\r\n"
	} else {
		headers += "Content-Type: " + msg.ContentType + "\r\n"
	}
	return headers
}

func (m *Mailer) makeBody(msg *Message, multipart bool) string {
	if !multipart {
		return "\r\n" + msg.Body + "\r\n"
	}
	return makeMultipartBody(msg)
}

func makeMultipartBody(msg *Message) string {
	result := "\r\n--" + boundary + "\r\n" +
		"Content-Type: " + msg.ContentType + "\r\n" +
		"\r\n" + msg.Body + "\r\n"

	for name, content := range msg.Attaches {
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

func parseNameAddrSlice(strs []string) (nameAddrs string, addrs []string) {
	var slice []string
	for _, s := range strs {
		nameAddr, addr := parseNameAddr(s)
		slice = append(slice, nameAddr)
		addrs = append(addrs, addr)
	}
	nameAddrs = strings.Join(slice, `,`)
	return
}

func parseNameAddr(s string) (nameAddr, addr string) {
	arr := strings.Split(s, `<`)
	nameAddr = mime.BEncoding.Encode(`utf-8`, arr[0]) + `<` + arr[1]
	addr = strings.TrimRight(arr[1], `>`)
	return
}
