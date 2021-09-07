// SPDX-License-Identifier: MIT

package console

import (
	"net/http"

	"golang.org/x/text/message"
)

type Response struct {
	http.ResponseWriter
	status int
}

func (r *Response) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *Response) WriteVisitLog(p *message.Printer, url string, succ, erro *Logger) {
	msg := p.Sprintf("visit url", r.status, url)
	if r.status > 399 {
		erro.Println(msg)
	} else {
		succ.Println(msg)
	}
}
