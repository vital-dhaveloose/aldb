package examples

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/vital-dhaveloose/aldb/activity"
	"github.com/vital-dhaveloose/aldb/activity/attributes"
	"github.com/vital-dhaveloose/aldb/activity/blob"
	"github.com/vital-dhaveloose/aldb/activity/participation"
	"github.com/vital-dhaveloose/aldb/common/datetime"
	"github.com/vital-dhaveloose/aldb/common/lang"
	"github.com/vital-dhaveloose/aldb/common/mediatype"
	"github.com/vital-dhaveloose/aldb/ref"
)

func TestCompile(_ *testing.T) {
	CreateExampleData()
}

func CreateExampleData() activity.Activity {
	leadRole := participation.ParticipationRole{
		ParticipationRoleRef: participation.ParticipationRoleRef{ParticipationRoleId: "lead"},
	}
	authorRole := participation.ParticipationRole{
		ParticipationRoleRef: participation.ParticipationRoleRef{ParticipationRoleId: "author"},
	}

	projoProjectManifest := attributes.Manifest{
		ManifestRef: ref.ManifestRef{Id: urlMustParse("http://projo.com/schemas/project")},
	}

	vital := participation.Person{
		Ref: participation.EntityRef{Host: "viwi.eu", EntityId: "vital.dhaveloose"},
	}

	msxProjectRef := ref.ActivityRef{Id: urlMustParse("aldb.clientcorp.eu/activities/project-x")}
	msxProjectStart, _ := time.Parse(time.RFC3339, "2019-07-25")
	msxProject := activity.Activity{
		ActivityRef: msxProjectRef,
		Label:       lang.LocalizableString{lang.LangAny: "Project X"},
		Participations: []participation.Participation{
			{
				ParticipationRef: ref.ParticipationRef{ActivityRef: msxProjectRef, ParticipationId: "1"},
				Entity:           &vital,
				Role:             &leadRole,
				Period:           datetime.Period{Start: msxProjectStart},
			},
		},
		Period: datetime.Period{
			Start: msxProjectStart,
		},
		AttributeSets: map[string]attributes.AttributeSet{
			"projo-attrs": {
				Manifest: &projoProjectManifest,
				Attributes: map[string]interface{}{
					"totalBudget":   map[string]interface{}{"currency": "EUR", "amount": float64(456000.00)},
					"priorityClass": "normal",
				},
			},
		},
	}

	rndRef := ref.ActivityRef{Id: urlMustParse("aldb.clientcorp.eu/activities/rnd")}
	rndStart, _ := time.Parse(time.RFC3339, "2020-07-25")
	rndProject := activity.Activity{
		ActivityRef: rndRef,
		Label:       lang.LocalizableString{lang.LangAny: "R&D"},
		Participations: []participation.Participation{
			{
				ParticipationRef: ref.ParticipationRef{ActivityRef: rndRef, ParticipationId: "1"},
				Entity:           &vital,
				Role:             &leadRole,
				Period:           datetime.Period{Start: rndStart},
			},
		},
		Period: datetime.Period{
			Start: rndStart,
		},
		Supers: []*activity.Activity{
			&msxProject,
		},
		AttributeSets: map[string]attributes.AttributeSet{
			"projo-attrs": {
				Manifest: &projoProjectManifest,
				Attributes: map[string]interface{}{
					"totalBudget":   map[string]interface{}{"currency": "EUR", "amount": float64(123000.00)},
					"priorityClass": "normal",
				},
			},
		},
	}

	someDocumentRef := ref.ActivityRef{Id: urlMustParse("aldb.clientcorp.eu/activities/doc-3")}
	someDocument := activity.Activity{
		ActivityRef: someDocumentRef,
		Label:       lang.LocalizableString{lang.LangAny: "some document"},
		Participations: []participation.Participation{
			{
				ParticipationRef: ref.ParticipationRef{ActivityRef: someDocumentRef, ParticipationId: "1"},
				Entity:           &vital,
				Role:             &authorRole,
			},
		},
		AttributeSets: map[string]attributes.AttributeSet{
			"text-attrs": {
				Manifest: &attributes.Manifest{ManifestRef: ref.ManifestRef{Id: urlMustParse("aldb.org/attribute-manifests/text")}},
				Attributes: map[string]interface{}{
					"language": "en-gb",
				},
			},
		},
		Supers: []*activity.Activity{
			&rndProject,
		},
		Blob: &blob.Blob{
			Manifest: &blob.BlobManifest{
				MediaType: mediatype.MediaTypeMustParse("text/plain; charset=UTF-8"),
				Size:      17,
			},
			Bytes: []byte("This is contents!"),
		},
	}

	fmt.Println(someDocument)
	return someDocument
}

func urlMustParse(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}
