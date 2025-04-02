package examples

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/vital-dhaveloose/aldb/model"
	"github.com/vital-dhaveloose/aldb/model/activity"
	"github.com/vital-dhaveloose/aldb/model/common/datetime"
	"github.com/vital-dhaveloose/aldb/model/common/lang"
	"github.com/vital-dhaveloose/aldb/model/common/mediatype"
	"github.com/vital-dhaveloose/aldb/model/ref"
)

func TestCompile(_ *testing.T) {
	CreateExampleData()
}

func CreateExampleData() activity.Activity {
	leadRole := activity.ParticipationRole{
		ParticipationRoleRef: activity.ParticipationRoleRef{ParticipationRoleId: "lead"},
	}
	authorRole := activity.ParticipationRole{
		ParticipationRoleRef: activity.ParticipationRoleRef{ParticipationRoleId: "author"},
	}

	projoProjectManifest := model.Manifest{
		ManifestRef: ref.ManifestRef{Id: urlMustParse("http://projo.com/schemas/project")},
	}

	vital := activity.Person{
		Ref: activity.EntityRef{Host: "viwi.eu", EntityId: "vital.dhaveloose"},
	}

	someDocumentRef := ref.NodeRef{Id: urlMustParse("aldb.clientcorp.eu/activities/doc-3")}
	rndRef := ref.NodeRef{Id: urlMustParse("aldb.clientcorp.eu/activities/rnd")}
	msxProjectRef := ref.NodeRef{Id: urlMustParse("aldb.clientcorp.eu/activities/project-x")}

	someDocument := activity.Activity{
		Node: model.Node{
			NodeRef: someDocumentRef,
			Label:   lang.LocalizableString{lang.LangAny: "some document"},
			AttributeSets: map[string]model.AttributeSet{
				"text-attrs": {
					Manifest: &model.Manifest{ManifestRef: ref.ManifestRef{Id: urlMustParse("aldb.org/attribute-manifests/text")}},
					Attributes: map[string]interface{}{
						"language": "en-gb",
					},
				},
			},
			Supers: []ref.NodeRef{
				rndRef,
			},
			Blob: &model.Blob{
				Manifest: &model.BlobManifest{
					MediaType: mediatype.MediaTypeMustParse("text/plain; charset=UTF-8"),
					Size:      17,
				},
				Bytes: []byte("This is contents!"),
			},
		},
		Participations: []activity.Participation{
			{
				ParticipationRef: ref.ParticipationRef{NodeRef: someDocumentRef, ParticipationId: "1"},
				Entity:           &vital,
				Role:             &authorRole,
			},
		},
	}

	rndStart, _ := time.Parse(time.RFC3339, "2020-07-25")
	rndProject := activity.Activity{
		Node: model.Node{
			NodeRef: rndRef,
			Label:   lang.LocalizableString{lang.LangAny: "R&D"},
			Supers: []ref.NodeRef{
				msxProjectRef,
			},
			AttributeSets: map[string]model.AttributeSet{
				"projo-attrs": {
					Manifest: &projoProjectManifest,
					Attributes: map[string]interface{}{
						"totalBudget":   map[string]interface{}{"currency": "EUR", "amount": float64(123000.00)},
						"priorityClass": "normal",
					},
				},
			},
		},
		Participations: []activity.Participation{
			{
				ParticipationRef: ref.ParticipationRef{NodeRef: rndRef, ParticipationId: "1"},
				Entity:           &vital,
				Role:             &leadRole,
				Period:           datetime.Period{Start: rndStart},
			},
		},
		Subs: []*model.Node{
			&someDocument.Node,
		},
	}

	msxProjectStart, _ := time.Parse(time.RFC3339, "2019-07-25")
	msxProject := activity.Activity{
		Node: model.Node{
			NodeRef: msxProjectRef,
			Label:   lang.LocalizableString{lang.LangAny: "Project X"},
			AttributeSets: map[string]model.AttributeSet{
				"projo-attrs": {
					Manifest: &projoProjectManifest,
					Attributes: map[string]interface{}{
						"totalBudget":   map[string]interface{}{"currency": "EUR", "amount": float64(456000.00)},
						"priorityClass": "normal",
					},
				},
			},
		},
		Participations: []activity.Participation{
			{
				ParticipationRef: ref.ParticipationRef{NodeRef: msxProjectRef, ParticipationId: "1"},
				Entity:           &vital,
				Role:             &leadRole,
				Period:           datetime.Period{Start: msxProjectStart},
			},
		},
		Subs: []*model.Node{
			&rndProject.Node,
		},
	}

	fmt.Println(msxProject)
	return someDocument
}

func urlMustParse(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}
