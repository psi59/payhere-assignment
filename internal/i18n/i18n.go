package i18n

import (
	"embed"
	"io/fs"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

//go:generate go run gen.go
//go:embed compiled/active.*.toml
var embedFS embed.FS

var bundle = i18n.NewBundle(language.English)

func init() {
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	if err := fs.WalkDir(embedFS, ".", func(path string, entry fs.DirEntry, err error) error {
		if entry.IsDir() {
			return nil
		}

		filename := entry.Name()
		if _, err := bundle.LoadMessageFileFS(embedFS, filepath.Join("compiled", filename)); err != nil {
			return errors.Wrap(err, "failed to load message")
		}

		return nil
	}); err != nil {
		log.Fatal().Err(err).Msg("failed to register language files")
	}
}

func T(lang language.Tag, id string, placeholder map[string]any) (msg string) {
	loc := i18n.NewLocalizer(bundle, lang.String())
	defer func() {
		if r := recover(); r != nil {
			msg = "undefined message"
		}
	}()

	msg = loc.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: placeholder,
	})

	return
}
