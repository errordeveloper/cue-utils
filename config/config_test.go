// Copyright 2020 Authors of Cilium
// Copyright 2022 Ilya Dmitrichenko
// SPDX-License-Identifier: Apache-2.0

package config_test

import (
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/errordeveloper/cue-utils/config"
)

func TestLoad(t *testing.T) {
	g := NewGomegaWithT(t)

	{
		err := (&Config{BaseDirectory: "./nonexistent"}).Load()

		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`unable to list avaliable config templates in "./nonexistent": open ./nonexistent: no such file or directory`))
	}

	{
		err := (&Config{BaseDirectory: "./empty"}).Load()

		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`no config templates found in "./empty"`))
	}

	{

		c := &Config{BaseDirectory: "testassets"}

		g.Expect(c.Load()).To(Succeed())

		g.Expect(c.HaveExistingTemplate("basic")).To(BeTrue())
	}

}
