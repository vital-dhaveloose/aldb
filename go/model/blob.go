package model

import (
	"github.com/vital-dhaveloose/aldb/model/common/mediatype"
)

type BlobManifest struct {
	MediaType mediatype.MediaType
	Size      int
}

type Blob struct {
	Manifest *BlobManifest
	Bytes    []byte
}
