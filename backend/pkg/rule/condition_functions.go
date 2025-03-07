package rule

import (
	"cmp"
	"fmt"

	"github.com/s0und0fs1lence/ads-zero/pkg/common"
)

func FuncOpGTE[K cmp.Ordered](c Condition, evt common.Event) (bool, error) {
	field := c.GetTargetField()
	value, err := evt.GetFieldValue(field.String())
	if err != nil {
		return false, err
	}
	targetValue := c.GetValue()
	target, ok := targetValue.(K)
	if !ok {
		return false, fmt.Errorf("invalid target value")
	}

	return value.(K) >= target, nil
}

func FuncOpGT[K cmp.Ordered](c Condition, evt common.Event) (bool, error) {
	field := c.GetTargetField()
	value, err := evt.GetFieldValue(field.String())
	if err != nil {
		return false, err
	}
	targetValue := c.GetValue()
	target, ok := targetValue.(K)
	if !ok {
		return false, fmt.Errorf("invalid target value")
	}

	return value.(K) > target, nil
}

func FuncOpLTE[K cmp.Ordered](c Condition, evt common.Event) (bool, error) {
	field := c.GetTargetField()
	value, err := evt.GetFieldValue(field.String())
	if err != nil {
		return false, err
	}
	targetValue := c.GetValue()
	target, ok := targetValue.(K)
	if !ok {
		return false, fmt.Errorf("invalid target value")
	}

	return value.(K) <= target, nil
}
func FuncOpLT[K cmp.Ordered](c Condition, evt common.Event) (bool, error) {
	field := c.GetTargetField()
	value, err := evt.GetFieldValue(field.String())
	if err != nil {
		return false, err
	}
	targetValue := c.GetValue()
	target, ok := targetValue.(K)
	if !ok {
		return false, fmt.Errorf("invalid target value")
	}

	return value.(K) < target, nil
}

func FuncOpEQ[K cmp.Ordered](c Condition, evt common.Event) (bool, error) {
	field := c.GetTargetField()
	value, err := evt.GetFieldValue(field.String())
	if err != nil {
		return false, err
	}
	targetValue := c.GetValue()
	target, ok := targetValue.(K)
	if !ok {
		return false, fmt.Errorf("invalid target value")
	}

	return value.(K) == target, nil
}

func FuncOpAND(c Condition, evt common.Event) (bool, error) {
	for _, child := range c.GetChildrens() {
		res, err := child.Exec(evt)
		if err != nil {
			return false, err
		}
		if !res {
			return false, nil
		}
	}
	return true, nil
}

func FuncOpOR(c Condition, evt common.Event) (bool, error) {
	for _, child := range c.GetChildrens() {
		res, err := child.Exec(evt)
		if err != nil {
			return false, err
		}
		if res {
			return true, nil
		}
	}
	return false, nil
}

func FuncOpNOT(c Condition, evt common.Event) (bool, error) {
	for _, child := range c.GetChildrens() {
		res, err := child.Exec(evt)
		if err != nil {
			return false, err
		}
		if res {
			return false, nil
		}
	}
	return true, nil
}
