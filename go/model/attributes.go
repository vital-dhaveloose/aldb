package model

import "github.com/vital-dhaveloose/aldb/model/ref"

const (
	AttrSetIdBlob = "blob-attrs"
)

type AttributeSet struct {
	Manifest *Manifest
	// Attributes string --> ( nil | string | float64 | bool | map[string]interface{} | []interface{} )
	Attributes map[string]interface{}
}

type Manifest struct {
	ref.ManifestRef
}
