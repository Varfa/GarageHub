package i18n

type Manager struct {
	translations    map[string]map[string]string
	defaultLanguage string
}

func (m *Manager) T(key string) string {
	return m.Translate(m.defaultLanguage, key)
}
func NewManager(
	languages []string,
	defaultLanguage string,
) (*Manager, error) {
	translations := make(map[string]map[string]string)

	for _, language := range languages {
		loaded, err := Load(language)
		if err != nil {
			return nil, err
		}

		translations[language] = loaded
	}

	return &Manager{
		translations:    translations,
		defaultLanguage: defaultLanguage,
	}, nil
}
func (m *Manager) Translate(language string, key string) string {
	languageTranslations, ok := m.translations[language]
	if !ok {
		languageTranslations = m.translations[m.defaultLanguage]
	}

	value, ok := languageTranslations[key]
	if ok {
		return value
	}

	defaultTranslations, ok := m.translations[m.defaultLanguage]
	if !ok {
		return key
	}

	value, ok = defaultTranslations[key]
	if !ok {
		return key
	}

	return value
}
func (m *Manager) HasLanguage(language string) bool {
	_, ok := m.translations[language]
	return ok
}
