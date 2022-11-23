// Copyright 2020 Authors of Cilium
// Copyright 2022 Ilya Dmitrichenko
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/errordeveloper/cue-utils/template"
)

type Config struct {
	BaseDirectory string

	templates map[string]*template.Generator
}

func (c *Config) Load() error {
	packagePaths := map[string]struct{}{}

	err := filepath.WalkDir(c.BaseDirectory, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && filepath.Ext(path) == ".cue" {
			packagePaths[filepath.Dir(path)] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("unable to list avaliable config templates in %q: %w", c.BaseDirectory, err)
	}

	c.templates = map[string]*template.Generator{}

	for packagePath := range packagePaths {
		template := template.NewGenerator(packagePath)
		if err := template.CompileAndValidate(); err != nil {
			return fmt.Errorf("unable to load config template from %q: %w", packagePaths, err)
		}
		c.templates[template.ImportPath] = template
	}

	if len(c.templates) == 0 {
		return fmt.Errorf("no config templates found in %q", c.BaseDirectory)
	}
	return nil
}

func (c *Config) HaveExistingTemplate(name string) bool {
	_, ok := c.templates[name]
	return ok
}

func (c *Config) Get(name string) (*template.Generator, error) {
	if !c.HaveExistingTemplate(name) {
		return nil, fmt.Errorf("unknown template %q", name)
	}
	return c.templates[name], nil
}

func (c *Config) WithResource(name string, obj interface{}) (*template.Generator, error) {
	template, err := c.Get(name)
	if err != nil {
		return nil, err
	}
	return template.WithResource(obj)
}

func (c *Config) WithDefaults(name string, obj interface{}) (*template.Generator, error) {
	template, err := c.Get(name)
	if err != nil {
		return nil, err
	}
	return template.WithDefaults(obj)
}

func (c *Config) ApplyDefaults(name string, obj interface{}) error {
	template, err := c.Get(name)
	if err != nil {
		return err
	}
	templateWithDefaults, err := template.WithDefaults(obj)
	if err != nil {
		return err
	}
	c.templates[name] = templateWithDefaults
	return nil
}

func (c *Config) ExistingTemplates() []string {
	templates := []string{}
	for template := range c.templates {
		templates = append(templates, template)
	}
	return templates
}
