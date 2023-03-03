package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

func ReadEnv() error {
	f, err := ioutil.ReadFile(".env")
	if err != nil {
		return errors.New("Failed to open .env, check the file's name or it's presence")
	}

	var newLineSplit []string = strings.Split(string(f), "\n")

	for i := 0; i < len(newLineSplit); i++ {
		var env = strings.Split(newLineSplit[i], "=")
		if len(env) == 2 {
			os.Setenv(env[0], env[1])
		}
	}

	return nil
}
