package aldb

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type MultiTypeUnmarshaler[U any] struct {
	GetCandidates func() []U
	SetResult     func(u U)
}

func (u *MultiTypeUnmarshaler[U]) UnmarshalJSON(bts []byte) error {
	cands := u.GetCandidates()
	for i := range cands {
		cand := cands[i]
		if um, castOk := any(cand).(json.Unmarshaler); castOk {
			err := um.UnmarshalJSON(bts)
			if err != nil {
				continue
			}
			u.SetResult(cand)
			return nil
		}
	}
	return fmt.Errorf("no candidate unmarshalled bytes without error")
}

func (u *MultiTypeUnmarshaler[U]) UnmarshalYAML(n *yaml.Node) error {
	cands := u.GetCandidates()
	for i := range cands {
		cand := cands[i]
		if um, castOk := any(cand).(yaml.Unmarshaler); castOk {
			err := um.UnmarshalYAML(n)
			if err != nil {
				continue
			}
			u.SetResult(cand)
			return nil
		}
	}
	return fmt.Errorf("no candidate unmarshalled bytes without error")
}
