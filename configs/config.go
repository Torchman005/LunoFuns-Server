package configs

import (
	"fmt"
	"os"
	"time"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Name	string `yaml:"name"`
	Version string `yaml:"version"`
	Env     string `yaml:"env"`
	Debug   bool   `yaml:"debug"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type DatabaseConfig struct {
	Driver string `yaml:"driver"`
	Host	 string `yaml:"host"`
	Port	 int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
	ParseTime bool `yaml:"parse_time"`
	Loc      string `yaml:"loc"`
	MaxIdleConns int `yaml:"max_idle_conns"`
	MaxOpenConns int `yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	LogLevel string `yaml:"log_level"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
	ExpireHours int `yaml:"expire_hours"`
	Issuer string `yaml:"issuer"`
	RefreshExpireDays int `yaml:"refresh_expire_days"`
}

type RedisConfig struct {
	Enable bool   `yaml:"enable"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Password string `yaml:"password"`
	DB int       `yaml:"db"`
	PoolSize int  `yaml:"pool_size"`
}

type LogConfig struct {
	Level string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
	FilePath string `yaml:"file_path"`
	MaxSize int `yaml:"max_size"`
	MaxBackups int `yaml:"max_backups"`
	MaxAge int `yaml:"max_age"`
	Compress bool `yaml:"compress"`
}

type CORSConfig struct {
	Enable bool `yaml:"enable"`
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
	AllowedHeaders []string `yaml:"allowed_headers"`
	AllowCredentials bool `yaml:"allow_credentials"`
	MaxAge int `yaml:"max_age"`
}

type RateLimitConfig struct {
	Enable bool `yaml:"enable"`
	RequestsPerSecond int `yaml:"requests_per_second"`
	Burst int `yaml:"burst"`
}

type SecurityConfig struct {
	CORS CORSConfig `yaml:"cors"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
}

type PasswordPolicyConfig struct {
	MinLength int `yaml:"min_length"`
	RequireUppercase bool `yaml:"require_uppercase"`
	RequireLowercase bool `yaml:"require_lowercase"`
	RequireNumber bool `yaml:"require_number"`
	RequireSpecialChars bool `yaml:"require_special_chars"`
	BcryptCost int `yaml:"bcrypt_cost"`
}

type EmailConfig struct {
	Enable bool `yaml:"enable"`
	SMTPHost string `yaml:"smtp_host"`
	SMTPPort int `yaml:"smtp_port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	FromEmail string `yaml:"from_email"`
	FromName string `yaml:"from_name"`
}

type Config struct {
	App AppConfig `yaml:"app"`
	Server ServerConfig `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`	
	JWT JWTConfig `yaml:"jwt"`
	Redis RedisConfig `yaml:"redis"`
	Log LogConfig `yaml:"log"`
	Security SecurityConfig `yaml:"security"`
	PasswordPolicy PasswordPolicyConfig `yaml:"password_policy"`
	Email EmailConfig `yaml:"email"`
}

var GlobalConfig *Config

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	GlobalConfig = &config
	return &config, nil
}

// 配置验证
func (c *Config) Validate() error {
    // 生产环境必须设置JWT密钥
    if c.App.Env == "production" {
        if c.JWT.Secret == "your-secret-key-change-in-production" {
            return fmt.Errorf("JWT secret must be changed in production")
        }
        if c.Database.Password == "" {
            return fmt.Errorf("database password cannot be empty in production")
        }
    }

	// 验证数据库配置
    if c.Database.Driver == "" {
        return fmt.Errorf("database driver is required")
    }
    
    return nil
}

// 获取DSN连接字符串
func (c *Config) GetDSN() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
        c.Database.Username,
        c.Database.Password,
        c.Database.Host,
        c.Database.Port,
        c.Database.Database,
        c.Database.Charset,
        c.Database.ParseTime,
        c.Database.Loc,
    )
}

// 获取服务器地址
func (c *Config) GetServerAddr() string {
    return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// 是否为开发环境
func (c *Config) IsDevelopment() bool {
    return c.App.Env == "development"
}

// 是否为生产环境
func (c *Config) IsProduction() bool {
    return c.App.Env == "production"
}

// 是否为测试环境
func (c *Config) IsStaging() bool {
    return c.App.Env == "staging"
}

// 获取Redis地址
func (c *Config) GetRedisAddr() string {
    return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

// 获取JWT过期时间
func (c *Config) GetJWTExpireDuration() time.Duration {
    return time.Duration(c.JWT.ExpireHours) * time.Hour
}
