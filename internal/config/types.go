package config

type Postgres struct {
	User     string `yaml:"user" env:"POSTGRES_USER" long:"pg-user" description:"Postgres user"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" long:"pg-pass" description:"Postgres password"`
	Host     string `yaml:"host" env:"POSTGRES_HOST" long:"pg-host" description:"Postgres host"`
	Port     string `yaml:"port" env:"POSTGRES_PORT" long:"pg-port" description:"Postgres port"`
	DB       string `yaml:"db" env:"POSTGRES_DB" long:"pg-db" description:"Postgres database name"`
	URI      string `yaml:"uri" env:"POSTGRES_URI" long:"pg-uri" description:"Postgres URI"`
	MaxConns string `yaml:"max_conns" env:"MAX_CONNS" long:"pg-max-conns" description:"Postgres max poll connection"`
}

type Auth struct {
	MaxBytesKey int `yaml:"max_bytes_key" env:"MAX_BYTES_KEY" long:"max-bytes-key" description:"Max bytes of API key"`
}

type Config struct {
	Postgres Postgres `yaml:"postgres"`
	Config   string   `env:"CONFIG" short:"c" long:"config" description:"Path to yaml config"`
	Addr     string   `yaml:"address" env:"ADDRESS" short:"a" long:"addresss" description:"Address for HTTP server"`
	Auth     Auth     `yaml:"auth"`
	IsDev    bool     `yaml:"is_dev" env:"IS_DEV" short:"d" long:"dev" description:"Dev mode"`
}
