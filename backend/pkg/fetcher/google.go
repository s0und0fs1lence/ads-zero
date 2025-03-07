package fetcher

import (
	"fmt"
	"time"

	"github.com/revealbot/google-ads-go/ads"
	"github.com/revealbot/google-ads-go/services"
)

func GoogleFetcher(accessToken string, start, end time.Time) (totalSpend uint64, err error) {
	// Create a client from credentials file
	ads.NewClient(&ads.GoogleAdsClientParams{})
	client, err := ads.NewClientFromStorage("google-ads.json")
	if err != nil {
		panic(err)
	}

	// Load the "GoogleAds" service
	googleAdsService := services.NewGoogleAdsServiceClient(client.Conn())

	// Create a search request
	request := services.SearchGoogleAdsRequest{
		CustomerId: "2984242032",
		Query:      "SELECT campaign.id, campaign.name FROM campaign ORDER BY campaign.id",
	}

	// Get the results
	response, err := googleAdsService.Search(client.Context(), &request)
	for _, row := range response.Results {
		campaign := row.Campaign
		fmt.Printf("id: %d, name: %s\n", campaign.Id, *campaign.Name)
	}
	return totalSpend, fmt.Errorf("UNIMPLEMENTED")
}
