package global

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"wf-engine/fleet"
)

func httpFetchRobotList() ([]*fleet.Robot, error) {
	url := "http://localhost:8000/api/robots/"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 300 {
		b := bytes.NewBuffer([]byte{})
		b.ReadFrom(res.Body)
		return nil, fmt.Errorf("encountered bad HTTP status code %d - %s", res.StatusCode, b.String())
	}

	robots := []*fleet.Robot{}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(bodyBytes, &robots)
	if err != nil {
		return nil, err
	}

	return robots, nil
}
