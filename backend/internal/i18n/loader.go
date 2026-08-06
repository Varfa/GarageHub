package i18n

import (
	"encoding/json"
	"errors"
	"os"
)

func loadLanguageFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	translations := make(map[string]string)
	if err := json.Unmarshal(data, &translations); err != nil {
		return nil, err
	}
	return translations, nil
}
func Load(language string) (map[string]string, error) {

	if language == "" {
		return nil, errors.New("language is empty")
	}
	path := "internal/i18n/" + language + ".json"

	return loadLanguageFile(path)
}
