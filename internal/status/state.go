package status

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type ProbeData struct {
	DeviceID string
	MacAddress, Rssi string
	PrevDetected     int64 //in milliseconds
	Timestamp        time.Time
}

type Queue interface {
	PublishToQueue(queueName string, payload []byte)
}

type RoomState struct {
	PotArrival, Arrived, PotDeparture map[string]ProbeData // Pot means Potential
	InUse                             *sync.Mutex
	mq Queue
}

var State RoomState

func InitializeRoomState(queue Queue) {
	State = RoomState{
		make(map[string]ProbeData),
		make(map[string]ProbeData),
		make(map[string]ProbeData),
		new(sync.Mutex),
		queue,
	}
}

func (state RoomState) InitializeCleanup() {
	state.CleanPotDeparture()
	state.CleanPotArrival()
	state.CleanArrived()
}

func ManageNewProbe(probe ProbeData) {
	State.InUse.Lock()
	defer State.InUse.Unlock()

	if State.CheckPotArrival(probe) {
		delete(State.PotArrival, probe.MacAddress)
		State.Arrived[probe.MacAddress] = probe
		State.mq.PublishToQueue(os.Getenv("queue"), probe.ProbeDataToMsg(true))

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
			state.PotDeparture[val.MacAddress] = val
			delete(state.Arrived, val.MacAddress)
		}
	}
	state.InUse.Unlock()
}

func (state RoomState) CleanPotDeparture() {
	state.InUse.Lock()
	for _, val := range state.PotDeparture {
		if time.Since(val.Timestamp) > time.Minute*3 {
			delete(state.PotDeparture, val.MacAddress)
			State.mq.PublishToQueue(os.Getenv("queue"), val.ProbeDataToMsg(false))
		}
	}
	state.InUse.Unlock()
}

type MSG struct {
	DeviceID string
	MacAddress string
	Active     bool //in milliseconds
	Timestamp        time.Time
}

func (data ProbeData) ProbeDataToMsg(active bool) []byte {
	doc, err := json.Marshal(MSG{
		data.DeviceID,
		data.MacAddress,
		active,
		data.Timestamp,
	})
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func (state RoomState) PublishMsg(queue string, p []byte) {
	state.mq.PublishToQueue(queue, p)
}



