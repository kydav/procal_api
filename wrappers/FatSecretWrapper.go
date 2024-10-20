package wrappers

import (
	"errors"
	"io"
	"net/http"
	"os"
)

type FatSecretWrapper interface {
	GetFoodIdFromBarcode(barcode string) (int, error)
}

type fatSecretWrapper struct{}

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

func NewFatSecretWrapper() FatSecretWrapper {
	return &fatSecretWrapper{}
}

func (fatSecretWrapper *fatSecretWrapper) GetFoodIdFromBarcode(barcode string) (int, error) {
	id := 000
	return id, nil
}

func (fatSecretWrapper *fatSecretWrapper) GetFoodFromId(id int) (Food, error) {
	food := Food{}
	return food, nil
}

func (fatSecretWrapper *fatSecretWrapper) SearchFoodsByName(id int) ([]Food, error) {
	foods := []Food{}
	return foods, nil
}

func (fatSecretWrapper *fatSecretWrapper) apiRequestWithPayload(path string, verb string, body io.Reader) ([]byte, error) {
	if os.Getenv("DOCUSIGN_API_BASE_URL") == "" {
		return nil, errors.New("not configured to hit Docusign")
	}
	url := os.Getenv("DOCUSIGN_API_BASE_URL") + "/v2.1/accounts/" + os.Getenv("DOCUSIGN_USER_API_ID") + "/" + path
	client := &http.Client{}
	req, err := http.NewRequest(verb, url, body)
	if err != nil {
		//log.Warn("error creating request")
		return nil, errors.New("error creating request: " + err.Error())
	}

	//auth_token, err := docusignWrapper.GetToken()
	if err != nil {
		return []byte{}, err
	}
	//req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", auth_token))
	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		//log.Warn("error during http request")
		return nil, errors.New("error sending request: " + err.Error())
	}
	//defer utils.CloseHandle(response.Body)

	if response.Status[0:1] != "2" {
		//responseData, err := io.ReadAll(response.Body)
		//log.Warn(responseData)
		//log.Warn(err)
		return nil, errors.New("Unsuccessful status code returned: " + response.Status)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		//log.Warn("unable to read body")
		return nil, errors.New("error sending request: " + err.Error())
	}
	return responseData, nil
}
