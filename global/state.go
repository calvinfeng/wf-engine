package global

import (
	"context"
	"time"
	"wf-engine/fleet"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// State is the global state.
var State *state

func init() {
	State = &state{
		GetRobotByStatus: make(chan RobotReqquest),
		update:           make(chan stateUpdate),
		robots:           make(map[string]*fleet.Robot),
	}
}

type state struct {
	GetRobotByStatus chan RobotReqquest
	update           chan stateUpdate
	robots           map[string]*fleet.Robot
}

func (s *state) Activate(ctx context.Context, updateDone chan struct{}) {
	go s.pollRobots(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case req := <-s.GetRobotByStatus:
			s.handleRobotRequest(req)
		case update := <-s.update:
			s.handleUpdate(update)
			select {
			case updateDone <- struct{}{}:
			default:
			}
		}
	}
}

func (s *state) pollRobots(ctx context.Context) {
	ticker := time.NewTicker(viper.GetDuration("global.polling_intv"))
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			robots, err := httpFetchRobotList()
			if err != nil {
				log.Error(err)
				continue
			}

			done := make(chan struct{})
			s.update <- stateUpdate{robots: robots, done: done}
			<-done
		}
	}
}

func (s *state) handleRobotRequest(req RobotReqquest) {
	if _, ok := s.robots[req.Robot]; !ok {
		req.Response <- nil
		return
	}

	if s.robots[req.Robot].Status != req.Status {
		req.Response <- nil
		return
	}

	copy := *s.robots[req.Robot]
	req.Response <- &copy
}

func (s *state) handleUpdate(update stateUpdate) {
	for _, robot := range update.robots {
		s.robots[robot.Name] = robot
	}

	update.done <- struct{}{}
}
