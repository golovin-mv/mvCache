package mutation

import "net/http"

type RemoveHeadersMutation struct {
	RemoveHeaders []string
}

func (r *RemoveHeadersMutation) Change(res *http.Response) {
	for _, h := range r.RemoveHeaders {
		res.Header.Del(h)
	}
}
