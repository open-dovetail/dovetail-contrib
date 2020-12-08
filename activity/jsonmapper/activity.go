package jsonmapper

import (
	"encoding/json"

	"github.com/project-flogo/core/activity"
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
}

// New creates a new Activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	return &Activity{}, nil
}

// Metadata implements activity.Activity.Metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements activity.Activity.Eval
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
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

	output := &Output{Code: 200, Message: "", Result: string(result)}
	ctx.SetOutputObject(output)
	return true, nil
}
