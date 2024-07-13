package lang

import (
	"github.com/vital-dhaveloose/aldb/base/aldberr"
	"github.com/vital-dhaveloose/aldb/util"
)

const (
	ErrorCodeLanguageNotFound = "common-lang-not-found"
)

type Lang string

const (
	LangAny Lang = "*"
	LangEn  Lang = "en"
)

//Localizable is something that given a target language and some parameters can return a localized string
//or an error.
type Localizable interface {
	Localize(lang Lang, params map[string]interface{}) (string, error)
}

const (
	//LocalizeParamKeyStrict indicates that only the specific given language can be used, not LangAny
	//or similar languages. The default value is false.
	LocalizeParamKeyStrict = "strict"
)

type LocalizableString map[Lang]string

func (s LocalizableString) Localize(lang Lang, params map[string]interface{}) (string, error) {
	errDet := map[string]interface{}{"lang": string(lang)}
	str, found := s[lang]
	if found {
		return str, nil
	}
	strict := util.GetEntryBool(params, LocalizeParamKeyStrict, false)
	if strict {
		return "", aldberr.New(ErrorCodeLanguageNotFound, "can't localize LocalizableString: language not found (strict)", errDet)
	}
	str, found = s[LangAny]
	if !found {
		return "", aldberr.New(ErrorCodeLanguageNotFound, "can't localize LocalizableString: no supported language found", errDet)
	}
	return str, nil
}
