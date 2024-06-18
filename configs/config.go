package configs

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

var cfg *conf

type conf struct {
	// mapstructure é um alias para resolver o problema de letra maiuscula e minuscula, deixando assim equivalente um ao outro no .env
	DBDriver      string `mapstructure:DB_DRIVER`
	DBHost        string `mapstructure:DB_HOST`
	DBPort        string `mapstructure:DB_PORT`
	DBUser        string `mapstructure:DB_USER`
	DBPassword    string `mapstructure:DB_PASSWORD`
	DBName        string `mapstructure:DB_NAME`
	WebserverPort string `mapstructure:WEB_SERVER_PORT`
	JWTSecret     string `mapstructure:JWT_SECRET`
	JwtExpiresIn  int    `mapstructure:JWT_EXPIRESIN`
	TokenAuth     *jwtauth.JWTAuth
}

func LoadConfig(path string) (*conf, error) {
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	// Adcionando essa linha as configurações sempre irão pegar do .env, não vai subscrever outras variáveis
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	cfg.TokenAuth = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)
	return cfg, err
}

// função config se configurada ela vai ser inicializada antes da função main
// func init() {

// }
