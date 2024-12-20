package config

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PprofDebug struct {
	Enable bool   `yaml:"enable"`
	Port   uint16 `yaml:"port"`
}

func (p PprofDebug) Addr() string {
	return fmt.Sprintf("localhost:%d", p.Port)
}

type Logger struct {
	Output      []string `yaml:"output"`
	Level       string   `yaml:"level"`
	ServerLevel string   `yaml:"server_level"`
}

type SSL struct {
	Enable  bool   `yaml:"enable"`
	Path    string `yaml:"path" env-default:"./../cert"`
	Public  string `yaml:"public"`
	Private string `yaml:"private"`
}

func (s SSL) GetPublic() string {
	return fmt.Sprintf("%s/%s", s.Path, s.Public)
}

func (s SSL) GetPrivate() string {
	return fmt.Sprintf("%s/%s", s.Path, s.Private)
}

type JWT struct {
	Secret        string
	ExpireAccess  int `yaml:"expire_access" env-default:"300"`
	ExpireRefresh int `yaml:"expire_refresh" env-default:"2592000"`
}

type Service struct {
	Port           uint16   `yaml:"port" env-default:"8080"`
	AllowedOrigins []string `yaml:"allowed_origins" env-default:"http://localhost:5173"`
}

type DbSQL struct {
	Name              string `yaml:"name"`
	User              string `yaml:"user"`
	Password          string
	Host              string `yaml:"host"`
	Port              uint16 `yaml:"port"`
	MaxConns          uint16 `yaml:"max_conns"`
	MinConns          uint16 `yaml:"min_conns"`
	MaxConnLifetime   uint32 `yaml:"max_conn_life_time"`
	MaxConnIdleTime   uint32 `yaml:"max_conn_idle_time"`
	HealthCheckPeriod uint32 `yaml:"health_check_period"`
	ConnectTimeout    uint32 `yaml:"connect_timeout"`
}

func (db *DbSQL) PgxPoolConfig() *pgxpool.Config {
	dbConfig, err := pgxpool.ParseConfig(db.GetConnStr())
	if err != nil {
		return nil
	}

	dbConfig.MaxConns = int32(db.MaxConns)
	dbConfig.MinConns = int32(db.MinConns)
	dbConfig.MaxConnLifetime = time.Duration(db.MaxConnLifetime) * time.Second
	dbConfig.MaxConnIdleTime = time.Duration(db.MaxConnIdleTime) * time.Second
	dbConfig.HealthCheckPeriod = time.Duration(db.HealthCheckPeriod) * time.Second
	dbConfig.ConnConfig.ConnectTimeout = time.Duration(db.ConnectTimeout) * time.Second

	return dbConfig
}

func (db *DbSQL) GetConnStr() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", db.User, db.Password, db.Host, db.Port, db.Name)
}

type DbRedis struct {
	Name     string `yaml:"name"`
	Password string
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DB       uint16 `yaml:"db"`
}

type Email struct {
	Address  string `yaml:"address"`
	Password string
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

func (e Email) AddrSvc() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}

type Kafka struct {
	Enable bool     `yaml:"enable"`
	Host   string   `yaml:"host"`
	Port   uint16   `yaml:"port"`
	Topics []string `yaml:"topics"`
}

func (k Kafka) Addr() string {
	return fmt.Sprintf("%s:%d", k.Host, k.Port)
}

type Google struct {
	SecretPath  string `yaml:"secret_path"`
	RedirectUrl string `yaml:"redirect_url"`
}

type Aes struct {
	Key string `yaml:"key"`
}
