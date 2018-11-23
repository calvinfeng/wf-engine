package workflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"wf-engine/fleet"
	"wf-engine/global"

	"github.com/spf13/viper"
)

func requestIDLERobot(name string) *fleet.Robot {
	resp := make(chan *fleet.Robot)
	global.State.GetRobotByStatus <- global.RobotReqquest{
		Robot:    name,
		Status:   "IDLE",
		Response: resp,
	}

	return <-resp
}

func waitForIDLERobot(name string) *fleet.Robot {
	var robot *fleet.Robot
	for robot == nil {
		resp := make(chan *fleet.Robot)
		global.State.GetRobotByStatus <- global.RobotReqquest{
			Robot:    name,
			Status:   "IDLE",
			Response: resp,
		}

		robot = <-resp
		time.Sleep(time.Second)
	}

	return robot
}

func httpSendRobotToNewPose(name string, pose fleet.Pose) error {
	data, err := json.Marshal(pose)
	if err != nil {
		return err
	}

	endpointf := "http://localhost:%d/api/robots/%s/send/"
	url := fmt.Sprintf(endpointf, viper.GetInt("http.port"), name)
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode >= 300 {
		b := bytes.NewBuffer([]byte{})
		b.ReadFrom(res.Body)
		return fmt.Errorf("encountered bad HTTP status code %d - %s", res.StatusCode, b.String())
	}

	return nil
}
