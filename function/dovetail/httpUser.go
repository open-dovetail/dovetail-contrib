package dovetail

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/project-flogo/core/data/coerce"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnHTTPUser{})
}

type fnHTTPUser struct {
}

func (fnHTTPUser) Name() string {
	return "httpUser"
}

func (fnHTTPUser) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeParams}, false
}

// Eval - extract username from HTTP headers when basic authentication is used
func (fnHTTPUser) Eval(params ...interface{}) (interface{}, error) {
	headers, err := coerce.ToParams(params[0])
	if err != nil {
		return nil, fmt.Errorf("httpUser function parameter [%+v] must be HTTP headers of type map[string]string", params[0])
	}

	auth, ok := headers["Authorization"]
	if !ok {
		// no user info
		fmt.Println("No authorization header found in the input")
		return "", nil
	}
	if !strings.HasPrefix(auth, "Basic ") || len(auth) < 7 {
		fmt.Printf("Auth header '%s' does not match Basic pattern\n", auth)
		return "", nil
	}

	data := []byte(auth[6:])
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(buf, data)
	if err != nil || n <= 0 {
		fmt.Printf("failed to decode auth token %s: %v\n", data, err)
		return "", nil
	}
	dstr := string(buf[:n])
	tokens := strings.Split(dstr, ":")
	if len(tokens) < 1 {
		fmt.Printf("auth token '%s' is not of format 'user:passwd'\n", dstr)
		return dstr, nil
	}
	return tokens[0], nil
}
