package ref

import "net/url"

type ActivityRef struct {
	Id      *url.URL
	Version string
}

//func (r ActivityRef) MarshalJSON() ([]byte, error) {
//	if r.Id == nil {
//		return []byte(""), nil
//	}
//	s := r.Id.String()
//	if len(r.Version) > 0 {
//		s = s + "/" + r.Version
//	}
//	return []byte(s), nil
//}
