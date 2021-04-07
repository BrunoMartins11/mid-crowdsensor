package status

import (
	"sync"
	"time"
)

type ProbeData struct {
	MacAddress, Rssi string
	PrevDetected     int64 //in milliseconds
	Timestamp        time.Time
}

type RoomState struct {
	PotArrival, Arrived, PotDeparture map[string]ProbeData // Pot means Potential
	InUse                             *sync.Mutex
}

var State RoomState

func InitializeRoomState() {
	State = RoomState{
		make(map[string]ProbeData),
		make(map[string]ProbeData),
		make(map[string]ProbeData),
		new(sync.Mutex),
	}
}

func (state RoomState) InitializeCleanup() {
	state.CleanPotArrival()
	state.CleanArrived()
	state.CleanPotDeparture()
}

func ManageNewProbe(probe ProbeData) {
	State.InUse.Lock()
	defer State.InUse.Unlock()

	if State.CheckPotArrival(probe) {
		delete(State.PotArrival, probe.MacAddress)
		State.Arrived[probe.MacAddress] = probe

	} else if State.CheckArrived(probe) {
		State.Arrived[probe.MacAddress] = probe

	} else if State.CheckPotDeparture(probe) {
		delete(State.PotDeparture, probe.MacAddress)
		State.Arrived[probe.MacAddress] = probe

	} else {
		State.PotArrival[probe.MacAddress] = probe
	}
}

func (state RoomState) CheckPotDeparture(probe ProbeData) bool {
	_, exists := state.PotDeparture[probe.MacAddress]
	return exists
}

func (state RoomState) CheckPotArrival(probe ProbeData) bool {
	_, exists := state.PotArrival[probe.MacAddress]
	return exists
}

func (state RoomState) CheckArrived(probe ProbeData) bool {
	_, exists := state.Arrived[probe.MacAddress]
	return exists
}

func (state RoomState) CleanPotArrival() {
	state.InUse.Lock()
	for _, val := range state.PotArrival {
		if time.Since(val.Timestamp) > time.Minute*3 {
			delete(state.PotArrival, val.MacAddress)
		}
	}
	state.InUse.Unlock()
}

func (state RoomState) CleanArrived() {
	state.InUse.Lock()
	for _, val := range state.Arrived {
		if time.Since(val.Timestamp) > time.Minute*3 {
			delete(state.Arrived, val.MacAddress)
			state.PotDeparture[val.MacAddress] = val
		}
	}
	state.InUse.Unlock()
}

func (state RoomState) CleanPotDeparture() {
	state.InUse.Lock()
	for _, val := range state.PotDeparture {
		if time.Since(val.Timestamp) > time.Minute*3 {
			delete(state.PotDeparture, val.MacAddress)
		}
	}
	state.InUse.Unlock()
}
