package dovetail

import (
	"testing"

	"github.com/project-flogo/core/data/coerce"
	"github.com/stretchr/testify/assert"
)

func TestFnHTTPUser_Eval(t *testing.T) {
	f := fnHTTPUser{}
	v, err := f.Eval(map[string]string{
		"Authorization": "Basic ZGVtb1VzZXI6ZGVtb1Bhc3M=",
	})
	assert.NoError(t, err, "eval httpUser should not throw error")
	u, _ := coerce.ToString(v)
	assert.Equal(t, "demoUser", u, "username should be 'demoUser'")
}
