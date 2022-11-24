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
		g.Expect(err.Error()).To(Equal(`unable to list avaliable config templates in "./nonexistent": lstat ./nonexistent: no such file or directory`))
	}

	{
		err := (&Config{BaseDirectory: "./testassets/empty"}).Load()

		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`no config templates found in "./testassets/empty"`))
	}

	{
		c := &Config{BaseDirectory: "testassets"}

		g.Expect(c.Load()).To(Succeed())

		g.Expect(c.HaveExistingTemplate("github.com/errordeveloper/cue-utils/config/testassets/basic")).To(BeTrue())
		g.Expect(c.HaveExistingTemplate("github.com/errordeveloper/cue-utils/config/testassets/nested:nested1")).To(BeTrue())
		g.Expect(c.HaveExistingTemplate("github.com/errordeveloper/cue-utils/config/testassets/nested/2:nested2")).To(BeTrue())
		g.Expect(c.HaveExistingTemplate("github.com/errordeveloper/cue-utils/config/testassets/empty")).To(BeFalse())

	}
}

func TestLoadWithImortPaths(t *testing.T) {
	g := NewGomegaWithT(t)

	{
		c := &Config{BaseDirectory: "./"}
		err := (c).Load()

		g.Expect(err).To(Not(HaveOccurred()))

		for _, template := range c.ExistingTemplates() {
			template, err := c.Get(template)
			g.Expect(err).To(Not(HaveOccurred()))
			t.Logf("%#v\n", template.ImportPath)
		}
	}

}
