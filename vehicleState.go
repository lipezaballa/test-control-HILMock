package main

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/rand"
)

const VEHICLE_STATE_LENGTH = 25

type VehicleState struct {
	YDistance   float64 `json:"yDistance"` //It goes between 22mm and 10 mm
	Current     float64 `json:"current"`
	Duty        byte    `json:"duty"`
	Temperature float64 `json:"temperature"`
}

func RandomVehicleState() VehicleState {
	VehicleState := &VehicleState{}
	//yDistance := perlin.NewPerlin(2, 2, 3, 5) //TODO: Implement Perlin Noise
	VehicleState.YDistance = float64(rand.Intn(13)+10) + (math.Round(rand.Float64()*100) / 100)
	VehicleState.Current = float64(rand.Intn(20)) + (math.Round(rand.Float64()*100) / 100)
	VehicleState.Duty = byte(rand.Intn(100))
	VehicleState.Temperature = float64(rand.Intn(40)+20) + (math.Round(rand.Float64()*100) / 100)
	return *VehicleState
}

func GetAllVehicleStates(data []byte) ([]VehicleState, error) {
	vehicleStateArray := []VehicleState{}
	reader := bytes.NewReader(data)
	var err error
	for i := 0; i <= len(data)-VEHICLE_STATE_LENGTH; i += VEHICLE_STATE_LENGTH {
		vehicleState := &VehicleState{}
		err = binary.Read(reader, binary.LittleEndian, vehicleState)
		if err != nil {
			break
		}
		vehicleStateArray = append(vehicleStateArray, *vehicleState)
	}
	return vehicleStateArray, err
}

func ConvertFloat64ToBytes(num float64) [8]byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(num))
	return buf
}

func GetBytesFromVehicleState(vehicleState VehicleState) []byte {

	buf1 := ConvertFloat64ToBytes(vehicleState.YDistance)
	buf2 := ConvertFloat64ToBytes(vehicleState.Current)
	var buf3 [1]byte = [1]byte{vehicleState.Duty}
	buf4 := ConvertFloat64ToBytes(vehicleState.Temperature)

	return append(append(append(buf1[:], buf2[:]...), buf3[:]...), buf4[:]...)
}

func GetAllBytesFromVehicleState(vehiclesState []VehicleState) []byte {
	var result []byte
	for _, vehicle := range vehiclesState {
		result = append(result, GetBytesFromVehicleState(vehicle)...)
	}
	return result
}
