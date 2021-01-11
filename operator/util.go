package operator

import (
	"encoding/json"
	"net/http"

	"github.com/school/school-operator/models"
)

var ClassUrl = "http://192.168.0.80:8080/apis/school.crd.io/v1/classes"

func GetClassStatus() (error, *models.Classes) {

	request, err := http.NewRequest(http.MethodGet, ClassUrl, nil)
	if err != nil {
		return err, nil
	}

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()

	res := models.Classes{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return err, nil
	}
	return nil, &res
}
