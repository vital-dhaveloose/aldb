package activity

import (
	"github.com/vital-dhaveloose/aldb/activity/attributes"
	"github.com/vital-dhaveloose/aldb/activity/blob"
	"github.com/vital-dhaveloose/aldb/activity/participation"
	"github.com/vital-dhaveloose/aldb/common/datetime"
	"github.com/vital-dhaveloose/aldb/common/lang"
	"github.com/vital-dhaveloose/aldb/ref"
)

type Activity struct {
	ref.ActivityRef
	//Label is a localizable name for the activity.
	Label lang.Localizable
	//Period during which the activiy is considered "current".
	Period datetime.Period
	//Participations is the list of participations of entities in the activity.
	Participations []participation.Participation
	Subs           []*Activity
	Supers         []*Activity
	//AttributeSets contain the structured content of the activity.
	AttributeSets map[string]attributes.AttributeSet

	// //RecordSchema is the schema that each record activity in this activity (relation is-record-in) has
	// //to satisfy with one of its attribute sets (referenced in the link, see Link#RecordAttributeSetRef).
	// RecordSchema *Schema

	//Blob contains the unstructured content of the activity.
	Blob *blob.Blob
}
