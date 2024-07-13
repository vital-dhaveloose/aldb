package participation

import (
	"github.com/vital-dhaveloose/aldb/common/datetime"
	"github.com/vital-dhaveloose/aldb/ref"
)

type Participation struct {
	ref.ParticipationRef
	Entity Entity
	Period datetime.Period
	Role   *ParticipationRole
}

type ParticipationRoleRef struct {
	ParticipationRoleId string
}

type ParticipationRole struct {
	ParticipationRoleRef
}
