package FatSecretWrapper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type FatSecretWrapper interface {
	GetFoodIdFromBarcode(ctx context.Context, barcode string) (FatSecretFoodId, error)
	GetFoodFromId(ctx context.Context, id int) (FatSecretFood, error)
	SearchFoodsByName(ctx context.Context, searchQuery string, page *string) (FatSecretFoodsSearch, error)
}

type fatSecretWrapper struct{}

func NewFatSecretWrapper() FatSecretWrapper {
	return &fatSecretWrapper{}
}

type FatSecretFood struct {
	Food Food `json:"food"`
}

type FatSecretFoodsArray struct {
	Food []Food `json:"food"`
}

type Food struct {
	FoodId    string           `json:"food_id"`
	FoodName  string           `json:"food_name"`
	FoodType  string           `json:"food_type"`
	Servings  FatSecretServing `json:"servings"`
	BrandName string           `json:"brand_name"`
}

type FatSecretServing struct {
	Serving []Serving `json:"serving"`
}

type Serving struct {
	ServingId              string `json:"serving_id"`
	ServingDescription     string `json:"serving_description"`
	MetricServingAmount    string `json:"metric_serving_amount"`
	MetricServingUnit      string `json:"metric_serving_unit"`
	NumberOfUnits          string `json:"number_of_units"`
	MeasurementDescription string `json:"measurement_description"`
	Calories               string `json:"calories"`
	Protein                string `json:"protein"`
	Fat                    string `json:"fat"`
}

type FatSecretFoodId struct {
	FoodId FoodId `json:"food_id"`
}

type FoodId struct {
	Value string `json:"value"`
}

type FatSecretFoodsSearch struct {
	FoodsSearch FoodsSearch `json:"foods_search"`
}

type FoodsSearch struct {
	MaxResults   string              `json:"max_results"`
	TotalResults string              `json:"total_results"`
	PageNumber   string              `json:"page_number"`
	Results      FatSecretFoodsArray `json:"results"`
}

type FatSecretTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

var (
	FatSecretToken            string
	FatSecretTokenExpiresTime time.Time
)

func (fatSecretWrapper *fatSecretWrapper) GetFoodIdFromBarcode(ctx context.Context, barcode string) (FatSecretFoodId, error) {
	var foodId FatSecretFoodId
	responseData, err := fatSecretWrapper.apiRequestWithPayload(ctx, fmt.Sprintf("food/barcode/find-by-id/v1?barcode=%s&format=json", barcode), http.MethodGet, nil)
	if err != nil {
		return foodId, err
	}
	food_id_unmarshal_error := json.Unmarshal(responseData, &foodId)
	if food_id_unmarshal_error != nil {
		return foodId, food_id_unmarshal_error
	}
	return foodId, nil
}

func (fatSecretWrapper *fatSecretWrapper) GetFoodFromId(ctx context.Context, id int) (FatSecretFood, error) {
	var food FatSecretFood
	responseData, err := fatSecretWrapper.apiRequestWithPayload(ctx, fmt.Sprintf("food/v4?food_id=%v&format=json", id), http.MethodGet, nil)
	if err != nil {
		return food, err
	}

	food_unmarshal_error := json.Unmarshal(responseData, &food)
	if food_unmarshal_error != nil {
		return food, food_unmarshal_error
	}
	return food, nil
}

func (fatSecretWrapper *fatSecretWrapper) SearchFoodsByName(ctx context.Context, searchQuery string, page *string) (FatSecretFoodsSearch, error) {
	food := FatSecretFoodsSearch{}
	var pageParams = ""
	if page != nil {
		pageParams = fmt.Sprintf("?page_number=%s", *page)
	}
	encodedQuery := url.QueryEscape(searchQuery)
	queryParams := fmt.Sprintf("foods/search/v3?search_expression=%s%s&format=json", encodedQuery, pageParams)
	responseData, err := fatSecretWrapper.apiRequestWithPayload(ctx, queryParams, http.MethodGet, nil)
	if err != nil {
		return food, err
	}
	food_unmarshal_error := json.Unmarshal(responseData, &food)
	if food_unmarshal_error != nil {
		return food, food_unmarshal_error
	}
	return food, nil
}

func (fatSecretWrapper *fatSecretWrapper) apiRequestWithPayload(ctx context.Context, path string, verb string, body io.Reader) ([]byte, error) {
	if os.Getenv("FAT_SECRET_BASE_URL") == "" {
		return nil, errors.New("missing fat secret base url")
	}
	url := os.Getenv("FAT_SECRET_BASE_URL") + path
	client := &http.Client{}
	req, err := http.NewRequest(verb, url, body)
	if err != nil {
		return nil, errors.New("error creating request: " + err.Error())
	}

	auth_token, err := fatSecretWrapper.GetToken(ctx)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", auth_token))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(req)
	if err != nil {
		return nil, errors.New("error sending request: " + err.Error())
	}

	if response.Status[0:1] != "2" {
		return nil, errors.New("Unsuccessful status code returned: " + response.Status)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("error sending request: " + err.Error())
	}
	return responseData, nil
}

func (fatSecretWrapper *fatSecretWrapper) GetToken(ctx context.Context) (string, error) {
	if FatSecretToken != "" && FatSecretTokenExpiresTime.After(time.Now()) {
		return FatSecretToken, nil
	}

	tokenUrl := os.Getenv("FAT_SECRET_TOKEN_URL")
	if tokenUrl == "" {
		return "", errors.New("missing fat secret base url")
	}
	client := &http.Client{}
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	req, err := http.NewRequest(http.MethodPost, tokenUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return "", errors.New("error creating request: " + err.Error())
	}
	clientId := os.Getenv("FAT_SECRET_CLIENT_ID")
	if tokenUrl == "" {
		return "", errors.New("missing fat secret client id")
	}
	clientSecret := os.Getenv("FAT_SECRET_CLIENT_SECRET")
	if tokenUrl == "" {
		return "", errors.New("missing fat secret client secret")
	}
	req.SetBasicAuth(clientId, clientSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	responseData, err := io.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		return "", errors.New("non 200 status returned when trying to get token")
	}

	if err != nil {
		return "", err
	}
	var result FatSecretTokenResponse
	if err := json.Unmarshal(responseData, &result); err != nil {
		return "", err
	}
	FatSecretToken = result.AccessToken
	FatSecretTokenExpiresTime = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)

	return result.AccessToken, nil
}
