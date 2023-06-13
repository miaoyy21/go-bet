package xmd

type UserBase struct {
	isDebug bool

	gold    int
	website string
	origin  string
	url     string
	cookie  string
	agent   string
	unix    string
	code    string
	device  string
	id      string
	token   string
}

func NewUserBase(debug bool, gold int, website string, origin string, url string, cookie string, agent string, unix string, code string, device string, id string, token string) UserBase {
	return UserBase{
		isDebug: debug,

		gold:    gold,
		website: website,
		origin:  origin,
		url:     url,
		cookie:  cookie,
		agent:   agent,
		unix:    unix,
		code:    code,
		device:  device,
		id:      id,
		token:   token,
	}
}
