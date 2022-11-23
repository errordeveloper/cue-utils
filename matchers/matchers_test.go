// Copyright 2022 Ilya Dmitrichenko
// SPDX-License-Identifier: Apache-2.0

package matchers

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestCUEMatchers(t *testing.T) {
	g := NewWithT(t)
	g.Expect(".").To(
		And(
			BeValidCUEPackage(),
			CUEValueHasLen(1),
			CUEValueMatchesJSON(`{"foo":{"bar":{}}}`),
			CUEValueHasKey("foo"),
			CUEValueHasKeyWithValue("foo", map[string]interface{}{"bar": map[string]interface{}{}}),
			CUEValueHasKeyWithValue("foo", HaveKey("bar")),
			Not(CUEValueHasKey("bar")),
			Not(CUEValueHasKey("_bar")),
		),
	)
	g.Expect("..").ToNot(BeValidCUEPackage())

	g.Expect(`{ "bar": 1 }`).ToNot(MatchCUESchema("bar: string"))
	g.Expect(`{ "bar": 1 }`).ToNot(MatchCUESchema("#Bar: bar: string\n#Bar"))
	g.Expect(`{ "foo": 1 }`).To(MatchCUESchema("bar: string")) // this matches because there is not actual schmea defined (i.e. #Bar)
	g.Expect(`{ "foo": 1 }`).ToNot(MatchCUESchema("#Bar: bar: string\n#Bar"))
	g.Expect(`{ "bar": "", "foo": 1 }`).ToNot(MatchCUESchema("#Bar: bar: string\n#Bar"))

	g.Expect(`{ "bar": 1 }`).To(MatchCUESchema("bar: int"))
	g.Expect(`[]`).ToNot(MatchCUESchema("bar: []"))
	g.Expect(`{ "foo": "bar" }`).To(MatchCUESchema("[_]: string"))
	g.Expect(`{ "foo": "bar" }`).To(MatchCUESchema("foo: string & =~ \"^bar$\""))
	g.Expect(`{ "foo": "brr" }`).ToNot(MatchCUESchema("foo: string & =~ \"^bar$\""))
	g.Expect(`{ }`).To(MatchCUESchema("[_]: string"))
	g.Expect(`{ "foo": [], "bar": {} }`).ToNot(MatchCUESchema("[_]: string"))
}
