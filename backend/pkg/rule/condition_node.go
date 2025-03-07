package rule

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/s0und0fs1lence/ads-zero/pkg/common"
)

type ConditionNode struct {
	Operator  Operator
	function  ConditionFunction
	Childrens []Condition
}

// SetValue implements Condition.
func (cn *ConditionNode) SetValue(value interface{}) {

}

func NewConditionNode(op Operator, fn ConditionFunction) (Condition, error) {
	if op == OpINVALID {
		log.Error().Str("function", "NewConditionNode").Msg("Invalid operator")
		return nil, fmt.Errorf("invalid operator")
	}
	return &ConditionNode{
		Operator:  op,
		function:  fn,
		Childrens: make([]Condition, 0),
	}, nil
}

// GetTargetField implements Condition.
func (cn *ConditionNode) GetValue() interface{} {
	return nil
}

// GetTargetField implements Condition.
func (cn *ConditionNode) GetTargetField() Column {
	return INVALID
}

// SetTargetField implements Condition.
func (cn *ConditionNode) SetTargetField(field Column) {
}

// Implement the Condition interface for ConditionNode
func (cn *ConditionNode) Append(child Condition) Condition {
	cn.Childrens = append(cn.Childrens, child)
	return cn
}

// GetChildrens implements Condition.
func (cn *ConditionNode) GetChildrens() []Condition {
	return cn.Childrens
}

func (cn *ConditionNode) ClearChildrens() {
	cn.Childrens = make([]Condition, 0)
}

func (cn *ConditionNode) Exec(evt common.Event) (bool, error) {
	return cn.function(cn, evt)
}

func (cn *ConditionNode) GetOperator() Operator {
	return cn.Operator
}

func (cn *ConditionNode) SetFunction(fn ConditionFunction) {
	cn.function = fn
}
