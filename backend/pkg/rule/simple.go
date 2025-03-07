package rule

import (
	"github.com/s0und0fs1lence/ads-zero/pkg/common"
)

type simpleRule struct {
	name      string
	id        string
	condition Condition
}

// Id implements Rule.
func (s *simpleRule) Id() string {
	return s.id
}

// Name implements Rule.
func (s *simpleRule) Name() string {
	return s.name
}

// Name implements Rule.
func (s *simpleRule) Value() interface{} {
	return s.condition.GetValue()
}

// Exec implements Rule.
func (s *simpleRule) Exec(taskResult common.Event) (bool, error) {

	return s.condition.Exec(taskResult)

}

func NewSimpleRule(column Column, operator Operator, value interface{}, name, id string) Rule {
	cond, err := NewConditionLeaf(operator)
	if err != nil {
		return nil
	}
	cond.SetTargetField(column)
	cond.SetValue(value)
	s := &simpleRule{
		name:      name,
		id:        id,
		condition: cond,
	}
	return s
}
