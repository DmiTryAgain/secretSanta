package calculator

type Calculator struct {
	restrictionByParticipant Restrictions
}

func NewCalculator(restrictions map[string][]string) *Calculator {
	return &Calculator{restrictionByParticipant: NewRestrictions(restrictions)}
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

			if restrictRecipientIDs := c.restrictionByParticipant[senderID]; restrictRecipientIDs.Contains(recipientID) {
				continue
			}

			res[senderID] = recipientID
			delete(remainsParticipants, recipientID)
			break
		}
	}

	//todo: Ахтунг! Провалидировать входные данные, чтобы не зарекурситься
	if len(res) != len(c.restrictionByParticipant) {
		res = c.CalculateRecipient()
	}

	return res
}
