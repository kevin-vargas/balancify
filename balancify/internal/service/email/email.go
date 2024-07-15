package email

import (
	"bytes"
)

type bytesBuffer struct {
	bytes.Buffer
}

func (b *bytesBuffer) space() {
	b.WriteString("\r\n")
}

type formatOption func(sb *bytesBuffer)

func WithSubject(s string) formatOption {
	return func(sb *bytesBuffer) {
		sb.WriteString("Subject: ")
		sb.WriteString(s)
		sb.space()
	}
}

func WithFrom(f string) formatOption {
	return func(sb *bytesBuffer) {
		sb.WriteString("From: ")
		sb.WriteString(f)
		sb.space()
	}
}
func WithReceiver(r string) formatOption {
	return func(sb *bytesBuffer) {
		sb.WriteString("To: ")
		sb.WriteString(r)
		sb.space()
	}
}

func WithHtmlBody(buf []byte) formatOption {
	return func(sb *bytesBuffer) {
		sb.WriteString("MIME-version: 1.0;")
		sb.space()
		sb.WriteString("Content-Type: text/html; charset=\"UTF-8\";")
		sb.space()
		sb.space()
		sb.Write(buf)
		sb.space()
	}
}

func Format(opts ...formatOption) []byte {
	var sb bytesBuffer
	for _, o := range opts {
		o(&sb)
	}
	return sb.Bytes()
}
