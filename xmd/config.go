package xmd

type Config struct {
	IsDebug    bool    `json:"is_debug"`
	IsExtra    bool    `json:"is_extra"`
	Rx         float64 `json:"rx"`
	DataSource string  `json:"datasource"`
	Gold       int     `json:"gold"`
	Origin     string  `json:"origin"`
	URL        string  `json:"url"`
	Cookie     string  `json:"cookie"`
	UserId     string  `json:"user_id"`
	Token      string  `json:"token"`
	Unix       string  `json:"unix"`
	KeyCode    string  `json:"key_code"`
	DeviceId   string  `json:"device_id"`
}