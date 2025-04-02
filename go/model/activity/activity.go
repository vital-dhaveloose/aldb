package activity

import (
	"github.com/vital-dhaveloose/aldb/model"
)

type Activity struct {
	model.Node
	//Participations is the list of participations of entities in the activity.
	Participations []Participation
	Subs           []*model.Node
	// //RecordSchema is the schema that each record activity in this activity (relation is-record-in) has
	// //to satisfy with one of its attribute sets (referenced in the link, see Link#RecordAttributeSetRef).
	// RecordSchema *Schema
}

// TODO common attribute sets:
// - time (during which a node is considered current, to place it on a timeline)
// - blob manifest (should be attribute set)
// - participation?
// - link to schema?
