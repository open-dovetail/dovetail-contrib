package jsonmapper

import (
	"encoding/json"
	"fmt"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
)

// Create a new logger
var logger = log.ChildLogger(log.RootLogger(), "activity-json-mapper")

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Activity is a stub for executing Hyperledger Fabric put operations
type Activity struct {
	compositeKeys map[string]interface{}
	name          string
	key           string
	attributes    []string
}

// New creates a new Activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	ctx.Logger().Infof("Create activity with InitContxt settings %v\n", ctx.Settings())
	if err := metadata.MapToStruct(ctx.Settings(), s, true); err != nil {
		ctx.Logger().Errorf("failed to configure activity %v\n", err)
		return nil, err
	}
	mapping, ok := s.CompositeKeys["mapping"].(map[string]interface{})
	var ck string
	var attrs []string
	if ok {
		for k, v := range mapping {
			ck = k
			values := v.([]interface{})
			for _, f := range values {
				attrs = append(attrs, f.(string))
			}
		}
	}
	return &Activity{
		compositeKeys: s.CompositeKeys,
		name:          s.Name,
		key:           ck,
		attributes:    attrs,
	}, nil
}

// Metadata implements activity.Activity.Metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements activity.Activity.Eval
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	logger.Infof("activity setting name %s keys %v keyName %s, keyAttrs: %v", a.name, a.compositeKeys, a.key, a.attributes)
	// check input args
	input := &Input{}
	if err = ctx.GetInputObject(input); err != nil {
		logger.Errorf("failed to resolve input data: %+v\n", err)
		output := &Output{Code: 400, Message: err.Error(), Result: ""}
		ctx.SetOutputObject(output)
		return false, err
	}

	if len(input.Data) == 0 {
		logger.Warnf("input data is empty\n")
		output := &Output{Code: 300, Message: "input data is empty", Result: ""}
		ctx.SetOutputObject(output)
		return true, nil
	}

	result, err := json.Marshal(input.Data)
	if err != nil {
		logger.Errorf("failed to serialize data: %+v\n", err)
		output := &Output{Code: 400, Message: err.Error(), Result: ""}
		ctx.SetOutputObject(output)
		return false, err
	}

	output := &Output{Code: 200, Message: fmt.Sprintf("%s - %v", a.name, a.compositeKeys), Result: string(result)}
	ctx.SetOutputObject(output)
	return true, nil
}
