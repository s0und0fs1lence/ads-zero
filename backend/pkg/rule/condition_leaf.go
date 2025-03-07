package rule

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/s0und0fs1lence/ads-zero/pkg/common"
)

type ConditionLeaf struct {
	Operator    Operator
	function    ConditionFunction
	targetField Column
	value       any
}

// SetValue implements Condition.
func (cl *ConditionLeaf) SetValue(value interface{}) {
	cl.value = value
	cl.function = cl.Operator.GetFunc(value)
}

func NewConditionLeaf(op Operator) (Condition, error) {
	if op == OpINVALID {
		log.Error().Str("function", "NewConditonLeaf").Msg("Invalid operator")
		return nil, fmt.Errorf("invalid operator")
	}
	return &ConditionLeaf{
		Operator: op,
		function: op.GetFunc(nil),
	}, nil
}

func (cl *ConditionLeaf) GetValue() interface{} {
	return cl.value
}

// SetTargetField implements Condition.
func (cl *ConditionLeaf) SetTargetField(field Column) {
	cl.targetField = field
}

// GetTargetField implements Condition.
func (cn *ConditionLeaf) GetTargetField() Column {
	return cn.targetField
}

// GetChildrens implements Condition.
func (cl *ConditionLeaf) GetChildrens() []Condition {
	return []Condition{}
}

// Implement the Condition interface for ConditionLeaf
func (cl *ConditionLeaf) Append(child Condition) Condition {
	// ConditionLeaf cannot have children, so return itself
	return cl
}

func (cl *ConditionLeaf) ClearChildrens() {
	// ConditionLeaf cannot have children, so do nothing
}

func (cl *ConditionLeaf) Exec(evt common.Event) (bool, error) {
	return cl.function(cl, evt)
}

func (cl *ConditionLeaf) GetOperator() Operator {
	return cl.Operator
}

func (cl *ConditionLeaf) SetFunction(fn ConditionFunction) {
	cl.function = fn
}
