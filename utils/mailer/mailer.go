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
	Addr, SenderNameAddr, SenderAddr string
}

type Message struct {
	Receivers   []People // [ name1, address1@example.com ], [ name2, address2@example.com ]
	ContentType string
	Title, Body string
	Attaches    map[string]string // { `filename.ext`: `content`, ... }
}

type People struct {
	Name, Addr string
}

func New(host, port string, sender People, password string) *Mailer {
	if host == `` || port == `` || sender.Addr == `` {
		return nil
	}

	mailer := Mailer{Addr: host + `:` + port}
	mailer.SenderNameAddr = encodeNameAddr(sender.Name, sender.Addr)
	mailer.SenderAddr = sender.Addr
	mailer.Auth = smtp.PlainAuth(``, sender.Addr, password, host)

	return &mailer
}

func (m *Mailer) Send(msg *Message) error {
	if m == nil || msg == nil || len(msg.Receivers) == 0 {
		return nil
	}
	if msg.ContentType == `` {
		msg.ContentType = `text/plain; charset=UTF-8`
	}
	rcvrAddrs := make([]string, len(msg.Receivers))
	for i, people := range msg.Receivers {
		rcvrAddrs[i] = people.Addr
	}
	return smtp.SendMail(m.Addr, m.Auth, m.SenderAddr, rcvrAddrs, m.makeMessage(msg))
}

const boundary = `f46d043c813270fc6b04c2d223da`

func (m *Mailer) makeMessage(msg *Message) []byte {
	var buf bytes.Buffer
	multipart := len(msg.Attaches) > 0
	buf.WriteString(m.makeHeaders(msg, multipart))
	buf.WriteString(m.makeBody(msg, multipart))
	return buf.Bytes()
}

func (m *Mailer) makeHeaders(msg *Message, multipart bool) string {
	headers := "From: " + m.SenderNameAddr + "\r\n" +
		"To: " + encodeNameAddrs(msg.Receivers) + "\r\n" +
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

func encodeNameAddrs(people []People) string {
	var slice []string
	for _, p := range people {
		slice = append(slice, encodeNameAddr(p.Name, p.Addr))
	}
	return strings.Join(slice, `,`)
}

func encodeNameAddr(name, addr string) string {
	if name == `` {
		if i := strings.IndexByte(addr, '@'); i >= 0 {
			name = addr[0:i]
		}
	}
	return mime.BEncoding.Encode(`utf-8`, name) + `<` + addr + `>`
}
