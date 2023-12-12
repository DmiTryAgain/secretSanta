package calculator

type Restrictions map[ParticipantID]ParticipantIDs

func NewRestrictions(in map[string][]string) (out Restrictions) {
	if in == nil {
		return
	}

	out = make(Restrictions, len(in))
	for k, v := range in {
		out[ParticipantID(k)] = NewParticipantIDs(v)
	}

	return
}

type ParticipantIDs []ParticipantID

func NewParticipantIDs(in []string) (out ParticipantIDs) {
	if in == nil {
		return
	}

	for _, id := range in {
		out = append(out, ParticipantID(id))
	}

	return
}

func (p ParticipantIDs) Contains(id ParticipantID) bool {
	for _, pID := range p {
		if pID == id {
			return true
		}
	}

	return false
}

type SenderByRecipient map[ParticipantID]ParticipantID
type ParticipantID string
