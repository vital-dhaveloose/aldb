package blob

import (
	"github.com/vital-dhaveloose/aldb/common/mediatype"
)

type BlobManifest struct {
	MediaType mediatype.MediaType
	Size      int
}

type Blob struct {
	Manifest *BlobManifest
	Bytes    []byte
}
