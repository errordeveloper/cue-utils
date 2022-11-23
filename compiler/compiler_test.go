// Copyright 2022 Ilya Dmitrichenko
// SPDX-License-Identifier: Apache-2.0

package compiler_test

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/errordeveloper/cue-utils/compiler"
)

func TestCUEBuildAll(t *testing.T) {
	g := NewWithT(t)
	val, err := NewCompiler().BuildAll("", ".")
	g.Expect(err).ToNot(HaveOccurred())
	data, err := json.Marshal(val)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(data).To(MatchJSON(`{"foo":{"bar":{}}}`))
}
