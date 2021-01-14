/*
SPDX-License-Identifier: BSD-3-Clause-Open-MPI
*/

package jsonmapper

import (
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}

func TestCreate(t *testing.T) {

	mf := mapper.NewFactory(resolve.GetBasicResolver())
	iCtx := test.NewActivityInitContext(Settings{}, mf)
	act, err := New(iCtx)
	assert.Nil(t, err)
	assert.NotNil(t, act, "activity should not be nil")
}

func TestEval(t *testing.T) {
	act := &Activity{}
	act.serialize = true
	tc := test.NewActivityContext(act.Metadata())
	input := &Input{Data: map[string]interface{}{"two": 2, "one": 1}}
	err := tc.SetInputObject(input)
	assert.Nil(t, err)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.Nil(t, err)
	assert.Equal(t, 200, output.Code)

	// Note that Golang map is serialized with keys sorted lexicographically in JSON
	// ref: https://golang.org/src/encoding/json/encode.go line 793
	assert.Equal(t, "{\"one\":1,\"two\":2}", output.Result)
}
