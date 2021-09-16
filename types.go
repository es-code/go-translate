package trans

type translationFiles map[string]map[string]map[string]interface{}

type appLocalConfig struct {
	defaultLang string
	supportedLanguages []string
	translationsFiles translationFiles
}
