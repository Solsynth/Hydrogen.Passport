package i18n

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var Bundle *i18n.Bundle

func InitInternationalization() {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("json", jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal)
	Bundle.LoadMessageFileFS(FS, "locale.en.json")
	Bundle.LoadMessageFileFS(FS, "locale.zh.json")
}
