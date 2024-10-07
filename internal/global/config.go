package global

type Config struct {
	JWT struct {
		Secret string
		Expire int64 // hour
		Issuer string
	}
}

var Conf *Config