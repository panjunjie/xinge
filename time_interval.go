package xinge

type TimeInterval struct {
	StartTime *TimePart `json:"start"`
	EndTime   *TimePart `json:"end"`
}

type TimePart struct {
	Hour int `json:"hour"`
	Min  int `json:"min"`
}

func DefaultTimeInterval() *TimeInterval {
	return &TimeInterval{StartTime: &TimePart{0, 0}, EndTime: &TimePart{23, 59}}
}

func (s *TimeInterval) IsValid() bool {
	if s.StartTime.Hour >= 0 && s.StartTime.Hour <= 23 &&
		s.StartTime.Min >= 0 && s.StartTime.Min <= 59 &&
		s.EndTime.Hour >= 0 && s.EndTime.Hour <= 23 &&
		s.EndTime.Min >= 0 && s.EndTime.Min <= 59 {
		return true
	}
	return false
}
