package configs

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type Config struct {
	DB_MASTER          string `mapstructure:"DB_MASTER"`
	GCS_BUCKET_NAME    string `mapstructure:"GCS_BUCKET_NAME"`
	GOOGLE_PROJECT_ID  string `mapstructure:"GOOGLE_PROJECT_ID"`
	GOOGLE_STORAGE_URL string `mapstructure:"GOOGLE_STORAGE_URL"`
}

var config Config

func LoadConfig() error {
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.BindEnv("SECRET_KEY")
	viper.BindEnv("DB_MASTER")
	viper.BindEnv("GCS_BUCKET_NAME")
	viper.BindEnv("GOOGLE_PROJECT_ID")
	viper.BindEnv("GOOGLE_STORAGE_URL")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return err
	}
	return nil
}

func GetConfig() *Config {
	return &config
}

func SetupDB(dbstr string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dbstr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
