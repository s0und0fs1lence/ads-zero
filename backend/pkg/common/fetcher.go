package common

import (
	"fmt"
	"sync"
	"time"

	"github.com/s0und0fs1lence/ads-zero/pkg/db"
)

type EntityType uint8

const (
	PROVIDER EntityType = iota
	BUSINESS
	ACCOUNT
	CAMPAIGN
	ADSET
	AD
)

type FetchError struct {
	EntityID   string
	EntityType EntityType
	DateRef    time.Time
	Err        error
}

type FetchTask struct {
	mx        sync.Mutex
	Start     time.Time
	End       time.Time
	Accounts  []db.DbAccountSpend
	Campaigns []db.DbCampaignSpend
	Errors    []FetchError
}

func NewFetchTask(start, end time.Time) *FetchTask {
	return &FetchTask{
		mx:        sync.Mutex{},
		Start:     start,
		End:       end,
		Accounts:  make([]db.DbAccountSpend, 0),
		Campaigns: make([]db.DbCampaignSpend, 0),
		Errors:    make([]FetchError, 0),
	}
}

func (t *FetchTask) Lock() {
	t.mx.Lock()
}
func (t *FetchTask) Unlock() {
	t.mx.Unlock()
}

func (t *FetchTask) GetTotalSpend() float64 {
	res := float64(0)
	for _, acc := range t.Accounts {
		res += acc.Spend
	}
	return res
}

// implement fetchTask
func (t *FetchTask) getDailySpendValue() (res float64) {
	for _, acc := range t.Accounts {
		res += acc.Spend
	}
	return res
}

func (t *FetchTask) GetFieldValue(field string) (res interface{}, err error) {

	switch field {
	case "DAILY_SPEND":
		return t.getDailySpendValue(), nil
	case "AVG_CPC":
		return t.getAvgCPCValue(), nil
	case "AVG_CPM":
		return t.getAvgCPMValue(), nil
	case "DAILY_CONVERSIONS":
		return t.getDailyConversionsValue(), nil
	default:
		return nil, fmt.Errorf("invalid field")
	}

}

func (t *FetchTask) getDailyConversionsValue() (res float64) {
	panic("unimplemented")
}

func (t *FetchTask) getAvgCPMValue() (res float64) {
	panic("unimplemented")
}

func (t *FetchTask) getAvgCPCValue() (res float64) {
	panic("unimplemented")
}

func (t *FetchTask) GetFieldAsInt64(field string) (res int64, err error) {
	r, err := t.GetFieldValue(field)
	if err != nil {
		return res, err
	}
	return int64(r.(float64)), nil
}
func (t *FetchTask) GetFieldAsUInt64(field string) (res uint64, err error) {
	r, err := t.GetFieldValue(field)
	if err != nil {
		return res, err
	}
	return uint64(r.(float64)), nil
}
func (t *FetchTask) GetFieldAsFloat64(field string) (res float64, err error) {
	r, err := t.GetFieldValue(field)
	if err != nil {
		return res, err
	}
	return r.(float64), nil
}
