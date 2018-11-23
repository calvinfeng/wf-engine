package fleet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// Robot is a robot.
type Robot struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	CurrentPose Pose   `json:"current_pose"`
}

// Pose is like a coordinate.
type Pose struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func newRobotListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		robots := store.GetRobots()

		bytes, err := json.Marshal(robots)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	}
}

func newSendRobotHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		decoder := json.NewDecoder(r.Body)

		target := Pose{}
		if err := decoder.Decode(&target); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		vars := mux.Vars(r)
		robot := store.GetRobot(vars["robot"])
		if robot == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("robot %s does not exist", vars["robot"])))
			return
		}

		go navigate(vars["robot"], robot.CurrentPose, target)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte{})
	}
}

func navigate(robot string, current, target Pose) {
	dX := (target.X - current.X) / 20
	dY := (target.Y - current.Y) / 20
	for i := 1; i <= 20; i++ {
		time.Sleep(viper.GetDuration("robot.update_intv"))
		newPose := Pose{X: current.X + dX*float64(i), Y: current.Y + dY*float64(i)}
		store.UpdateRobot(robot, "WORKING", newPose)
	}

	store.UpdateRobot(robot, "IDLE", target)
}
