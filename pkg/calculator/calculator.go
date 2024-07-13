package calculator

type Calculator struct {
	restrictionByParticipant RestrictionsByParticipant
}

func NewCalculator(participants map[string]string, restrictions map[string][]string) (*Calculator, error) {
	rs := NewRestrictionsByParticipant(participants, restrictions)
	if err := rs.validate(); err != nil {
		return nil, err
	}
	return &Calculator{restrictionByParticipant: rs}, nil

}

func (c *Calculator) CalculateRecipient() SenderByRecipient {
	res := make(SenderByRecipient, len(c.restrictionByParticipant))

	allParticipants := make(map[ParticipantID]struct{})
	remainsParticipants := make(map[ParticipantID]struct{})
	for p := range c.restrictionByParticipant {
		allParticipants[p] = struct{}{}
		remainsParticipants[p] = struct{}{}
	}

	for senderID := range allParticipants {
		for recipientID := range remainsParticipants {

			if senderID == recipientID {
				continue
			}

			if restrictRecipientIDs := c.restrictionByParticipant[senderID]; restrictRecipientIDs.contains(recipientID) {
				continue
			}

			res[senderID] = recipientID
			delete(remainsParticipants, recipientID)
			break
		}
	}

	if len(res) != len(c.restrictionByParticipant) {
		res = c.CalculateRecipient()
	}

	return res
}
