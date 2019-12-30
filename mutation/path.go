package mutation

import (
	"net/http"
)

type Path struct {
	Routes map[string]string
}

func (p *Path) Change(r *http.Request) {
	if v, ok := p.Routes[r.URL.Path]; ok {
		r.URL.Path = v
	}
}
