package scheduler

type Config struct {
	Interval int `mapstructure:"interval" validate:"required,min=5,max=300"`
}
