package fetcher

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	fb "github.com/huandu/facebook/v2"
	"github.com/s0und0fs1lence/ads-zero/pkg/common"
	"github.com/s0und0fs1lence/ads-zero/pkg/configuration"
	"github.com/s0und0fs1lence/ads-zero/pkg/db"
)

type fbRequest struct {
	bUrl   string
	params fb.Params
}

type campaignInsights struct {
	campaignID   string
	campaignName string
	spend        float64
}
type accountInsights struct {
	account   accountInfo
	start     time.Time
	end       time.Time
	dateRef   time.Time
	campaigns map[string]campaignInsights
	spend     float64
}

var (
	allRequests = map[string]fbRequest{
		"get_all_business_managers": {
			bUrl: "/me/businesses",
			params: fb.Params{
				"fields": "id,link,name,is_hidden,profile_picture_uri,created_time,verification_status",
			},
		},
		"owned_ad_accounts": {
			bUrl: "owned_ad_accounts",
			params: fb.Params{
				"fields": "business,id,account_id,account_status,currency,created_time,owner,timezone_id,timezone_name,timezone_offset_hours_utc,name",
				"ids":    "",
				"limit":  50000000,
			},
		},
		"client_ad_accounts": {
			bUrl: "client_ad_accounts",
			params: fb.Params{
				"fields": "business,id,account_id,account_status,currency,created_time,owner,timezone_id,timezone_name,timezone_offset_hours_utc,name",
				"ids":    "",
				"limit":  50000000,
			},
		},
		"account_insights": {
			bUrl: "",
			params: fb.Params{
				"fields":     "spend", //TODO: add more fields if needed
				"time_range": "",
				"level":      "account",
				"limit":      500000,
			},
		},
	}
	//TODO: add app id and app secret
	globalApp = fb.New(
		configuration.Config().GetString(configuration.FacebookAppID),
		configuration.Config().GetString(configuration.FacebookAppSecret),
	)
)

type businessManagerInfo struct {
	Id                string `facebook:",required"`
	Name              string `facebook:"name"`
	Link              string `facebook:"link"`
	IsHidden          bool   `facebook:"is_hidden"`
	ProfilePictureUri string `facebook:"profile_picture_uri"`
	// CreatedBy         struct {
	// 	id   string `facebook:"id"`
	// 	name string `facebook:"name"`
	// } `facebook:"created_by"`

	CreatedTime        time.Time `facebook:"created_time"`
	VerificationStatus string    `facebook:"verification_time"`
}

type tRange struct {
	Since string `json:"since"`
	Until string `json:"until"`
}

type accountInfo struct {
	Id       string `facebook:",required"`
	Name     string `facebook:"name"`
	Business struct {
		Id   string
		Name string
	}
	AccountId          string
	Status             string
	Currency           string
	CreatedTime        time.Time `facebook:"created_time"`
	VerificationStatus string    `facebook:"verification_time"`
}

func fetchAllPages(paging *fb.PagingResult) []fb.Result {
	var dest []fb.Result
	dest = append(dest, paging.Data()...)
	for {
		// get next page.
		noMore, err := paging.Next()
		if err != nil {
			panic(err)
		}
		if noMore {
			// No more results available
			break
		}
		// append current page of results to slice of Result
		dest = append(dest, paging.Data()...)
	}
	return dest
}

func fetchBusinessMenagers(task *common.FetchTask, session *fb.Session) ([]businessManagerInfo, error) {

	curReq := allRequests["get_all_business_managers"]

	res, err := session.Get(curReq.bUrl, curReq.params)
	if err != nil {
		return nil, err
	}
	paging, err := res.Paging(session)
	if err != nil {
		return nil, err
	}

	data := fetchAllPages(paging)

	toReturn := make([]businessManagerInfo, len(data))
	for idx, item := range data {
		var elem businessManagerInfo
		if err := item.Decode(&elem); err != nil {
			task.Errors = append(task.Errors, common.FetchError{
				EntityType: common.BUSINESS,
				Err:        err,
			})
			continue
		}
		toReturn[idx] = elem
	}

	return toReturn, nil
}

func parseAccountInfo(task *common.FetchTask, data []interface{}) []accountInfo {
	var ok bool
	result := make([]accountInfo, 0)
	for _, elem := range data {
		m, isValid := elem.(map[string]interface{})
		if !isValid {
			continue
		}

		t := accountInfo{}
		t.Id, ok = m["id"].(string)
		if !ok {
			continue
		}
		t.AccountId, ok = m["account_id"].(string)
		if !ok {
			continue
		}
		t.Name, ok = m["name"].(string)
		if !ok {
			task.Lock()
			task.Errors = append(task.Errors, common.FetchError{
				EntityID:   t.AccountId,
				EntityType: common.ACCOUNT,
				Err:        fmt.Errorf("the `name` field was not present in the json reply"),
			})
			task.Unlock()
			continue
		}
		status, ok := m["account_status"].(json.Number)
		if !ok {
			task.Lock()
			task.Errors = append(task.Errors, common.FetchError{
				EntityID:   t.AccountId,
				EntityType: common.ACCOUNT,
				Err:        fmt.Errorf("the `account_status` field was not present in the json reply"),
			})
			task.Unlock()
			continue
		}
		switch status.String() {
		case "1":
			t.Status = "ACTIVE"
		case "2", "101":
			t.Status = "INACTIVE"
		default:
			t.Status = "UNKNOWN"
		}
		// t.Status = status.String()
		t.Currency, ok = m["currency"].(string)
		if !ok {
			task.Lock()
			task.Errors = append(task.Errors, common.FetchError{
				EntityID:   t.AccountId,
				EntityType: common.ACCOUNT,
				Err:        fmt.Errorf("the `currency` field was not present in the json reply"),
			})
			task.Unlock()
			continue
		}
		bmMap, ok := m["business"].(map[string]interface{})
		if !ok {
			task.Lock()
			task.Errors = append(task.Errors, common.FetchError{
				EntityID:   t.AccountId,
				EntityType: common.ACCOUNT,
				Err:        fmt.Errorf("the `business` field was not present in the json reply"),
			})
			task.Unlock()
			continue
		}
		t.Business.Id, ok = bmMap["id"].(string)
		if !ok {
			task.Lock()
			task.Errors = append(task.Errors, common.FetchError{
				EntityID:   t.AccountId,
				EntityType: common.ACCOUNT,
				Err:        fmt.Errorf("the `business.id` field was not present in the json reply"),
			})
			task.Unlock()
			continue
		}
		t.Business.Name, ok = bmMap["name"].(string)
		if !ok {
			task.Lock()
			task.Errors = append(task.Errors, common.FetchError{
				EntityID:   t.AccountId,
				EntityType: common.ACCOUNT,
				Err:        fmt.Errorf("the `business.name` field was not present in the json reply"),
			})
			task.Unlock()
			continue
		}
		ctime, ok := m["created_time"].(string)
		if !ok {
			task.Lock()
			task.Errors = append(task.Errors, common.FetchError{
				EntityID:   t.AccountId,
				EntityType: common.ACCOUNT,
				Err:        fmt.Errorf("the `created_time` field was not present in the json reply"),
			})
			task.Unlock()
			continue
		}
		created_time, err := time.Parse(time.RFC3339, ctime)
		if err != nil {
			continue
		}
		t.CreatedTime = created_time
		//TODO: parse all fields

		result = append(result, t)
	}
	return result
}

func fetchAccountsByBm(task *common.FetchTask, session *fb.Session, bmInfo []businessManagerInfo) ([]accountInfo, error) {
	result := make([]accountInfo, 0)
	mx := sync.Mutex{}
	bmIds := make([]string, len(bmInfo))
	for idx := range bmInfo {
		bmIds[idx] = bmInfo[idx].Id
		// bmIds = append(bmIds, bmInfo[idx].Id)
	}
	g := new(errgroup.Group)

	g.Go(func() error {
		curReq := allRequests["owned_ad_accounts"]
		curReq.params["ids"] = strings.Join(bmIds, ",")

		res, err := session.Get(curReq.bUrl, curReq.params)
		if err != nil {
			return err
		}

		for k, v := range res {
			if !strings.Contains(k, "include") {
				mp, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				dt, ok := mp["data"]
				if !ok {
					continue
				}
				lst, ok := dt.([]interface{})
				if !ok {
					continue
				}
				info := parseAccountInfo(task, lst)

				mx.Lock()
				result = append(result, info...)

				mx.Unlock()
			}

		}

		return nil
	})
	g.Go(func() error {
		curReq := allRequests["client_ad_accounts"]
		curReq.params["ids"] = strings.Join(bmIds, ",")

		res, err := session.Get(curReq.bUrl, curReq.params)
		if err != nil {
			return err
		}

		for k, v := range res {
			if !strings.Contains(k, "include") {
				mp, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				dt, ok := mp["data"]
				if !ok {
					continue
				}
				lst, ok := dt.([]interface{})
				if !ok {
					continue
				}

				info := parseAccountInfo(task, lst)
				mx.Lock()
				result = append(result, info...)

				mx.Unlock()
			}

		}
		return nil
	})
	// Wait for all HTTP fetches to complete.
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}

func getDateRef(start, end interface{}) (*time.Time, error) {
	start_str, ok := start.(string)
	if !ok {
		return nil, fmt.Errorf("could not convert start to string")
	}
	end_str, ok := end.(string)
	if !ok {
		return nil, fmt.Errorf("could not convert end to string")
	}
	start_time, err := time.Parse(time.DateOnly, start_str)
	if err != nil {
		return nil, err
	}
	end_time, err := time.Parse(time.DateOnly, end_str)
	if err != nil {
		return nil, err
	}
	if start_time.Equal(end_time) {
		return &start_time, nil
	}
	return nil, fmt.Errorf("the date differ")
}

func fetchAccountSpend(session *fb.Session, accountInfo []accountInfo, start, end time.Time) ([]accountInsights, error) {
	g := new(errgroup.Group)
	mx := sync.Mutex{}
	res := make([]accountInsights, 0)
	timeRange := tRange{
		Since: start.Format("2006-01-02"),
		Until: end.Format("2006-01-02"),
	}

	chunks := slices.Chunk(accountInfo, 50)
	for c := range chunks {
		g.Go(func() error {
			for _, info := range c {
				req := fbRequest{
					bUrl: fmt.Sprintf("%s/insights", info.Id),
					params: fb.Params{
						"fields":     "spend,campaign_id,campaign_name", //TODO: add more fields if needed
						"time_range": fmt.Sprintf("{'since':'%s','until': '%s'}", timeRange.Since, timeRange.Until),
						"level":      "campaign",
						"limit":      500000,
					},
				}
				response, err := session.Get(req.bUrl, req.params)
				if err != nil {
					return err
				}
				paging, err := response.Paging(session)
				if err != nil {
					return err
				}
				toRet := accountInsights{
					account:   info,
					start:     start,
					end:       end,
					campaigns: make(map[string]campaignInsights),
					spend:     float64(0),
				}

				finalSpend := float64(0)
				data := fetchAllPages(paging)

				date_ref, err := getDateRef(start.Format(time.DateOnly), end.Format(time.DateOnly))
				if err != nil {
					return err
				}
				toRet.dateRef = *date_ref
				for _, pg := range data {
					date_start := pg.GetField("date_start")
					date_stop := pg.GetField("date_stop")
					date_ref, err := getDateRef(date_start, date_stop)
					if err != nil {
						return err
					}
					toRet.dateRef = *date_ref
					sp := pg.GetField("spend")
					var currSpend float64
					if sp != nil {
						strRep, ok := sp.(string)
						if !ok {
							continue
						}
						currSpend, err = strconv.ParseFloat(strRep, 64)
						if err != nil {
							return err
						}
						finalSpend += currSpend

					}
					campaign_id := pg.GetField("campaign_id")
					if campaign_id != nil {
						campaign_name := pg.GetField("campaign_name")
						c, exist := toRet.campaigns[campaign_id.(string)]
						if exist {

							c.campaignName = campaign_name.(string)
							c.campaignID = campaign_id.(string)
							c.spend = currSpend
							toRet.campaigns[campaign_id.(string)] = c
						} else {
							i := campaignInsights{
								campaignID:   campaign_id.(string),
								campaignName: campaign_name.(string),
								spend:        currSpend,
							}
							toRet.campaigns[campaign_id.(string)] = i
						}

					}
				}
				toRet.spend = finalSpend

				mx.Lock()
				res = append(res, toRet)
				mx.Unlock()
			}
			return nil
		})
	}
	// Wait for all HTTP fetches to complete.
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return res, nil
}

// TODO:
// current implementation is lacking some advanced stuff like rate limiting retry etc.
// also, we should define a structure wich we should ideally send into a channel o return as a whole list, and after insert this data in the database
// we should also improve error checking, and maybe tweek the goroutines for smaller clients.

func facebookFetcher(accessToken, appID, appSecret string, start, end time.Time) (tasks *common.FetchTask, err error) {
	app := fb.New(
		appID,
		appSecret,
	)
	session := app.Session(accessToken)
	if session == nil {
		return nil, fmt.Errorf("could not use the provided access token")
	}
	session.RFC3339Timestamps = true
	session.Version = "v21.0"
	task := common.NewFetchTask(start, end)

	//FETCH business managers
	bms, err := fetchBusinessMenagers(task, session)
	if err != nil {
		return task, err
	}
	accounts, err := fetchAccountsByBm(task, session, bms)
	if err != nil {
		return task, err
	}
	spend, err := fetchAccountSpend(session, accounts, start, end)
	if err != nil {
		return task, err
	}
	for _, s := range spend {

		var accBase db.DbAccountSpend
		accInfo := s.account

		accBase.AccountID = accInfo.AccountId
		accBase.ProviderType = db.Facebook
		accBase.Status = accInfo.Status
		accBase.Spend = s.spend
		accBase.AccountName = accInfo.Name
		accBase.BusinessID = accInfo.Business.Id
		accBase.BusinessName = accInfo.Business.Name
		accBase.DateRef = s.dateRef
		accBase.NumberOfCampaigns = uint16(len(s.campaigns))
		accBase.UpdatedAt = time.Now().UTC()
		//TODO: add all the necessary fields
		task.Accounts = append(task.Accounts, accBase)
		for _, v := range s.campaigns {
			campBase := db.DbCampaignSpend{
				AccountID:    accBase.AccountID,
				ProviderType: db.Facebook,
				AccountName:  accInfo.Name,
				BusinessID:   accInfo.Business.Id,
				BusinessName: accInfo.Business.Name,
				CampaignID:   v.campaignID,
				CampaignName: v.campaignName,
				Spend:        v.spend,
				DateRef:      s.dateRef,
				UpdatedAt:    time.Now().UTC(),
				Status:       db.UnknownStatus.String(),
			}
			task.Campaigns = append(task.Campaigns, campBase)
		}

	}

	return task, nil
}
