package config

//SMSConfig 短信配置结构
type SMSConfig struct {
	Product        string   `yaml:"product"`
	AppKey         string   `yaml:"app_key"`
	AppSecret      string   `yaml:"app_secret"`
	SignName       string   `yaml:"sign_name"`
	DailySendLimit int      `yaml:"daily_send_limit"`
	Templates      []string `yaml:"templates"`
}
