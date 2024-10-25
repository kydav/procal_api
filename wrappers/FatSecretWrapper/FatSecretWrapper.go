package FatSecretWrapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type FatSecretWrapper interface {
	GetFoodIdFromBarcode(barcode string) (int, error)
	GetFoodFromId(id int) (FatSecretFood, error)
	SearchFoodsByName(id int) ([]Food, error)
}

type fatSecretWrapper struct{}

func NewFatSecretWrapper() FatSecretWrapper {
	return &fatSecretWrapper{}
}

type FatSecretFood struct {
	Food Food `json:"food"`
}

type Food struct {
	FoodId   int64            `json:"food_id"`
	FoodName string           `json:"food_name"`
	FoodType string           `json:"food_type"`
	Servings FatSecretServing `json:"servings"`
}

type FatSecretServing struct {
	Serving []Serving `json:"serving"`
}

type Serving struct {
	ServingId              int64   `json:"serving_id"`
	ServingDescription     string  `json:"serving_description"`
	MetricServingAmount    float64 `json:"metric_serving_amount"`
	MetricServingUnit      string  `json:"metric_serving_unit"`
	NumberOfUnits          float64 `json:"number_of_units"`
	MeasurementDescription string  `json:"measurement_description"`
	Calories               float64 `json:"calories"`
	Protein                float64 `json:"protein"`
}

type FatSecretTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type FoodIdQuery struct {
	FoodId int `json:"food_id"`
}

func (fatSecretWrapper *fatSecretWrapper) GetFoodIdFromBarcode(barcode string) (int, error) {
	id := 000
	return id, nil
}

func (fatSecretWrapper *fatSecretWrapper) GetFoodFromId(id int) (FatSecretFood, error) {
	var food FatSecretFood
	responseData, _ := fatSecretWrapper.apiRequestWithPayload(fmt.Sprintf("food/v4?food_id=%v&format=json", id), http.MethodGet, nil)
	//responseData, _ := fatSecretWrapper.apiRequestWithPayload("food/barcode/find-by-id/v1?barcode=0041570054161&format=json", http.MethodGet, nil)

	food_unmarshal_error := json.Unmarshal(responseData, &food)
	if food_unmarshal_error != nil {
		return food, food_unmarshal_error
	}
	return food, nil
}

func (fatSecretWrapper *fatSecretWrapper) SearchFoodsByName(id int) ([]Food, error) {
	foods := []Food{}
	return foods, nil
}

func (fatSecretWrapper *fatSecretWrapper) apiRequestWithPayload(path string, verb string, body io.Reader) ([]byte, error) {
	if os.Getenv("FAT_SECRET_BASE_URL") == "" {
		return nil, errors.New("not configured to hit fat secret")
	}
	url := os.Getenv("FAT_SECRET_BASE_URL") + path
	client := &http.Client{}
	req, err := http.NewRequest(verb, url, body)
	if err != nil {
		log.Warn("error creating request")
		return nil, errors.New("error creating request: " + err.Error())
	}

	auth_token, err := fatSecretWrapper.GetToken()
	if err != nil {
		return []byte{}, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", auth_token))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//req.   Header.Add("Parameters", "method=foods.search&search_expression=toast&format=json")
	response, err := client.Do(req)
	if err != nil {
		log.Warn("error during http request")
		return nil, errors.New("error sending request: " + err.Error())
	}

	if response.Status[0:1] != "2" {
		responseData, err := io.ReadAll(response.Body)
		log.Warn(response)
		log.Warn(responseData)
		log.Warn(err)
		return nil, errors.New("Unsuccessful status code returned: " + response.Status)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Warn("unable to read body")
		return nil, errors.New("error sending request: " + err.Error())
	}
	return responseData, nil
}

func (fatSecretWrapper *fatSecretWrapper) GetToken() (string, error) {
	tokenUrl := os.Getenv("FAT_SECRET_TOKEN_URL")
	client := &http.Client{}
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	//data.Set("scope", "basic")
	req, err := http.NewRequest(http.MethodPost, tokenUrl, strings.NewReader(data.Encode()))
	if err != nil {
		log.Warn("error creating request")
		return "", errors.New("error creating request: " + err.Error())
	}
	clientId := os.Getenv("FAT_SECRET_CLIENT_ID")
	clientSecret := os.Getenv("FAT_SECRET_CLIENT_SECRET")
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

	return result.AccessToken, nil
}
