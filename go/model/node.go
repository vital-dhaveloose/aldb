package model

import (
	"github.com/vital-dhaveloose/aldb/model/common/lang"
	"github.com/vital-dhaveloose/aldb/model/ref"
)

type Node struct {
	ref.NodeRef
	//Label is a localizable name for the node. Optional, default nil.
	Label lang.Localizable
	//Supers contains references to all activities this node is part of.
	Supers []ref.NodeRef
	//AttributeSets contain the structured content of the node.
	AttributeSets map[string]AttributeSet
	//Blob contains the unstructured content of the node.
	Blob *Blob
}
