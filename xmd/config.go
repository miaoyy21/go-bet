package xmd

type Config struct {
	IsDebug    bool    `json:"is_debug"`
	Secs       float64 `json:"secs"`
	Wx         float64 `json:"wx"`
	Exp        float64 `json:"exp"`
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
