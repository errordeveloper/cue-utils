// Copyright 2020 Authors of Cilium
// Copyright 2022 Ilya Dmitrichenko
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"
	"os"

	"github.com/errordeveloper/cue-utils/template"
)

type Config struct {
	BaseDirectory string

	templates map[string]*template.Generator
}

func (c *Config) Load() error {
	entries, err := os.ReadDir(c.BaseDirectory)
	if err != nil {
		return fmt.Errorf("unable to list avaliable config templates in %q: %w", c.BaseDirectory, err)
	}

	c.templates = map[string]*template.Generator{}

	for _, entry := range entries {
		if entry.IsDir() {
			// both path.Join and filpath.Join break this by striping leading `./`,
			// just like Go, relative package path in must be prefixed with `./`
			// (or `../`)
			fullPath := c.BaseDirectory + "/" + entry.Name()
			template := template.NewGenerator(fullPath)
			if err := template.CompileAndValidate(); err != nil {
				return fmt.Errorf("unable to load config template from %q: %w", fullPath, err)
			}
			c.templates[entry.Name()] = template
		}
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
