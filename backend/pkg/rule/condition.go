package rule

import (
	"github.com/s0und0fs1lence/ads-zero/pkg/common"
)

type ConditionFunction func(c Condition, evt common.Event) (bool, error)

type Condition interface {
	Append(child Condition) (parent Condition)
	ClearChildrens()
	//TODO: get valid input for the func
	Exec(evt common.Event) (bool, error)
	GetOperator() Operator
	GetValue() interface{}
	SetValue(value interface{})
	SetFunction(fn ConditionFunction)
	SetTargetField(field Column)
	GetTargetField() Column
	GetChildrens() []Condition
}
