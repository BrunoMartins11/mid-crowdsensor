package status_test

import (
	"github.com/BrunoMartins11/mid-crowdsensor/internal/status"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCheckPotDeparture(t *testing.T) {
	status.InitializeRoomState()
	data := status.ProbeData{
		MacAddress:   "123",
		Rssi:         "23",
		PrevDetected: 12,
		Timestamp:    time.Now(),
	}
	assert.Equal(t, false, status.State.CheckPotDeparture(data))

	status.State.PotDeparture["123"] = data
	assert.Equal(t, true, status.State.CheckPotDeparture(data))
}

func TestCheckPotArrival(t *testing.T) {
	status.InitializeRoomState()
	data := status.ProbeData{
		MacAddress:   "123",
		Rssi:         "23",
		PrevDetected: 12,
		Timestamp:    time.Now(),
	}
	assert.Equal(t, false, status.State.CheckPotArrival(data))

	status.State.PotArrival["123"] = data
	assert.Equal(t, true, status.State.CheckPotArrival(data))
}

func TestCheckArrived(t *testing.T) {
	status.InitializeRoomState()
	data := status.ProbeData{
		MacAddress:   "123",
		Rssi:         "23",
		PrevDetected: 12,
		Timestamp:    time.Now(),
	}
	assert.Equal(t, false, status.State.CheckArrived(data))

	status.State.Arrived["123"] = data
	assert.Equal(t, true, status.State.CheckArrived(data))
}

func TestManageNewProbeArrival(t *testing.T){
	status.InitializeRoomState()
	data := status.ProbeData{
		MacAddress:   "123",
		Rssi:         "23",
		PrevDetected: 12,
		Timestamp:    time.Now(),
	}

	status.ManageNewProbe(data)
	assert.Equal(t, data, status.State.PotArrival["123"])

	status.ManageNewProbe(data)
	assert.Equal(t, data, status.State.Arrived["123"])
	data.Rssi = "44"
	status.ManageNewProbe(data)
	assert.Equal(t, data, status.State.Arrived["123"])

	_, ok := status.State.PotArrival["123"]
	assert.Equal(t, false, ok)
}

func TestManageNewProbeDeparture(t *testing.T){
	status.InitializeRoomState()
	data := status.ProbeData{
		MacAddress:   "123",
		Rssi:         "23",
		PrevDetected: 12,
		Timestamp:    time.Now(),
	}
	status.State.PotDeparture[data.MacAddress] = data

	_, ok := status.State.PotArrival["123"]
	assert.Equal(t, false, ok)

	status.ManageNewProbe(data)
	assert.Equal(t, data, status.State.Arrived["123"])
}

func TestCleanUpPotArrival(t *testing.T){
	status.InitializeRoomState()
	loc := time.FixedZone("UTC-8", -8*60*60)
	data := status.ProbeData{
		MacAddress:   "123",
		Rssi:         "23",
		PrevDetected: 12,
		Timestamp:    time.Date(1998, time.August, 1, 1,1,1,1, loc),
	}
	data1 := status.ProbeData{
		MacAddress:   "1234",
		Rssi:         "23",
		PrevDetected: 12,
		Timestamp:    time.Now(),
	}
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
	loc := time.FixedZone("UTC-8", -8*60*60)
	data := status.ProbeData{
		MacAddress:   "123",
		Rssi:         "23",
		PrevDetected: 12,
		Timestamp:    time.Date(1998, time.August, 1, 1,1,1,1, loc),
	}
	data1 := status.ProbeData{
		MacAddress:   "1234",
		Rssi:         "23",
		PrevDetected: 12,
		Timestamp:    time.Now(),
	}
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
	loc := time.FixedZone("UTC-8", -8*60*60)
	data := status.ProbeData{
		MacAddress:   "123",
		Rssi:         "23",
		PrevDetected: 12,
		Timestamp:    time.Date(1998, time.August, 1, 1,1,1,1, loc),
	}
	data1 := status.ProbeData{
		MacAddress:   "1234",
		Rssi:         "23",
		PrevDetected: 12,
		Timestamp:    time.Now(),
	}
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
	loc := time.FixedZone("UTC-8", -8*60*60)
	data := status.ProbeData{
		MacAddress:   "123",
		Rssi:         "23",
		PrevDetected: 12,
		Timestamp:    time.Date(1998, time.August, 1, 1,1,1,1, loc),
	}
	data1 := status.ProbeData{
		MacAddress:   "1234",
		Rssi:         "23",
		PrevDetected: 12,
		Timestamp:    time.Now(),
	}
	status.State.PotArrival[data.MacAddress] = data
	status.State.Arrived[data.MacAddress] = data
	status.State.PotArrival[data.MacAddress] = data

	status.State.PotArrival[data1.MacAddress] = data1
	status.State.Arrived[data1.MacAddress] = data1
	status.State.PotDeparture[data1.MacAddress] = data1

	status.State.InitializeCleanup()
	_, ok := status.State.PotArrival["123"]
	assert.Equal(t, false, ok)
	_, ok = status.State.PotArrival["1234"]
	assert.Equal(t, true, ok)

	_, ok = status.State.Arrived["123"]
	assert.Equal(t, false, ok)
	_, ok = status.State.Arrived["1234"]
	assert.Equal(t, true, ok)

	_, ok = status.State.PotDeparture["123"]
	assert.Equal(t, false, ok)
	_, ok = status.State.PotDeparture["1234"]
	assert.Equal(t, true, ok)
}


