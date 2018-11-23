package global

import "wf-engine/fleet"

// RobotReqquest is request for robot from global state.
type RobotReqquest struct {
	Robot    string
	Status   string
	Response chan *fleet.Robot
}

type stateUpdate struct {
	robots []*fleet.Robot
	done   chan struct{}
}
