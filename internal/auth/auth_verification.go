package auth

import (
	"log"
	"net/http"
	"os"
)

func IsValidToken(token string) bool {
	client := &http.Client{}
	//Submit request
	request, err := http.NewRequest("GET", os.Getenv("IOT_AUTH") + "/validate", nil)
	if err != nil {
		log.Fatalln(err)
	}
	//Add header option
	request.Header.Add("Authorization", "Bearer " +  token)

	//Processing returned results
	response, err := client.Do(request)

	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()

	if response.StatusCode == 401 || response.StatusCode == 400 {
		return false
	}

	return true
}
