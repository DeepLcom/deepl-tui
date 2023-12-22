package handlers

import (
	"strings"

	"github.com/cluttrdev/deepl-go/deepl"
)

type glossary struct {
	Info    deepl.GlossaryInfo
	Entries []deepl.GlossaryEntry
}

// GlossariesHandler manages glossaries.
type GlossariesHandler struct {
	languages  []deepl.LanguagePair
	glossaries []deepl.GlossaryInfo
	entries    map[string][]deepl.GlossaryEntry
}

// FetchLanguages retreives the list of supported glossary langues pairs.
func (h *GlossariesHandler) FetchLanguages(client *deepl.Translator) error {
	langs, err := client.GetGlossaryLanguagePairs()
	if err != nil {
		return err
	}

	h.languages = langs
	return nil
}

// FetchGlossaries retreives all available glossaries.
func (h *GlossariesHandler) FetchGlossaries(client *deepl.Translator) error {
	infos, err := client.ListGlossaries()
	if err != nil {
		return err
	}

	h.glossaries = infos
	return nil
}

// FetchEntries retreives the entries of a single glossary.
func (h *GlossariesHandler) FetchEntries(client *deepl.Translator, id string) ([]deepl.GlossaryEntry, error) {
	entries, err := client.GetGlossaryEntries(id)
	if err != nil {
		return nil, err
	}

	if h.entries == nil {
		h.entries = make(map[string][]deepl.GlossaryEntry)
	}
	h.entries[id] = make([]deepl.GlossaryEntry, len(entries))
	copy(h.entries[id], entries)
	return entries, nil
}

// List returns the list of available glossaries.
func (h *GlossariesHandler) List() []deepl.GlossaryInfo {
	infos := make([]deepl.GlossaryInfo, len(h.glossaries))
	copy(infos, h.glossaries)
	return infos
}

// Get returns meta information for a single glossary.
// If there is no information available with the given glossary ID, the second
// return value will be `false`.
//
// Use [GlossaryHandler.FetchGlossaries] to fetch all available glossaries.
func (h *GlossariesHandler) Get(id string) (deepl.GlossaryInfo, bool) {
	for _, g := range h.glossaries {
		if g.GlossaryId == id {
			return g, true
		}
	}
	return deepl.GlossaryInfo{}, false
}

// Entries returns the entries of a single glossary.
// If there is no data available for the given glossary ID, the second
// return value will be `false`.
//
// Use [GlossaryHandler.FetchEntries] to fetch entries for a glossary.
func (h *GlossariesHandler) Entries(id string) ([]deepl.GlossaryEntry, bool) {
	entries, ok := h.entries[id]
	if !ok {
		return nil, false
	}

	e := make([]deepl.GlossaryEntry, len(entries))
	copy(e, entries)
	return e, true
}

// FindName returns the glossary id for the first glossary found that has the
// given `name`.
// If no glossary is found, an empty string is returned.
func (h *GlossariesHandler) FindName(name string) string {
	for _, g := range h.glossaries {
		if g.Name == name {
			return g.GlossaryId
		}
	}
	return ""
}

// GetTargetLangs returns all supported glossary target languages for the
// given source language.
// If `source` is an empty string, all supported target languages are returned.
func (h *GlossariesHandler) GetTargetLangs(source string) []string {
	var targets []string
	if source != "" {
		for _, pair := range h.languages {
			if pair.SourceLang == strings.ToLower(source) {
				targets = append(targets, pair.TargetLang)
			}
		}
	} else {
		m := make(map[string]struct{})
		for _, pair := range h.languages {
			m[pair.TargetLang] = struct{}{}
		}
		for k := range m {
			targets = append(targets, k)
		}
	}
	return targets
}

// GetSourceLangs returns all supported glossary source languages for the
// given target language.
// If `target` is an empty string, all supported source languages are returned.
func (h *GlossariesHandler) GetSourceLangs(target string) []string {
	var sources []string
	if target != "" {
		for _, pair := range h.languages {
			if pair.TargetLang == strings.ToLower(target) {
				sources = append(sources, pair.SourceLang)
			}
		}
	} else {
		m := make(map[string]struct{})
		for _, pair := range h.languages {
			m[pair.SourceLang] = struct{}{}
		}
		for k := range m {
			sources = append(sources, k)
		}
	}
	return sources
}

// Create creates a new glossary.
func (h *GlossariesHandler) Create(client *deepl.Translator, name string, source string, target string, entries [][2]string) error {
	entries_ := make([]deepl.GlossaryEntry, 0, len(entries))
	for _, entry := range entries {
		entries_ = append(entries_, deepl.GlossaryEntry{
			Source: entry[0],
			Target: entry[1],
		})
	}

	info, err := client.CreateGlossary(name, source, target, entries_)
	if err != nil {
		return err
	}

	if info != nil {
		h.glossaries = append(h.glossaries, *info)
	}
	return nil
}

// Delete deletes a single glossary.
func (h *GlossariesHandler) Delete(client *deepl.Translator, id string) error {
	return client.DeleteGlossary(id)
}
