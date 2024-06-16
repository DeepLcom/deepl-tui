package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/cluttrdev/deepl-go/deepl"

	"github.com/DeepLcom/deepl-tui/internal/handlers"
	"github.com/DeepLcom/deepl-tui/internal/ui"
)

// Application is the main type composing the ui and translator client.
type Application struct {
	ui         *ui.UI
	translator *deepl.Translator

	textChanged chan struct{}

	sourceLangs []string
	sourceLang  string
	targetLangs []string
	targetLang  string

	formality string

	glossaries handlers.GlossariesHandler
	glossaryID string
}

// NewApplication creates and returns a new apllication.
func NewApplication(t *deepl.Translator) *Application {
	tui := ui.NewUI()
	tui.EnableMouse(true)
	tui.EnablePaste(false)

	return &Application{
		ui:         tui,
		translator: t,
	}
}

// Run initializes the application and runs the main loop.
func (app *Application) Run() error {
	app.textChanged = make(chan struct{})
	defer close(app.textChanged)

	if err := app.setLanguageOptions(); err != nil {
		app.ui.SetFooter(err.Error())
	}

	if err := app.setFormalityOptions(); err != nil {
		app.ui.SetFooter(err.Error())
	}

	app.setupGlossaryHandling()

	app.ui.SetInputTextChangedFunc(func() {
		app.textChanged <- struct{}{}
	})

	go func() {
		period := 500 * time.Millisecond
		ticker := time.NewTicker(period)

		var changed bool
		for {
			select {
			case _, ok := <-app.textChanged:
				if !ok {
					break
				}
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

func (app *Application) setError(err error) {
	app.ui.SetFooter(fmt.Sprintf("Error: %v", err))
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
			app.sourceLang = app.sourceLangs[index]
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
			app.targetLang = app.targetLangs[index]
			app.updateTranslation()
		},
	)

	return nil
}

func (app *Application) setFormalityOptions() error {
	app.ui.SetFormalityOptions(
		[]string{"Automatic", "Formal tone", "Informal tone"},
		func(text string, index int) {
			switch text {
			case "Automatic":
				app.formality = ""
			case "Formal tone":
				app.formality = "prefer_more"
			case "Informal tone":
				app.formality = "prefer_less"
			}
			app.updateTranslation()
		},
	)

	return nil
}

func (app *Application) setupGlossaryHandling() {
	if err := app.updateGlossaries(); err != nil {
		app.ui.SetFooter(err.Error())
	}

	app.ui.SetGlossaryDataFunc(func(id string) (deepl.GlossaryInfo, []deepl.GlossaryEntry) {
		info, ok := app.glossaries.Get(id)
		if !ok {
			return info, nil
		}
		entries, ok := app.glossaries.Entries(id)
		if !ok {
			var err error
			entries, err = app.glossaries.FetchEntries(app.translator, id)
			if err != nil {
				app.ui.SetFooter(err.Error())
			}
		}
		return info, entries
	})

	app.ui.SetGlossarySelectedFunc(func(id string) {
		app.glossaryID = id
		app.updateTranslation()
	})

	app.ui.SetGlossaryCreateFunc(func(name string, source string, target string, entries [][2]string) {
		if err := app.glossaries.Create(app.translator, name, source, target, entries); err != nil {
			app.ui.SetFooter(err.Error())
			return
		}

		if err := app.updateGlossaries(); err != nil {
			app.ui.SetFooter(err.Error())
		}
	})

	app.ui.SetGlossaryUpdateFunc(func(id string, name string, entries [][2]string) {
		info, ok := app.glossaries.Get(id)
		if !ok {
			app.ui.SetFooter(fmt.Sprintf("Unknown glossary id: %s", id))
			return
		}

		if err := app.glossaries.Create(app.translator, name, info.SourceLang, info.TargetLang, entries); err != nil {
			app.ui.SetFooter(err.Error())
			return
		}

		if err := app.glossaries.Delete(app.translator, id); err != nil {
			app.ui.SetFooter(err.Error())
		}

		if err := app.updateGlossaries(); err != nil {
			app.ui.SetFooter(err.Error())
		}
	})

	app.ui.SetGlossaryDeleteFunc(func(id string) {
		if err := app.glossaries.Delete(app.translator, id); err != nil {
			app.ui.SetFooter(err.Error())
		}

		if err := app.updateGlossaries(); err != nil {
			app.ui.SetFooter(err.Error())
		}
	})
}

func (app *Application) updateTranslation() {
	go func() {
		app.ui.QueueUpdateDraw(func() {
			app.ui.ClearOutputText()

			text := app.ui.GetInputText()
			if text == "" {
				return
			} else if app.targetLang == "" {
				app.setError(fmt.Errorf("Target language not set"))
				return
			}

			var opts []deepl.TranslateOption
			if app.sourceLang != "" {
				opts = append(opts, deepl.WithSourceLang(app.sourceLang))
			}
			if app.formality != "" {
				opts = append(opts, deepl.WithFormality(app.formality))
			}
			if app.glossaryID != "" {
				opts = append(opts, deepl.WithGlossaryID(app.glossaryID))
			}

			translations, err := app.translator.TranslateText([]string{text}, app.targetLang, opts...)
			if err != nil {
				app.setError(err)
				return
			}

			for _, translation := range translations {
				if err := app.ui.WriteOutputText(strings.NewReader(translation.Text)); err != nil {
					app.setError(err)
					return
				}
			}
		})
	}()
}

func (app *Application) updateGlossaries() (err error) {
	if err := app.glossaries.FetchGlossaries(app.translator); err != nil {
		return err
	}

	var opts [][2]string
	for _, info := range app.glossaries.List() {
		opts = append(opts, [2]string{info.GlossaryId, info.Name})
	}

	// sort by name
	sort.SliceStable(opts, func(i, j int) bool {
		return opts[i][1] < opts[j][1]
	})

	app.ui.SetGlossaryOptions(opts)

	// >>>>>>>>
	if err := app.glossaries.FetchLanguages(app.translator); err != nil {
		return err
	}
	langs := app.glossaries.GetSourceLangs("")

	app.ui.SetGlossaryLanguageOptions(langs)
	// <<<<<<<<

	return nil
}
