package FatSecretWrapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type FatSecretWrapper interface {
	GetFoodIdFromBarcode(barcode string) (int, error)
	GetFoodFromId(id int) (Food, error)
	SearchFoodsByName(id int) ([]Food, error)
}

type fatSecretWrapper struct{}

func NewFatSecretWrapper() FatSecretWrapper {
	return &fatSecretWrapper{}
}

type Food struct {
	FoodId   string    `json:"food_id"`
	FoodName string    `json:"food_name"`
	FoodType string    `json:"food_type"`
	Servings []Serving `json:"servings"`
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
}

type FatSecretTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (fatSecretWrapper *fatSecretWrapper) GetFoodIdFromBarcode(barcode string) (int, error) {
	id := 000
	return id, nil
}

func (fatSecretWrapper *fatSecretWrapper) GetFoodFromId(id int) (Food, error) {
	responseData, _ := fatSecretWrapper.apiRequestWithPayload("food/v4", http.MethodGet, nil)
	var food Food
	if food_unmarshal_error := json.Unmarshal(responseData, &food); food_unmarshal_error == nil {
		return food, food_unmarshal_error
	}
	return food, nil
}

func (fatSecretWrapper *fatSecretWrapper) SearchFoodsByName(id int) ([]Food, error) {
	foods := []Food{}
	return foods, nil
}

func (fatSecretWrapper *fatSecretWrapper) GetToken() (string, error) {
	tokenUrl := os.Getenv("FAT_SECRET_TOKEN_URL")
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, tokenUrl, nil)
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
	req.Header.Add("Content-Type", "application/json")
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
