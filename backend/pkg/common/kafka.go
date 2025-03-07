package common

import (
	"encoding/json"
	"time"
)

type ScheduleMessage struct {
	ClientID string    `json:"client_id"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	//TODO: add all the necessary fields
}

func (s *ScheduleMessage) UnmarshalJSON(data []byte) error {
	var mp map[string]interface{}
	if err := json.Unmarshal(data, &mp); err != nil {
		return err
	}
	value, ok := mp["client_id"]
	if ok {
		clientID, ok2 := value.(string)
		if ok2 {
			s.ClientID = clientID
		}

	}
	value, ok = mp["start"]
	if ok {
		start, ok2 := value.(string)
		if ok2 {
			ts, err := time.Parse(time.DateOnly, start)
			if err != nil {
				return err
			}
			s.Start = ts
		}

	}
	value, ok = mp["end"]
	if ok {
		end, ok2 := value.(string)
		if ok2 {
			ts, err := time.Parse(time.DateOnly, end)
			if err != nil {
				return err
			}
			s.End = ts
		}

	}
	return nil
}

func (s *ScheduleMessage) MarshalJSON() ([]byte, error) {
	type Alias ScheduleMessage
	return json.Marshal(&struct {
		*Alias
		Start string `json:"start"`
		End   string `json:"end"`
	}{
		Alias: (*Alias)(s),
		Start: s.Start.Format(time.DateOnly),
		End:   s.End.Format(time.DateOnly),
	})
}

type NotifyMessage struct {
	ClientID     string  `json:"client_id"`
	CurrentSpend float64 `json:"current_spend"`
	ClientEmail  string  `json:"client_email"`
}
