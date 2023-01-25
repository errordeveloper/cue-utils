// Copyright 2020 Authors of Cilium
// Copyright 2022 Ilya Dmitrichenko
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"fmt"

	"cuelang.org/go/cue"

	"github.com/errordeveloper/cue-utils/compiler"
	"github.com/errordeveloper/cue-utils/errors"
)

const (
	templateKey = "template"
	defaultsKey = "defaults"
	resourceKey = "resource"
)

type Generator struct {
	dir  string
	args []string
	cue  *compiler.Compiler

	Value      cue.Value
	ImportPath string
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
	val, err := g.cue.BuildAll(g.dir, g.args...)
	if err != nil {
		return err
	}

	g.Value = val.Value
	g.ImportPath = val.ImportPath
	return nil
}

func (g *Generator) with(key string, obj interface{}) (*Generator, error) {
	keyPath := cue.ParsePath(key)
	if err := keyPath.Err(); err != nil {
		return nil, err
	}
	val := g.Value.FillPath(keyPath, obj)
	if err := val.Err(); err != nil {
		return nil, errors.Describe(fmt.Sprintf("unable to fill path %q", key), err)
	}
	return &Generator{
		dir:   g.dir,
		Value: val,
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
	val := g.Value.LookupPath(templateKeyPath)
	if err := val.Err(); err != nil {
		return nil, fmt.Errorf("unable to lookup path %q: %w", templateKey, err)
	}

	data, err := val.MarshalJSON()
	if err != nil {
		return nil, errors.Describe("unable to render JSON", err)
	}
	return data, nil
}
