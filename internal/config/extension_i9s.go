package config

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	viewToExtension = make(map[string][]Extension)
	i9sExtLock      = sync.RWMutex{}
)

func GetViewToExtension() map[string][]Extension {
	i9sExtLock.RLock()
	defer i9sExtLock.RUnlock()

	return viewToExtension
}

var i9sExtDir = filepath.Join(K9sHome(), "i9s", "ext")

type Extensions struct {
	Extension []Extension `yaml:"extensions"`
}

type Extension struct {
	View    string   `yaml:"view"`
	Name    string   `yaml:"name"`
	Args    []string `yaml:"args"`
	Command string   `yaml:"command"`
}

func (p Extension) String() string {
	return fmt.Sprintf("[%s-%s] %s(%s)", p.View, p.Name, p.Command, strings.Join(p.Args, " "))
}

func InitI9sExtension() {
	loadI9sExtensions()
}

func loadI9sExtensions() {
	pp := &Extensions{}

	if err := filepath.WalkDir(i9sExtDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Error().Err(err).Msgf("load i9s ext: fail to walk %s, skip...", path)
			return nil
		}
		if d.IsDir() {
			return nil
		}

		bs, err := os.ReadFile(path)
		if err != nil {
			log.Error().Err(err).Msgf("load i9s ext: fail to read %s, skip...", path)
			return nil
		}

		dec := yaml.NewDecoder(bytes.NewReader(bs))
		var ext Extension

		for {
			err := dec.Decode(&ext)
			if err != nil {
				if err != io.EOF {
					log.Error().Err(err).Msgf("load i9s ext: fail to decode from %s, skip...", path)
				}
				break
			}
			pp.Extension = append(pp.Extension, ext)
			ext = Extension{}
		}

		return nil
	}); err != nil {
		log.Error().Err(err).Msgf("load i9s ext walk dir %s met err", i9sExtDir)
	}

	newViewToExtension := make(map[string][]Extension)
	for _, ext := range pp.Extension {
		newViewToExtension[ext.View] = append(newViewToExtension[ext.View], ext)
	}
	i9sExtLock.Lock()
	viewToExtension = newViewToExtension
	i9sExtLock.Unlock()

	log.Info().Msgf("load i9s ext %s get: %+v", i9sExtDir, newViewToExtension)
}
