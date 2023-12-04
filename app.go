package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cluttrdev/deepl-go/deepl"

	"github.com/DeepLcom/deepl-tui/internal/ui"
)

type Application struct {
	ui         *ui.UI
	translator *deepl.Translator

	input chan string
	text  string

	sourceLangs []string
	sourceLang  *string
	targetLangs []string
	targetLang  *string
}

func NewApplication(t *deepl.Translator) (*Application, error) {
	tui := ui.NewUI()
	tui.EnableMouse(true)

	return &Application{
		ui:         tui,
		translator: t,

		input: make(chan string),
	}, nil
}

func (app *Application) Run() error {
	defer close(app.input)

	if err := app.setLanguageOptions(); err != nil {
		app.ui.SetFooter(err.Error())
	}

	app.ui.SetInputTextChangedFunc(func() {
		app.input <- app.ui.GetInputText()
	})

	go func() {
		period := 500 * time.Millisecond
		ticker := time.NewTicker(period)

		var changed bool
		for {
			select {
			case text, ok := <-app.input:
				if !ok {
					break
				}
				app.text = text
				changed = true
				ticker.Reset(period)
			case <-ticker.C:
				if changed {
					app.updateTranslation()
					changed = false
				}
			}
		}
	}()

	if err := app.ui.Run(); err != nil {
		return err
	}

	return nil
}

func (app *Application) setLanguageOptions() error {
	sourceLangs, err := app.translator.GetLanguages("source")
	if err != nil {
		return fmt.Errorf("error getting source languages: %w", err)
	}

	app.sourceLangs = make([]string, 1, len(sourceLangs)+1)
	var sourceLangOpts = make([]string, 1, len(sourceLangs)+1)
	sourceLangOpts[0] = "Detect language"
	for _, lang := range sourceLangs {
		app.sourceLangs = append(app.sourceLangs, lang.Code)
		sourceLangOpts = append(sourceLangOpts, lang.Name)
	}

	app.ui.SetSourceLangOptions(
		sourceLangOpts,
		func(text string, index int) {
			app.sourceLang = &app.sourceLangs[index]
			app.updateTranslation()
		},
	)

	targetLangs, err := app.translator.GetLanguages("target")
	if err != nil {
		return fmt.Errorf("error getting target languages: %w", err)
	}

	app.targetLangs = make([]string, 0, len(targetLangs))
	var targetLangOpts = make([]string, 0, len(targetLangs))
	for _, lang := range targetLangs {
		app.targetLangs = append(app.targetLangs, lang.Code)
		targetLangOpts = append(targetLangOpts, lang.Name)
	}

	app.ui.SetTargetLangOptions(
		targetLangOpts,
		func(text string, index int) {
			app.targetLang = &app.targetLangs[index]
			app.updateTranslation()
		},
	)

	return nil
}

func (app *Application) updateTranslation() (err error) {
	if app.text == "" {
		return nil
	} else if app.targetLang == nil {
		return errors.New("Target language not set")
	}

	text := []string{app.text}
	targetLang := *app.targetLang

	var opts []deepl.TranslateOption
	if app.sourceLang != nil && *app.sourceLang != "" {
		opts = append(opts, deepl.WithSourceLang(*app.sourceLang))
	}

	translations, err := app.translator.TranslateText(text, targetLang, opts...)
	if err != nil {
		return err
	}

	app.ui.ClearOutputText()
	for _, translation := range translations {
		if err := app.ui.WriteOutputText(strings.NewReader(translation.Text)); err != nil {
			return err
		}
	}

	return nil
}
