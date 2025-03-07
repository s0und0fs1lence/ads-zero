package rule

import (
	"testing"
	"time"

	"github.com/s0und0fs1lence/ads-zero/pkg/common"
	"github.com/s0und0fs1lence/ads-zero/pkg/db"
)

func TestSimpleRule(t *testing.T) {
	s := NewSimpleRule(DAILY_SPEND, OpGT, float64(10), "1", "1")
	task := common.NewFetchTask(time.Now(), time.Now())
	task.Accounts = append(task.Accounts, db.DbAccountSpend{
		Spend: 20,
	})
	match, err := s.Exec(task)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatal("Expected match")
	}
}
