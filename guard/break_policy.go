package guard

import "net/http"

type BreakPolicy struct {
}

func (p *BreakPolicy) Reject(w http.ResponseWriter, r *http.Request) {
	hj, ok := w.(http.Hijacker)

	// The rw can't be hijacked, return early.
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// Hijack the rw.
	conn, _, err := hj.Hijack()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}

	// Close the hijacked raw tcp connection.
	if err := conn.Close(); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
