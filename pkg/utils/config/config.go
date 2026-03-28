// @AI_GENERATED
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// AppConfig is the root configuration struct for the application.
type AppConfig struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Log       LogConfig       `mapstructure:"log"`
	Migration MigrationConfig `mapstructure:"migration"`
	MCP       MCPConfig       `mapstructure:"mcp"`
	Pprof     PprofConfig     `mapstructure:"pprof"`
	Models    ModelsConfig    `mapstructure:"models"`
	Skills    SkillsConfig    `mapstructure:"skills"`
	Memory    MemoryConfig    `mapstructure:"memory"`
}

// MemoryConfig holds memory context engine configuration.
type MemoryConfig struct {
	Enabled             bool    `mapstructure:"enabled"`
	EmbeddingBaseURL    string  `mapstructure:"embedding_base_url"`
	EmbeddingAPIKey     string  `mapstructure:"embedding_api_key"`
	EmbeddingModel      string  `mapstructure:"embedding_model"`
	EmbeddingDim        int     `mapstructure:"embedding_dim"`
	HybridVectorWeight  float32 `mapstructure:"hybrid_vector_weight"`
	HybridKeywordWeight float32 `mapstructure:"hybrid_keyword_weight"`
	SearchLimit         int     `mapstructure:"search_limit"`
}

// SkillsConfig holds skill loader settings.
type SkillsConfig struct {
	// Dir overrides the workspace skills directory (default: ./skills).
	Dir string `mapstructure:"dir"`
	// ExtraDirs are additional skill directories (lowest priority, before bundled).
	ExtraDirs []string `mapstructure:"extra_dirs"`
}

// ModelsConfig holds AI model provider configuration.
type ModelsConfig struct {
	// Default model used when a session does not specify one.
	DefaultProvider string `mapstructure:"default_provider"`
	DefaultModel    string `mapstructure:"default_model"`
	// Providers maps provider name → provider-level settings.
	Providers map[string]ProviderConfig `mapstructure:"providers"`
}

// ProviderConfig holds settings for a single AI provider.
type ProviderConfig struct {
	// APIKeys is a list of API keys for round-robin rotation.
	APIKeys []string `mapstructure:"api_keys"`
	// BaseURL overrides the default API endpoint (useful for OpenAI-compatible providers).
	BaseURL string `mapstructure:"base_url"`
	// Timeout for requests to this provider (default: 60s).
	Timeout time.Duration `mapstructure:"timeout"`
}

// MCPServerConfig holds configuration for a single MCP server.
type MCPServerConfig struct {
	Name    string   `mapstructure:"name"`
	Command string   `mapstructure:"command"`
	Args    []string `mapstructure:"args"`
	Env     []string `mapstructure:"env"`
	// RequireApproval requires human approval before ANY tool from this server executes.
	RequireApproval bool `mapstructure:"require_approval"`
	// DangerousTools lists specific tool names from this server that require approval.
	// If RequireApproval is true, this list is ignored (all tools require approval).
	DangerousTools []string `mapstructure:"dangerous_tools"`
}

// MCPConfig holds MCP tool server configuration.
type MCPConfig struct {
	Enabled bool              `mapstructure:"enabled"`
	Servers []MCPServerConfig `mapstructure:"servers"`
}

// PprofConfig holds pprof profiling settings.
type PprofConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Addr    string `mapstructure:"addr"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig holds PostgreSQL connection settings.
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// RedisConfig holds Redis connection settings.
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// JWTConfig holds JWT authentication settings.
type JWTConfig struct {
	Secret          string        `mapstructure:"secret"`
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
}

// LogConfig holds logging settings.
// Defined here (not reused from logger package) to avoid circular dependencies.
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// MigrationConfig holds database migration settings.
type MigrationConfig struct {
	Enabled     bool          `mapstructure:"enabled"`
	SourceType  string        `mapstructure:"source_type"`
	SourcePath  string        `mapstructure:"source_path"`
	LockTimeout time.Duration `mapstructure:"lock_timeout"`
	TableName   string        `mapstructure:"table_name"`
}

// LoadConfig reads configuration from the YAML file at path,
// applies environment variable overrides (prefix OPENCLAW),
// and returns a populated AppConfig.
func LoadConfig(path string) (*AppConfig, error) {
	v := viper.New()

	// Set defaults for all sections.
	setDefaults(v)

	// Enable environment variable override.
	v.SetEnvPrefix("OPENCLAW")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Load YAML config file.
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file %q: %w", path, err)
	}

	var cfg AppConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// setDefaults registers reasonable default values for all config sections.
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", 30*time.Second)
	v.SetDefault("server.write_timeout", 30*time.Second)

	// Database defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "openclaw")
	v.SetDefault("database.password", "openclaw")
	v.SetDefault("database.dbname", "openclaw")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("database.conn_max_lifetime", time.Hour)

	// Redis defaults
	v.SetDefault("redis.addr", "localhost:6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)

	// JWT defaults
	v.SetDefault("jwt.secret", "change-me-in-production")
	v.SetDefault("jwt.access_token_ttl", 24*time.Hour)
	v.SetDefault("jwt.refresh_token_ttl", 168*time.Hour)

	// Log defaults
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "console")

	// Migration defaults
	v.SetDefault("migration.enabled", true)
	v.SetDefault("migration.source_type", "filesystem")
	v.SetDefault("migration.source_path", "./migrations")
	v.SetDefault("migration.lock_timeout", 30*time.Second)
	v.SetDefault("migration.table_name", "schema_migrations")

	// MCP defaults
	v.SetDefault("mcp.enabled", false)

	// Pprof defaults
	v.SetDefault("pprof.enabled", false)
	v.SetDefault("pprof.addr", "localhost:6060")

	// Skills defaults
	v.SetDefault("skills.dir", "./skills")

	// Memory defaults
	v.SetDefault("memory.enabled", false)
	v.SetDefault("memory.embedding_model", "BAAI/bge-m3")
	v.SetDefault("memory.embedding_dim", 1024)
	v.SetDefault("memory.hybrid_vector_weight", 0.7)
	v.SetDefault("memory.hybrid_keyword_weight", 0.3)
	v.SetDefault("memory.search_limit", 10)

	// Models defaults
	v.SetDefault("models.default_provider", "ollama")
	v.SetDefault("models.default_model", "llama3")
}

// @AI_GENERATED: end
