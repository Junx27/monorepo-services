package config

type Config struct {
	Port string `env:"PORT" envDefault:"8000"`
	Env  string `env:"ENV" envDefault:"dev"`

	DatabaseUsername string `env:"DB_USERNAME"`
	DatabasePassword string `env:"DB_PASSWORD"`
	DatabaseHost     string `env:"DB_HOST"`
	DatabasePort     string `env:"DB_PORT"`
	DatabaseName     string `env:"DB_NAME"`

	MigrationPath string `env:"MIGRATION_PATH" envDefault:"/app/migrations"`

	AMQPHost string `env:"AMQP_HOST"`
}
