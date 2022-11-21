// Copyright 2020 Authors of Cilium
// Copyright 2022 Ilya Dmitrichenko
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"fmt"

	"cuelang.org/go/cue"

	"github.com/errordeveloper/cue-utils/compiler"
)

const (
	templateKey = "template"
	defaultsKey = "defaults"
	resourceKey = "resource"
)

type Generator struct {
	dir      string
	args     []string
	cue      *compiler.Compiler
	template cue.Value
}

func NewGenerator(dir string, args ...string) *Generator {
	if len(args) == 0 {
		args = []string{"."}
	}
	return &Generator{
		args: args,
		dir:  dir,
		cue:  compiler.NewCompiler(),
	}
}

func (g *Generator) CompileAndValidate() error {
	template, err := g.cue.BuildAll(g.dir, g.args...)
	if err != nil {
		return err
	}

	g.template = template

	return nil
}

func (g *Generator) with(key string, obj interface{}) (*Generator, error) {
	keyPath := cue.ParsePath(key)
	if err := keyPath.Err(); err != nil {
		return nil, err
	}
	result := g.template.FillPath(keyPath, obj)
	if err := result.Err(); err != nil {
		return nil, fmt.Errorf("error after filling %q: %w", key, err)
	}
	return &Generator{
		dir:      g.dir,
		template: result,
	}, nil
}

func (g *Generator) WithDefaults(obj interface{}) (*Generator, error) {
	return g.with(defaultsKey, obj)
}

func (g *Generator) WithResource(obj interface{}) (*Generator, error) {
	return g.with(resourceKey, obj)
}

func (g *Generator) RenderJSON() ([]byte, error) {
	templateKeyPath := cue.ParsePath(templateKey)
	if err := templateKeyPath.Err(); err != nil {
		return nil, err
	}
	value := g.template.LookupPath(templateKeyPath)
	if err := value.Err(); err != nil {
		return nil, fmt.Errorf("unable to lookup %q: %w", templateKey, err)
	}
	return value.MarshalJSON()
}
