package env

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type envVariables struct {
	ApiKey    string
	WalletKey string
	SteloApi  string
}

var Variables envVariables

func Initialize() error {
	godotenv.Load()

	apiKey, err := getEnvVar("API_KEY")
	if err != nil {
		return err
	}

	walletKey, err := getEnvVar("WALLET_KEY")
	if err != nil {
		return err
	}

	steloApi, err := getEnvVar("STELO_API")
	if err != nil {
		return err
	}

	Variables = envVariables{
		ApiKey:    apiKey,
		WalletKey: walletKey,
		SteloApi:  steloApi,
	}

	return nil
}

func getEnvVar(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", errors.New("Environment variable \"" + key + "\" not found")
	}

	return value, nil
}
