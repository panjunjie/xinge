package xinge

type Style struct {
	BuilderId int    `json:"builder_id,omitempty"`
	Ring      int    `json:"ring,omitempty"`
	Vibrate   int    `json:"vibrate,omitempty"`
	Clearable int    `json:"clearable,omitempty"`
	NId       int    `json:"n_id,omitempty"`
	RingRaw   string `json:"ring_raw,omitempty"`
	Lights    int    `json:"lights,omitempty"`
	IconType  int    `json:"icon_type,omitempty"`
	IconRes   string `json:"icon_res,omitempty"`
	StyleId   int    `json:"style_id,omitempty"`
	SmallIcon string `json:"small_icon,omitempty"`
}

func NewStyle(builderId int) *Style {
	return NewStyleFull(builderId, 0, 1, 1, 0, 1, 0, 1)
}

func NewStyleBase(builderId int, ring int, vibrate int, clearable int, nId int) *Style {
	return NewStyleFull(builderId, ring, vibrate, clearable, nId, 1, 0, 1)
}

func NewStyleFull(builderId int, ring int, vibrate int, clearable int, nId int, lights int, iconType int, styleId int) *Style {
	return &Style{BuilderId: builderId,
		Ring:      ring,
		Vibrate:   vibrate,
		Clearable: clearable,
		NId:       nId,
		Lights:    lights,
		IconType:  iconType,
		StyleId:   styleId,
	}
}

func (s *Style) IsValid() bool {
	if s.Ring < 0 || s.Ring > 1 {
		return false
	}

	if s.Vibrate < 0 || s.Vibrate > 1 {
		return false
	}

	if s.Clearable < 0 || s.Clearable > 1 {
		return false
	}

	if s.Lights < 0 || s.Lights > 1 {
		return false
	}

	if s.IconType < 0 || s.IconType > 1 {
		return false
	}

	if s.StyleId < 0 || s.StyleId > 1 {
		return false
	}

	return true
}
