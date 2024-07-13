package participation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/vital-dhaveloose/aldb/base/aldberr"
	"github.com/vital-dhaveloose/aldb/common/lang"
)

type EntityRef struct {
	Host     string
	EntityId string
}

func (r *EntityRef) ToName() string {
	if !r.IsComplete() {
		return ""
	}
	return fmt.Sprintf("entities/%s", r.EntityId)
}

func (r *EntityRef) FromName(name string) error {
	if len(name) == 0 {
		return nil
	}
	if r == nil {
		*r = EntityRef{}
	}
	if !strings.HasPrefix(name, "entity/") {
		return errors.New("cannot convert name to EntityRef: name doesn't start with entity/")
	}
	r.EntityId = strings.TrimPrefix(name, "persons/")
	return nil
}

func (r *EntityRef) IsComplete() bool {
	return r != nil && len(r.EntityId) > 0
}

//Entity is a person, a group of people, an organisation or a computer system. It can be authenticated
// and or referenced from the model (notably in Participation).
type Entity interface {
	EntityRef() EntityRef
}

//region Person

type Person struct {
	Ref  EntityRef
	Name PersonName
}

func (p *Person) EntityRef() EntityRef {
	if p == nil {
		return EntityRef{}
	}
	return p.Ref
}

type PersonName struct {
	Given, Family  string
	OtherGivens    []string
	Prefix, Suffix string
}

//endregion

//region Organisation TODO may be duplicate modelling of Activity

type Organisation struct {
	Ref  EntityRef
	Name LocalizableOrganisationName
}

func (p *Organisation) EntityRef() EntityRef {
	if p == nil {
		return EntityRef{}
	}
	return p.Ref
}

type LocalizableOrganisationName map[lang.Lang]OrganisationName

type OrganisationName struct {
	Abbreviation, Short, Long string
}

const (
	//LocizeParamKeyField allows specifying the field to give localize. Allowed values are "short" (default),
	//"abbreviation" and "long".
	LocizeParamKeyField = "pattern"
)

func (lon LocalizableOrganisationName) Localize(language lang.Lang, params map[string]interface{}) (string, error) {
	errDet := map[string]interface{}{"lang": string(language)}
	//TODO support for strict param
	//TODO support for pattern param
	on, found := lon[language]
	if !found {
		return "", aldberr.New(lang.ErrorCodeLanguageNotFound, "cannot localize organisation name: language not found", errDet)
	}
	return on.Short, nil
}

//endregion
