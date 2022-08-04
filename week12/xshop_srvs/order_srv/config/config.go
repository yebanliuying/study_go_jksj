package config

type SrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type RedisConfig struct {
	Host   string `mapstructure:"host" json:"host"`
	Port   int    `mapstructure:"port" json:"port"`
	Expire int    `mapstructure:"expire" json:"expire"`
}

type JaegerConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type ServerConfig struct {
	Name       string       `mapstructure:"name" json:"name"`
	Host       string       `mapstructure:"host" json:"host"`
	Tags       []string     `mapstructure:"tags" json:"tags"`
	RedisInfo  RedisConfig  `mapstructure:"redis" json:"redis"`
	MysqlInfo  MysqlConfig  `mapstructure:"mysql" json:"mysql"`
	ConsulInfo ConsulConfig `mapstructure:"consul" json:"consul"`
	JaegerInfo JaegerConfig `mapstructure:"jaeger" json:"jaeger"`

	//商品微服务配置
	GoodsSrvInfo SrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	//库存微服务配置
	InventorySrvInfo SrvConfig `mapstructure:"inventory_srv" json:"inventory_srv"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
