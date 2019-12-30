// TODO: rename to transform
package mutation

import "net/http"

type RequestMutation interface {
	Change(r *http.Request)
}

type ResponseMutation interface {
	Change(r *http.Response)
}

type MutationConfig struct {
	Headers       map[string]string `yaml:"request-headers"`
	RemoveHeaders []string          `yaml:"remove-headers"`
	Path          map[string]string
}

type HeaderMutation struct {
	Headers map[string]string
}

func (h *HeaderMutation) Change(r *http.Request) {
	if len(h.Headers) == 0 {
		return
	}

	for k, v := range h.Headers {
		r.Header.Set(k, v)
	}

}
