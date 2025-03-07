package rule

import (
	"strings"

	"github.com/s0und0fs1lence/ads-zero/pkg/common"
)

type Column uint

const (
	INVALID Column = iota
	DAILY_SPEND
)

func (c Column) String() string {
	switch c {
	case DAILY_SPEND:
		return "DAILY_SPEND"
	}
	return ""
}

func ColumnFromString(s string) Column {
	s = strings.ToLower(s)
	switch s {
	case "daily_spend":
		return DAILY_SPEND
	default:
		return INVALID
	}
}

type Operator uint

const (
	OpINVALID Operator = iota
	OpROOT
	OpAND
	OpOR
	OpNOT
	OpEQ
	OpNotEQ
	OpLT
	OpLTE
	OpGT
	OpGTE
)

func (o Operator) String() string {
	switch o {
	case OpROOT:
		return "OpROOT"
	case OpAND:
		return "OpAND"
	case OpOR:
		return "OpOR"
	case OpNOT:
		return "OpNOT"
	case OpEQ:
		return "OpEQ"
	case OpNotEQ:
		return "OpNotEQ"
	case OpLT:
		return "OpLT"
	case OpLTE:
		return "OpLTE"
	case OpGT:
		return "OpGT"
	case OpGTE:
		return "OpGTE"
	}
	return "OpINVALID"
}

func (o Operator) GetFunc(valueType any) ConditionFunction {

	switch valueType.(type) {
	case int, int32, int64:
		switch o {
		case OpAND:
			return FuncOpAND
		case OpOR:
			return FuncOpOR
		case OpNOT:
			return FuncOpNOT
		case OpEQ:
			return FuncOpEQ[int64]
		// case OpNotEQ:
		// 	return "OpNotEQ"
		case OpLT:
			return FuncOpLT[int64]
		case OpLTE:
			return FuncOpLTE[int64]
		case OpGT:
			return FuncOpGT[int64]
		case OpGTE:
			return FuncOpGTE[int64]
		}
	case float32, float64:
		switch o {
		case OpAND:
			return FuncOpAND
		case OpOR:
			return FuncOpOR
		case OpNOT:
			return FuncOpNOT
		case OpEQ:
			return FuncOpEQ[float64]
		// case OpNotEQ:
		// 	return "OpNotEQ"
		case OpLT:
			return FuncOpLT[float64]
		case OpLTE:
			return FuncOpLTE[float64]
		case OpGT:
			return FuncOpGT[float64]
		case OpGTE:
			return FuncOpGTE[float64]
		}
	case string:
		switch o {
		case OpAND:
			return FuncOpAND
		case OpOR:
			return FuncOpOR
		case OpNOT:
			return FuncOpNOT
		case OpEQ:
			return FuncOpEQ[string]
		// case OpNotEQ:
		// 	return "OpNotEQ"
		case OpLT:
			return FuncOpLT[string]
		case OpLTE:
			return FuncOpLTE[string]
		case OpGT:
			return FuncOpGT[string]
		case OpGTE:
			return FuncOpGTE[string]
		}
	}

	return nil
}

// Function to match a string to an operator (case insensitive)
func OperatorFromString(s string) Operator {
	s = strings.ToLower(s)
	switch s {
	case "oproot":
		return OpROOT
	case "opand":
		return OpAND
	case "opor":
		return OpOR
	case "opnot":
		return OpNOT
	case "opeq":
		return OpEQ
	case "opnoteq":
		return OpNotEQ
	case "oplt":
		return OpLT
	case "oplte":
		return OpLTE
	case "opgt":
		return OpGT
	case "opgte":
		return OpGTE
	default:
		return OpINVALID
	}
}

type Rule interface {
	Exec(taskResult common.Event) (bool, error)
	Name() string
	Id() string
	Value() interface{}
}
