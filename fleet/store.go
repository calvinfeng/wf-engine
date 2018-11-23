package fleet

import (
	"fmt"
	"sync"
)

var store *Store

func init() {
	store = &Store{
		robots: make(map[string]*Robot),
		mutex:  &sync.Mutex{},
	}

	for i := 1; i <= 3; i++ {
		robot := &Robot{
			Name:        fmt.Sprintf("freight%d", i),
			Status:      "IDLE",
			CurrentPose: Pose{0, 0},
		}
		store.robots[robot.Name] = robot
	}
}

// Store keeps a list of resources on this mock server.
type Store struct {
	robots map[string]*Robot
	mutex  *sync.Mutex
}

// GetRobot checks whether a robot exists in store.
func (s Store) GetRobot(name string) *Robot {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if r, ok := s.robots[name]; ok {
		copy := *r
		return &copy
	}

	return nil
}

// GetRobots returns a list of robot in store.
func (s *Store) GetRobots() []*Robot {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	results := make([]*Robot, 0, len(s.robots))
	for _, r := range s.robots {
		copy := *r
		results = append(results, &copy)
	}

	return results
}

// UpdateRobot modifies x-y coordinate of a robot in store.
func (s *Store) UpdateRobot(name, status string, pose Pose) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.robots[name]; ok {
		s.robots[name].Status = status
		s.robots[name].CurrentPose = pose
	}
}
