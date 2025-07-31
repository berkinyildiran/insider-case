package cache

type Config struct {
	Host string `mapstructure:"host" validate:"required,hostname|ip"`
	Port int    `mapstructure:"port" validate:"required,min=1,max=65535"`
}
