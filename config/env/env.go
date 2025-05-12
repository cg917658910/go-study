package env

var Envs *Env

type MYSQL struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     string `mapstructure:"poet" json:"port" yaml:"port"`
	User     string `mapstructure:"user" json:"user" yaml:"user"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DBName   string `mapstructure:"db_name" json:"db_name" yaml:"db_name"`
}

type Env struct {
	MYSQL MYSQL `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
}
