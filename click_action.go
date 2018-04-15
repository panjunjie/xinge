package xinge

const (
	TYPE_ACTIVITY = 1
	TYPE_URL      = 2
	TYPE_INTENT   = 3
)

type ClickAction struct {
	ActionType                  int      `json:"action_type,omitempty"`
	Browser                     *Browser `json:"browser,omitempty"`
	Activity                    string   `json:"activity,omitempty"`
	Intent                      string   `json:"intent,omitempty"`
	AtyAttr                     *AtyAttr `json:"aty_attr,omitempty"`
	PackageName                 string   `json:"package_name,omitempty"`
	PackageDownloadUrl          string   `json:"-"`
	ConfirmOnPackageDownloadUrl int      `json:"-"`
}

func NewClickAction() *ClickAction {
	return &ClickAction{
		ActionType:                  TYPE_ACTIVITY,
		Browser:                     NewBrowser(),
		Activity:                    "",
		AtyAttr:                     NewAtyAttr(),
		PackageName:                 "",
		PackageDownloadUrl:          "",
		ConfirmOnPackageDownloadUrl: 1,
	}
}

func NewSimplekAction(packageName, activity string) *ClickAction {
	action := NewClickAction()
	action.Activity = activity
	action.PackageName = ""
	return action
}

func (s *ClickAction) SetPackageName(packageName string) {
	s.PackageName = packageName
}

func (s *ClickAction) SetActivity(activity string) {
	s.Activity = activity
}

func (s *ClickAction) SetIntent(intent string) {
	s.Intent = intent
}

func (s *ClickAction) SetActionType(actionType int) {
	s.ActionType = actionType
}

func (s *ClickAction) SetPackageDownloadUrl(packageDownloadUrl string) {
	s.PackageDownloadUrl = packageDownloadUrl
}

func (s *ClickAction) SetConfirmOnPackageDownloadUrl(confirmOnPackageDownloadUrl int) {
	s.ConfirmOnPackageDownloadUrl = confirmOnPackageDownloadUrl
}

func (s *ClickAction) SetBrowser(browser *Browser) {
	s.Browser = browser
}

func (s *ClickAction) SetAtyAttr(atyAttr *AtyAttr) {
	s.AtyAttr = atyAttr
}

func (s *ClickAction) IsValid() bool {
	if s.ActionType < TYPE_ACTIVITY || s.ActionType > TYPE_INTENT {
		return false
	}

	if s.ActionType == TYPE_URL {
		if s.Browser.Url == "" || s.Browser.ConfirmOnUrl < 0 || s.Browser.ConfirmOnUrl > 1 {
			return false
		}
		return true
	}

	if s.ActionType == TYPE_INTENT {
		if s.Intent == "" {
			return false
		}
		return true
	}

	return true
}

type Browser struct {
	Url          string `json:"url,omitempty"`
	ConfirmOnUrl int    `json:"confirm,omitempty"`
}

func NewBrowser() *Browser {
	return &Browser{
		Url:          "",
		ConfirmOnUrl: 0,
	}
}

func (s *Browser) SetUrl(url string) {
	s.Url = url
}

func (s *Browser) SetConfirmOnUrl(confirmOnUrl int) {
	s.ConfirmOnUrl = confirmOnUrl
}

type AtyAttr struct {
	AtyAttrIntentFlag        int `json:"if,omitempty"`
	AtyAttrPendingIntentFlag int `json:"pf,omitempty"`
}

func NewAtyAttr() *AtyAttr {
	return &AtyAttr{
		AtyAttrIntentFlag:        0,
		AtyAttrPendingIntentFlag: 0,
	}
}

func (s *AtyAttr) SetAtyAttrIntentFlag(atyAttrIntentFlag int) {
	s.AtyAttrIntentFlag = atyAttrIntentFlag
}

func (s *AtyAttr) SetConfirmOnUrl(atyAttrPendingIntentFlag int) {
	s.AtyAttrPendingIntentFlag = atyAttrPendingIntentFlag
}
