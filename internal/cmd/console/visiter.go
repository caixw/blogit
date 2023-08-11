// SPDX-License-Identifier: MIT

package console

import (
	"net/http"

	"github.com/issue9/localeutil"
	"golang.org/x/text/message"
)

type visiterResponse struct {
	http.ResponseWriter
	status int
}

func Visiter(next http.Handler, p *message.Printer, succ, erro *Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := &visiterResponse{ResponseWriter: w}
		next.ServeHTTP(ww, r)
		ww.printLog(p, r.URL.String(), succ, erro)
	})
}

func (r *visiterResponse) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *visiterResponse) printLog(p *message.Printer, url string, succ, erro *Logger) {
	msg := localeutil.Phrase("visit url %d %s", r.status, url).LocaleString(p)
	if r.status > 399 {
		erro.Println(msg)
	} else {
		succ.Println(msg)
	}
}
