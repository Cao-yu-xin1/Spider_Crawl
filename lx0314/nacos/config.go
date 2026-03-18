package nacos

type AppConfig struct {
	Mysql  Mysql
	Redis  Redis
	AliPay AliPay
	Consul Consul
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

type Consul struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	ServiceName string `yaml:"serviceName"`
	ServicePort int    `yaml:"servicePort"`
	TTL         int    `yaml:"ttl"`
}

type AliPay struct {
	PrivateKey string `yaml:"privateKey"`
	AppId      string `yaml:"appId"`
	NotifyURL  string `yaml:"notifyURL"`
	ReturnURL  string `yaml:"returnURL"`
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
