package config

type AppConfig struct {
	Mysql  Mysql
	Redis  Redis
	AliPay AliPay
}

type Mysql struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	Database int    `yaml:"database"`
}

type AliPay struct {
	PrivateKey string
	AppId      string
	NotifyURL  string
	ReturnURL  string
}

type Nacos struct {
	Host        string
	Port        int
	NamespaceId string
	DataId      string
	Group       string
	Username    string
	Password    string
}
