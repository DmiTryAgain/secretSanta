package calculator

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	minParticipantsCount = 3
)

var (
	ErrNotEnoughParticipants = fmt.Errorf("слишком мало участников. Участников должно быть больше %d", minParticipantsCount)
	ErrIncorrectRestrictions = errors.New("некорректные ограничения")
)

type RestrictionsByParticipant map[ParticipantID]ParticipantIDs

func NewRestrictionsByParticipant(in map[string][]string) (out RestrictionsByParticipant) {
	if in == nil {
		return
	}

	out = make(RestrictionsByParticipant, len(in))
	for k, v := range in {
		out[ParticipantID(k)] = NewParticipantIDs(v)
	}

	return
}

// participants Возвращает всех участников
func (r RestrictionsByParticipant) participants() ParticipantIDs {
	res := make(ParticipantIDs, 0, len(r))
	for p := range r {
		res = append(res, p)
	}

	return res
}

func (r RestrictionsByParticipant) validate() error {
	// Выходим, если участников мало
	if len(r) < minParticipantsCount {
		return ErrNotEnoughParticipants
	}

	possibleRecipientsByParticipant := make(map[ParticipantID]ParticipantIDs)
	participants := r.participants()

	for p, rs := range r {
		rsm := rs.toMap()
		var remainsParticipants ParticipantIDs
		// Пробегаемся по участникам и заполняем только доступных к получению участников
		for _, participant := range participants {
			if _, ok := rsm[participant]; ok || participant == p {
				continue
			}

			remainsParticipants = append(remainsParticipants, participant)
		}

		// Если хоть у кого-то уже нет доступных получателей, выходим с ошибкой
		if len(remainsParticipants) == 0 {
			return errors.Wrap(ErrIncorrectRestrictions, fmt.Sprintf("у участника %s нет доступных получателей", p))
		}

		// Запоминаем возможных получателей для каждого участника
		possibleRecipientsByParticipant[p] = remainsParticipants
	}

	// Смотрим, чтобы возможных вариантов было больше, чем учасников с такими же распределением вариантов
	for p, pr := range possibleRecipientsByParticipant {
		countByParticipant := map[ParticipantID]struct{}{p: {}}
		for pToCompare, prToCompare := range possibleRecipientsByParticipant {
			// Скипаем одного и того же участника
			if p == pToCompare {
				continue
			}

			if pr.isEqualValues(prToCompare) {
				countByParticipant[pToCompare] = struct{}{}
			}
		}

		// Если участников с одинаковыми получателями больше, чем самих получателей, распределить не получится
		if len(countByParticipant) > len(pr) {
			errText := "У участников:\n"
			for id := range countByParticipant {
				errText += fmt.Sprintf("%s\n", id)
			}
			errText += "количество получателей меньше, чем их самих. Распределить не получится"
			return errors.Wrap(ErrIncorrectRestrictions, errText)
		}
	}

	return nil
}

type ParticipantID string
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

func (pp ParticipantIDs) toMap() map[ParticipantID]struct{} {
	res := make(map[ParticipantID]struct{}, len(pp))
	for _, p := range pp {
		res[p] = struct{}{}
	}

	return res
}

// nolint:unused
// subtractParticipants Убирает из слайса участников, которые переданы в аргументе
func (pp ParticipantIDs) subtractParticipants(pIDs ParticipantIDs) ParticipantIDs {
	// Возвращаем всех участников, если пришла пустота
	if len(pIDs) == 0 {
		return pp
	}

	// Для удобства помещаем делаем из слайса мапу
	pm := pIDs.toMap()

	// Заполняем результат
	res := make(ParticipantIDs, 0)
	for _, p := range pp {
		// Если такой есть, не добавляем его
		if _, ok := pm[p]; ok {
			continue
		}

		res = append(res, p)
	}

	return res
}

func (pp ParticipantIDs) contains(id ParticipantID) bool {
	for _, pID := range pp {
		if pID == id {
			return true
		}
	}

	return false
}

// isEqualValues Проверяет, одинаковые ли по значениям 2 слайса с участниками
func (pp ParticipantIDs) isEqualValues(ps ParticipantIDs) bool {
	if len(pp) != len(ps) {
		return false
	}

	psm := ps.toMap()
	for _, pID := range pp {
		if _, ok := psm[pID]; !ok {
			return false
		}
	}

	return true
}

type SenderByRecipient map[ParticipantID]ParticipantID
