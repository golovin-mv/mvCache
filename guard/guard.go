package guard

import "net/http"

type GuardCongif struct {
	Enable           bool
	Policy           string
	AvariableDomains []string `yaml:"avariable-domains"`
}

type Guard struct {
	c *GuardCongif
	p Policy
	r []Rule
}

type Policy interface {
	Reject(w http.ResponseWriter, r *http.Request)
}

func (g *Guard) Guard(w http.ResponseWriter, r *http.Request) bool {
	if !g.checkByRules(r) {
		g.p.Reject(w, r)
		return false
	}

	return true
}

func NewGuard(c *GuardCongif) *Guard {
	g := new(Guard)
	g.c = c

	switch c.Policy {

	case "forbidden":
		g.p = &ForbiddenPolicy{}
	case "break":
		g.p = &BreakPolicy{}
	default:
		g.p = &ForbiddenPolicy{}
	}

	g.addRules()
	return g
}

func (g *Guard) addRules() {
	if len(g.c.AvariableDomains) > 0 {
		g.r = append(g.r, &AvariableDomainsRule{g.c.AvariableDomains})
	}
}

func (g *Guard) checkByRules(req *http.Request) bool {
	if g.r == nil || len(g.r) == 0 {
		return true
	}
	for _, rule := range g.r {
		if !rule.Check(req) {
			return false
		}
	}

	return true
}
