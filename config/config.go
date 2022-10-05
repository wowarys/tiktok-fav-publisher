package config

const (
	CacheTypeSet   = "set"
	CacheTypeRedis = "redis"
)

type Config struct {
	DBAddr string `env:"DB_ADDR" envDefault:"localhost:6379"`

	TelegramToken string `env:"TG_TOKEN,required"`
	ChannelID     int64  `env:"CHANNEL_ID,required"`

	TikTokUsername string `env:"TIKTOK_USERNAME,required"`
	CacheType      string `env:"CACHE_TYPE" envDefault:"set"`
}
