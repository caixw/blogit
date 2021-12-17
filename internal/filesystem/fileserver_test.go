// SPDX-License-Identifier: MIT

package filesystem

import (
	"net/http"
	"os"
	"testing"

	"github.com/issue9/assert/v2"
	"github.com/issue9/assert/v2/rest"
)

func TestFileServer(t *testing.T) {
	a := assert.New(t, false)

	fs := FileServer(os.DirFS("./"), nil)
	a.NotNil(fs)
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	s := rest.NewServer(a, http.DefaultServeMux, nil)

	// vars.IndexFilename
	s.Get("/assets/").Do(nil).Status(http.StatusNotFound)

	// vars.IndexFilename
	s.Get("/assets/testdata").Do(nil).Status(http.StatusOK).
		BodyFunc(func(a *assert.Assertion, body []byte) {
			a.Contains(string(body), "<html>")
		})
	s.Get("/assets/testdata/").Do(nil).Status(http.StatusOK).
		BodyFunc(func(a *assert.Assertion, body []byte) {
			a.Contains(string(body), "<html>")
		})

	s.Get("/assets/filesystem.go").Do(nil).
		Status(http.StatusOK).
		BodyFunc(func(a *assert.Assertion, body []byte) {
			a.Contains(string(body), "package filesystem")
		})

	// 404
	s.Get("/assets/not-exists").Do(nil).Status(http.StatusNotFound)
}
