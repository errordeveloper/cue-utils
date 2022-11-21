// Copyright 2022 Ilya Dmitrichenko
// SPDX-License-Identifier: Apache-2.0

package compiler

import (
	"encoding/json"
	"fmt"
	"sync"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"

	"github.com/errordeveloper/cue-utils/errors"
)

var sharedCUEMutex = &sync.Mutex{
	// load.Instances is not thread-safe (https://github.com/cue-lang/cue/issues/1043#issuecomment-1016729326)
}

type Compiler struct {
	ctx   *cue.Context
	mutex *sync.Mutex
}

func NewCompiler() *Compiler {
	return &Compiler{
		ctx:   cuecontext.New(),
		mutex: sharedCUEMutex,
	}
}

func (c *Compiler) BuildAll(dir string, args ...string) (cue.Value, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	loadedInstances := load.Instances(args, &load.Config{Dir: dir})
	for _, loadedInstance := range loadedInstances {
		if loadedInstance.Err != nil {
			return cue.Value{}, errors.Describe(fmt.Sprintf("failed to load instances (dir: %q, args: %v)", dir, args), loadedInstance.Err)
		}
	}

	builtInstances, err := c.ctx.BuildInstances(loadedInstances)
	if err != nil {
		return cue.Value{}, errors.Describe(fmt.Sprintf("failed to build instances (dir: %q, args: %v)", dir, args), err)
	}
	for _, builtInstance := range builtInstances {
		if err := builtInstance.Value().Validate(); err != nil {
			return cue.Value{}, errors.Describe("validation failure", err)
		}
	}
	if len(builtInstances) != 1 {
		return cue.Value{}, fmt.Errorf("unexpected: more then one instance built")
	}
	return builtInstances[0], nil
}

func (c *Compiler) CompileString(src string, options ...cue.BuildOption) cue.Value {
	return c.ctx.CompileString(src, options...)
}

func (c *Compiler) MarshalValueJSON(v cue.Value) ([]byte, error) {
	return json.Marshal(v)
}
