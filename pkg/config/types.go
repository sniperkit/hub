package config

// ServiceConfig contains configuration information for a service
type ServiceConfig struct {
	Token string
	User  string
}

// OutputConfig sontains configuration information for an output
type OutputConfig struct {
	SpinnerIndex    int    `yaml:"spinner_index" yaml:"spinner_index" toml:"spinner_index"`
	SpinnerInterval int    `yaml:"spinner_interval" yaml:"spinner_interval" toml:"spinner_interval"`
	SpinnerColor    string `yaml:"spinner_color" yaml:"spinner_color" toml:"spinner_color"`
}

type StorageConfig struct {
	DB  DBConfig    `json:"db" yaml:"db" toml:"db"`
	IDX IndexConfig `json:"index" yaml:"index" toml:"index"`
	S3  DBConfig    `json:"s3" yaml:"s3" toml:"s3"`
}

type IndexConfig struct {
	IndexPath string `json:"index_path" yaml:"index_path" toml:"index_path"`
}

type DBConfig struct {
	Adapter     string `env:"STORAGE_DB_ADAPTER" default:"sqlite"`
	DSN         string `env:"STORAGE_DB_DSN"`
	Name        string `env:"STORAGE_DB_NAME" default:"hub_db"`
	Host        string `env:"STORAGE_DB_HOST" default:"localhost"`
	Port        string `env:"STORAGE_DB_PORT" default:"3306"`
	User        string `env:"STORAGE_DB_USER"`
	Password    string `env:"STORAGE_DB_PASSWORD"`
	AutoMigrate bool   `env:"STORAGE_DB_AUTO_MIGRATE"`
}

type S3Config struct {
	AccessKeyID     string `env:"AWS_ACCESS_KEY_ID" json:"-" yaml:"-" toml:"-"`
	SecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY" json:"-" yaml:"-" toml:"-"`
	Region          string `env:"AWS_REGION" json:"region" yaml:"region" toml:"region"`
	Bucket          string `env:"AWS_BUCKET" json:"bucket" yaml:"bucket" toml:"bucket"`
}

// Config contains configuration information
type Config struct {
	Storage  DBConfig                  `json:"storage" yaml:"storage" toml:"storage"`
	Services map[string]*ServiceConfig `json:"services" yaml:"services" toml:"services"`
	Outputs  map[string]*OutputConfig  `json:"outputs" yaml:"outputs" toml:"outputs"`
}
