//go:build e2e

package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type GetRoomsSuite struct {
	suite.Suite
}

func (g *GetRoomsSuite) Test_OnNoErrors_GetRooms() {
	if _, ok := os.LookupEnv("API_URL"); !ok {
		g.Require().FailNow("environment variable API_URL not set")
	}

	apiUrl := os.Getenv("API_URL")
	response, err := http.Post(fmt.Sprintf("%s/api/sign-up", apiUrl), "application/json", bytes.NewBuffer([]byte(`
		{
			"name": "John Doe",
			"email": "john.doe@gmail.com",
			"password": "123456789"
		}
	`)))
	g.Require().NoError(err)
	g.Equal(201, response.StatusCode)

	response, err = http.Post(fmt.Sprintf("%s/api/login-with-email-and-password", apiUrl), "application/json", bytes.NewBuffer([]byte(`
		{
			"email": "john.doe@gmail.com",
			"password": "123456789"
		}
	`)))
	g.Require().NoError(err)
	g.Equal(200, response.StatusCode)
	defer response.Body.Close()

	var responseBody map[string]any
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	g.Require().NoError(err)

	accessToken := responseBody["data"].(map[string]any)["accessToken"].(string)

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/api/rooms", apiUrl), nil)
	g.Require().NoError(err)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", accessToken)

	client := &http.Client{}
	response, err = client.Do(request)
	g.Require().NoError(err)
	g.Equal(200, response.StatusCode)
	defer response.Body.Close()
	roomsResponseBody, err := io.ReadAll(response.Body)
	g.Require().NoError(err)

	g.JSONEq(`
		{
			"statusCode": 200,
			"statusText": "OK",
			"data": []
		}
	`, string(roomsResponseBody))
}

func TestGetRooms(t *testing.T) {
	suite.Run(t, new(GetRoomsSuite))
}
