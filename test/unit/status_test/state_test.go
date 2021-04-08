package status_test

import (
	"github.com/BrunoMartins11/mid-crowdsensor/internal/status"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)
var loc = time.FixedZone("UTC-8", -8*60*60)

var data = status.ProbeData{
MacAddress:   "123",
Rssi:         "23",
PrevDetected: 12,
Timestamp:    time.Date(1998, time.August, 1, 1,1,1,1, loc),
}
var data1 = status.ProbeData{
MacAddress:   "1234",
Rssi:         "23",
PrevDetected: 12,
Timestamp:    time.Now(),
}

var data2 = status.ProbeData{
MacAddress:   "12345",
Rssi:         "23",
PrevDetected: 12,
Timestamp:    time.Date(1998, time.August, 1, 1,1,1,1, loc),
}

func TestCheckPotDeparture(t *testing.T) {
	status.InitializeRoomState()
	assert.Equal(t, false, status.State.CheckPotDeparture(data1))

	status.State.PotDeparture["1234"] = data1
	assert.Equal(t, true, status.State.CheckPotDeparture(data1))
}

func TestCheckPotArrival(t *testing.T) {
	status.InitializeRoomState()
	assert.Equal(t, false, status.State.CheckPotArrival(data1))

	status.State.PotArrival["1234"] = data1
	assert.Equal(t, true, status.State.CheckPotArrival(data1))
}

func TestCheckArrived(t *testing.T) {
	status.InitializeRoomState()

	assert.Equal(t, false, status.State.CheckArrived(data1))

	status.State.Arrived["1234"] = data
	assert.Equal(t, true, status.State.CheckArrived(data1))
}

func TestManageNewProbeArrival(t *testing.T){
	status.InitializeRoomState()
	status.ManageNewProbe(data1)
	assert.Equal(t, data1, status.State.PotArrival[data1.MacAddress])

	status.ManageNewProbe(data1)
	assert.Equal(t, data1, status.State.Arrived[data1.MacAddress])
	data1.Rssi = "44"
	status.ManageNewProbe(data1)
	assert.Equal(t, data1, status.State.Arrived[data1.MacAddress])

	_, ok := status.State.PotArrival["1234"]
	assert.Equal(t, false, ok)
}

func TestManageNewProbeDeparture(t *testing.T){
	status.InitializeRoomState()
	status.State.PotDeparture[data1.MacAddress] = data1

	_, ok := status.State.PotArrival["1234"]
	assert.Equal(t, false, ok)

	status.ManageNewProbe(data1)
	assert.Equal(t, data1, status.State.Arrived["1234"])
}

func TestCleanUpPotArrival(t *testing.T){
	status.InitializeRoomState()
	status.State.PotArrival[data.MacAddress] = data
	status.State.PotArrival[data1.MacAddress] = data1
	status.State.CleanPotArrival()
	_, ok := status.State.PotArrival["123"]
	assert.Equal(t, false, ok)
	_, ok = status.State.PotArrival["1234"]
	assert.Equal(t, true, ok)
}

func TestCleanUpArrived(t *testing.T){
	status.InitializeRoomState()
	status.State.Arrived[data.MacAddress] = data
	status.State.Arrived[data1.MacAddress] = data1

	status.State.CleanArrived()
	_, ok := status.State.Arrived["123"]
	assert.Equal(t, false, ok)

	_, ok = status.State.Arrived["1234"]
	assert.Equal(t, true, ok)
}

func TestCleanUpPotDeparture(t *testing.T){
	status.InitializeRoomState()
	status.State.PotDeparture[data.MacAddress] = data
	status.State.PotDeparture[data1.MacAddress] = data1

	status.State.CleanPotDeparture()
	_, ok := status.State.PotDeparture["123"]
	assert.Equal(t, false, ok)

	_, ok = status.State.PotDeparture["1234"]
	assert.Equal(t, true, ok)
}

func TestInitializeCleanup(t *testing.T){
	status.InitializeRoomState()
	status.State.PotArrival[data.MacAddress] = data
	status.State.Arrived[data.MacAddress] = data
	status.State.PotDeparture[data2.MacAddress] = data2

	status.State.PotArrival[data1.MacAddress] = data1
	status.State.Arrived[data1.MacAddress] = data1
	status.State.PotDeparture[data1.MacAddress] = data1

	status.State.InitializeCleanup()

	_, ok := status.State.PotDeparture["123"]
	assert.Equal(t, true, ok)
	_, ok = status.State.PotDeparture["1234"]
	assert.Equal(t, true, ok)
	_, ok = status.State.PotDeparture["12345"]
	assert.Equal(t, false, ok)


	_, ok = status.State.PotArrival["123"]
	assert.Equal(t, false, ok)
	_, ok = status.State.PotArrival["1234"]
	assert.Equal(t, true, ok)

	_, ok = status.State.Arrived["123"]
	assert.Equal(t, false, ok)
	_, ok = status.State.Arrived["1234"]
	assert.Equal(t, true, ok)
}


