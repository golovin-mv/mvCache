package guard

import (
	"net/http"
)

type ForbiddenPolicy struct {
}

func (p *ForbiddenPolicy) Reject(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(""))
}
