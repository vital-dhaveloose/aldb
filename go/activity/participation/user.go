package participation

import (
	"github.com/vital-dhaveloose/aldb/common/datetime"
	"github.com/vital-dhaveloose/aldb/common/lang"
	"github.com/vital-dhaveloose/aldb/ref"
)

//User represents an Entity that can use the system in an authenticated way for a certain context.
type User struct {
	Entity
	Context UserContext
}

type UserContextRef struct {
	EntityRef     ref.Ref
	UserContextId string
}

type UserContext struct {
	UserContextRef
	Organisation Organisation
	ValidPeriod  datetime.Period
	Description  lang.LocalizableString
}
