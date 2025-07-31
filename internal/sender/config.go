package sender

type Config struct {
	Address    string `mapstructure:"address" validate:"required,url"`
	CacheTTL   int    `mapstructure:"cache_ttl" validate:"required,min=60,max=3600"`
	FetchLimit int    `mapstructure:"fetch_limit" validate:"required,min=1,max=1000"`
}
