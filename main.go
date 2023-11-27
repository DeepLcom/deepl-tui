package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/cluttrdev/deepl-go/deepl"
)

func main() {
	if err := execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func execute() error {
	auth_key := parseAuthKey()
	if auth_key == "" {
		return errors.New("Missing required authentication key.")
	}

	translator, err := deepl.NewTranslator(auth_key)
	if err != nil {
		return err
	}

	app, err := NewApplication(translator)
	if err != nil {
		return err
	}

	return app.Run()
}

func parseAuthKey() string {
	// parse args
	f := flag.String("auth-key", "", "the authentication key as given in your DeepL account.")
	flag.Parse()
	if *f != "" {
		return *f
	}

	// parse env
	return os.Getenv("DEEPL_AUTH_KEY")
}
