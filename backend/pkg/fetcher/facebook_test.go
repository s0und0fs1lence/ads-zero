package fetcher

import (
	"fmt"
	"testing"
	"time"

	"github.com/s0und0fs1lence/ads-zero/pkg/common"
	"github.com/s0und0fs1lence/ads-zero/pkg/db"
)

func TestFacebook(t *testing.T) {
	dbConn, err := db.NewClickhouseService(nil)
	if err != nil {
		t.Fatal(err)
	}
	c := NewClient().WithDbSvc(dbConn).WithUserId("300883896877187315")
	task := common.NewFetchTask(time.Now(), time.Now())
	task.Accounts = append(task.Accounts, db.DbAccountSpend{
		Spend: 20,
	})
	data, err := c.ExecuteRules(task)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range data {
		fmt.Printf("data: %v\n", d)
	}
	fmt.Printf("task.Accounts: %v\n", task.Accounts)
}
