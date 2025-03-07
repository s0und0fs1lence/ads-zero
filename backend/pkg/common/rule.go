package common

type Event interface {
	GetFieldValue(field string) (interface{}, error)
	GetFieldAsInt64(field string) (int64, error)
	GetFieldAsUInt64(field string) (uint64, error)
	GetFieldAsFloat64(field string) (float64, error)
}

type RuleResult struct {
	RuleName  string
	RuleID    string
	Result    bool
	Threshold any
}
