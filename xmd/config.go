package xmd

type Config struct {
	IsDebug    bool    `json:"is_debug"`
	IsExtra    bool    `json:"is_extra"`
	Secs       float64 `json:"secs"`
	Fn         string  `json:"fn"`
	Wx         float64 `json:"wx"`
	Rx         float64 `json:"rx"`
	DataSource string  `json:"datasource"`
	Gold       int     `json:"gold"`
	Origin     string  `json:"origin"`
	URL        string  `json:"url"`
	Cookie     string  `json:"cookie"`
	Agent      string  `json:"agent"`
	UserId     string  `json:"user_id"`
	Token      string  `json:"token"`
	Unix       string  `json:"unix"`
	KeyCode    string  `json:"key_code"`
	DeviceId   string  `json:"device_id"`
}
