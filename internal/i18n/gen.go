//go:build ignore
// +build ignore

package main

import (
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog/log"
)

const msgIDFilename = "msg_generated.go"

func main() {
	if err := os.Remove(msgIDFilename); err != nil && !os.IsNotExist(err) {
		log.Fatal().Msgf("failed to remove generated file: %+v", err)
	}

	files, err := ioutil.ReadDir("compiled")
	if err != nil {
		log.Fatal().Msgf("failed to read locale files: %v", err)
	}

	var c msgMap
	for _, file := range files {
		filename := file.Name()
		if file.IsDir() || !strings.HasSuffix(filename, "toml") {
			continue
		}

		if !strings.Contains(filename, "active.en.toml") {
			continue
		}

		f, err := os.Open("compiled/" + filename)
		if err != nil {
			log.Fatal().Msgf("failed to open file: %v", err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatal().Msgf("failed to close file: %v", err)
			}
		}()

		c = make(msgMap)
		decoder := toml.NewDecoder(f)
		if _, err := decoder.Decode(&c); err != nil {
			log.Fatal().Msgf("failed to decode toml file: %v", err)
		}
	}

	tmpl, err := template.New("msg").Parse(`package i18n
	
	// Code generated by go generate; DO NOT EDIT.
	
	const (
	{{- range $id, $msg := . }}
		{{ $id }} = "{{ $id }}"
	{{- end }}
	)
	`)
	if err != nil {
		log.Fatal().Msgf("Failed to parse i18n template")
	}

	f, err := os.OpenFile(msgIDFilename, os.O_RDWR|os.O_CREATE, 0o755)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open file")
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal().Err(err).Msg("failed to close file")
		}
	}()

	if err := tmpl.Execute(f, c); err != nil {
		log.Fatal().Msgf("Failed to execute repository template: %v", err)
	}
}

type msgMap map[string]any
