// Copyright 2022 Ilya Dmitrichenko
// SPDX-License-Identifier: Apache-2.0

package matchers

import (
	"encoding/json"
	"fmt"

	"cuelang.org/go/cue"
	"github.com/errordeveloper/cue-utils/compiler"
	"github.com/errordeveloper/cue-utils/errors"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"
)

func init() {
	format.MaxLength = 10000
}

func BeValidCUEPackage() types.GomegaMatcher {
	return &beValidCUEPackageMatcher{
		cue: compiler.NewCompiler(),
	}
}

type beValidCUEPackageMatcher struct {
	reasonErr error
	cue       *compiler.Compiler
}

func (m *beValidCUEPackageMatcher) Match(actual interface{}) (bool, error) {
	_, err := m.cue.BuildAll("", actual.(string))
	m.reasonErr = errors.Describe("failed to build", err)
	return err == nil, nil
}

func (m *beValidCUEPackageMatcher) FailureMessage(actual interface{}) string {
	msg := format.Message(actual, "to be a valid CUE package")
	if m.reasonErr != nil {
		msg += "\nerror: " + m.reasonErr.Error()
	}
	return msg
}

func (m *beValidCUEPackageMatcher) NegatedFailureMessage(actual interface{}) string {
	msg := format.Message(actual, "to be an invalid CUE package")
	if m.reasonErr != nil {
		msg += "\nerror: " + m.reasonErr.Error()
	}
	return msg
}

func CUEValueMatchesJSON(json interface{}) types.GomegaMatcher {
	return &cueValueMatchesJSONMatcher{
		MatchJSONMatcher: &matchers.MatchJSONMatcher{
			JSONToMatch: json,
		},
		cue: compiler.NewCompiler(),
	}
}

type cueValueMatchesJSONMatcher struct {
	*matchers.MatchJSONMatcher
	cue        *compiler.Compiler
	actualJSON []byte
}

func (m *cueValueMatchesJSONMatcher) Match(actual interface{}) (bool, error) {
	value, err := m.cue.BuildAll("", actual.(string))
	if err != nil {
		return false, err
	}
	actualJSON, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	m.actualJSON = actualJSON
	return m.MatchJSONMatcher.Match(m.actualJSON)
}

func (m *cueValueMatchesJSONMatcher) FailureMessage(actual interface{}) string {
	return m.MatchJSONMatcher.FailureMessage(m.actualJSON)
}

func (m *cueValueMatchesJSONMatcher) NegatedFailureMessage(actual interface{}) string {
	return m.MatchJSONMatcher.NegatedFailureMessage(m.actualJSON)
}

func newCUEValueHelper() *cueValueHelper {
	return &cueValueHelper{
		cue: compiler.NewCompiler(),
	}
}

type cueValueHelper struct {
	cue *compiler.Compiler
	dir string
	obj interface{}
}

func (c *cueValueHelper) matchWith(input string, m types.GomegaMatcher) (bool, error) {
	value, err := c.cue.BuildAll(c.dir, input)
	if err != nil {
		return false, err
	}
	data, err := value.MarshalJSON()
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(data, &c.obj); err != nil {
		return false, err
	}
	return m.Match(c.obj)
}

func CUEValueHasLen(count int) types.GomegaMatcher {
	return &cueValueHasLen{
		HaveLenMatcher: &matchers.HaveLenMatcher{
			Count: count,
		},
		cueValueHelper: newCUEValueHelper(),
	}
}

type cueValueHasLen struct {
	*cueValueHelper
	*matchers.HaveLenMatcher
}

func (m *cueValueHasLen) Match(actual interface{}) (bool, error) {
	return m.matchWith(actual.(string), m.HaveLenMatcher)
}

func (m *cueValueHasLen) FailureMessage(interface{}) string {
	return m.HaveLenMatcher.FailureMessage(m.obj)
}

func (m *cueValueHasLen) NegatedFailureMessage(interface{}) string {
	return m.HaveLenMatcher.NegatedFailureMessage(m.obj)
}

func CUEValueHasKey(key interface{}) types.GomegaMatcher {
	return &cueValueHasKey{
		HaveKeyMatcher: &matchers.HaveKeyMatcher{
			Key: key,
		},
		cueValueHelper: newCUEValueHelper(),
	}
}

type cueValueHasKey struct {
	*cueValueHelper
	*matchers.HaveKeyMatcher
}

func (m *cueValueHasKey) Match(actual interface{}) (bool, error) {
	return m.matchWith(actual.(string), m.HaveKeyMatcher)
}

func (m *cueValueHasKey) FailureMessage(interface{}) string {
	return m.HaveKeyMatcher.FailureMessage(m.obj)
}

func (m *cueValueHasKey) NegatedFailureMessage(interface{}) string {
	return m.HaveKeyMatcher.NegatedFailureMessage(m.obj)
}

func CUEValueHasKeyWithValue(key, value interface{}) types.GomegaMatcher {
	return &cueValueHasKeyWithValue{
		HaveKeyWithValueMatcher: &matchers.HaveKeyWithValueMatcher{
			Key:   key,
			Value: value,
		},
		cueValueHelper: newCUEValueHelper(),
	}
}

type cueValueHasKeyWithValue struct {
	*cueValueHelper
	*matchers.HaveKeyWithValueMatcher
}

func (m *cueValueHasKeyWithValue) Match(actual interface{}) (bool, error) {
	return m.matchWith(actual.(string), m.HaveKeyWithValueMatcher)
}

func (m *cueValueHasKeyWithValue) FailureMessage(interface{}) string {
	return m.HaveKeyWithValueMatcher.FailureMessage(m.obj)
}

func (m *cueValueHasKeyWithValue) NegatedFailureMessage(interface{}) string {
	return m.HaveKeyWithValueMatcher.NegatedFailureMessage(m.obj)
}

func MatchCUESchema(schema string) *matchCUESchemaMatcher {
	return &matchCUESchemaMatcher{
		schema: schema,
		cue:    compiler.NewCompiler(),
	}
}

func (m *matchCUESchemaMatcher) WithImportPath(v string) *matchCUESchemaMatcher {
	m.importPath = v
	return m
}

func (m *matchCUESchemaMatcher) WithWorkDir(v string) *matchCUESchemaMatcher {
	m.workDir = v
	return m
}

var _ types.GomegaMatcher = &matchCUESchemaMatcher{}

type matchCUESchemaMatcher struct {
	importPath, schema string
	workDir            string
	reasonErr          error
	cue                *compiler.Compiler
}

func (m *matchCUESchemaMatcher) Match(actual interface{}) (bool, error) {
	var schema cue.Value
	if m.importPath != "" {
		imported, err := m.cue.BuildAll(m.workDir, m.importPath)
		if err != nil {
			return false, err
		}
		schema = imported.LookupPath(cue.ParsePath(m.schema))
	} else {
		schema = m.cue.CompileString(m.schema, cue.Filename("test_schema.cue"))
	}
	if schema.Err() != nil {
		return false, errors.Describe("failed to compile schema", schema.Err())
	}

	jsonValue := ""
	switch actual := actual.(type) {
	case string:
		jsonValue = actual
	case []byte:
		jsonValue = string(actual)
	default:
		return false, fmt.Errorf("unexpected type %T", actual)
	}

	valueDefinition := fmt.Sprintf("import \"encoding/json\"\njson.Unmarshal(%q)", jsonValue)
	value := m.cue.CompileString(valueDefinition, cue.Filename("test_value.cue"))
	if value.Err() != nil {
		return false, errors.Describe("failed to compile value", value.Err())
	}
	unified := schema.Unify(value)
	if unified.Err() != nil {
		m.reasonErr = errors.Describe("failed to unify", unified.Err())
		return false, nil
	}
	if err := unified.Validate(); err != nil {
		m.reasonErr = errors.Describe("failed to unify", unified.Err())
	}
	return m.reasonErr == nil, nil
}

func (m *matchCUESchemaMatcher) FailureMessage(actual interface{}) string {
	if b, ok := actual.([]byte); ok {
		actual = string(b)
	}
	msg := format.Message(actual, "to match CUE schema", m.schema)
	if m.reasonErr != nil {
		msg += "\nerror: " + m.reasonErr.Error()
	}
	return msg
}

func (m *matchCUESchemaMatcher) NegatedFailureMessage(actual interface{}) string {
	if b, ok := actual.([]byte); ok {
		actual = string(b)
	}
	msg := format.Message(actual, "to NOT match CUE schema", m.schema)
	if m.reasonErr != nil {
		msg += "\nerror: " + m.reasonErr.Error()
	}
	return msg
}
