package guard

import "net/http"

type Rule interface {
	Check(r *http.Request) bool
}

type AvariableDomainsRule struct {
	domains []string
}

func (a *AvariableDomainsRule) Check(r *http.Request) bool {
	for _, n := range a.domains {
		if r.Host == n {
			return true
		}
	}
	return false
}
