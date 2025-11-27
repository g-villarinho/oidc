package config

import "time"

const (
	development = "development"
	staging     = "staging"
	production  = "production"
)

type Config struct {
	Env       string    `mapstructure:"env"`
	Postgres  Postgres  `mapstructure:"postgres"`
	Redis     Redis     `mapstructure:"redis"`
	Cors      Cors      `mapstructure:"cors"`
	Key       Key       `mapstructure:"key"`
	RateLimit RateLimit `mapstructure:"ratelimit"`
	Session   Session   `mapstructure:"session"`
	Server    Server    `mapstructure:"server"`
	URL       URL       `mapstructure:"url"`
}

type Server struct {
	Port            int           `mapstructure:"port"`
	Host            string        `mapstructure:"host"`
	ShutdownTimeout time.Duration `mapstructure:"shutdowntimeout"`
}

type Postgres struct {
	Host        string        `mapstructure:"host"`
	Port        int           `mapstructure:"port"`
	DBName      string        `mapstructure:"dbname"`
	User        string        `mapstructure:"user"`
	Password    string        `mapstructure:"password"`
	SSLMode     string        `mapstructure:"sslmode"`
	TimeZone    string        `mapstructure:"timezone"`
	MaxConn     int32         `mapstructure:"maxconn"`
	MinConn     int32         `mapstructure:"minconn"`
	MaxIdle     int           `mapstructure:"maxidle"`
	MaxLifeTime time.Duration `mapstructure:"maxlifetime"`
}

type Redis struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	DB           int    `mapstructure:"db"`
	Password     string `mapstructure:"password"`
	MinIdleConns int    `mapstructure:"minidleconns"`
	PoolSize     int    `mapstructure:"poolsize"`
}

type Cors struct {
	AllowedOrigins []string `mapstructure:"allowedorigins"`
	AllowedMethods []string `mapstructure:"allowedmethods"`
	AllowedHeaders []string `mapstructure:"allowedheaders"`
}

type RateLimit struct {
	MaxRequests int           `mapstructure:"maxrequests"`
	Window      time.Duration `mapstructure:"window"`
}

type Key struct {
	PrivateKey string `mapstructure:"privatekey"`
	PublicKey  string `mapstructure:"publickey"`
}

type Session struct {
	Secret        string        `mapstructure:"secret"`
	Duration      time.Duration `mapstructure:"duration"`
	CookieOptions CookieOptions `mapstructure:"cookieoptions"`
}

type CookieOptions struct {
	Name     string `mapstructure:"name"`
	Secure   bool   `mapstructure:"secure"`
	HTTPOnly bool   `mapstructure:"httpOnly"`
	SameSite string `mapstructure:"sameSite"`
}

type URL struct {
	APIBaseURL string `mapstructure:"apibaseurl"`
	AppBaseURL string `mapstructure:"appbaseurl"`
}

func (e *Config) IsDevelopment() bool {
	return e.Env == development
}

func (e *Config) IsStaging() bool {
	return e.Env == staging
}

func (e *Config) IsProduction() bool {
	return e.Env == production
}
