package server

type Config struct {
	Port int `mapstructure:"port" validate:"required,min=1,max=65535"`
}
