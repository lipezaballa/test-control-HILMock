package main

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	trace "github.com/rs/zerolog/log"
)

const START_MSG = "start_simulation"
const FINISH_SIMULATION = "finish_simulation"

type HilMock struct {
	backConn *websocket.Conn
}

func NewHilMock() *HilMock {
	return &HilMock{}
}
func (hilMock *HilMock) SetBackConn(conn *websocket.Conn) {
	hilMock.backConn = conn
}

func (hilMock *HilMock) startIDLE() {
	trace.Info().Msg("IDLE")
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			_, msgByte, err := hilMock.backConn.ReadMessage()
			if err != nil {
				trace.Error().Err(err).Msg("error receiving message in IDLE")
			} else {
				msg := string(msgByte)
				switch msg {
				case START_MSG:

					errStarting := hilMock.backConn.WriteMessage(websocket.BinaryMessage, []byte(START_MSG))
					if errStarting != nil {
						trace.Error().Err(errStarting).Msg("Error sending message of starting simultaion to backend")
						break
					}
					fmt.Println(START_MSG)

					err := hilMock.startSimulationState()
					trace.Info().Msg("IDLE")

					if err != nil {
						return
					}
				}
			}
		}

	}
}

func (hilMock *HilMock) startSimulationState() error {
	errChan := make(chan error)
	done := make(chan struct{})
	//dataChan := make(chan VehicleState) FIXME: Is it necessary to store them?
	//orderChan := make(chan Order) FIXME: Is it necessary to store them?
	stopChan := make(chan struct{})
	trace.Info().Msg("Simulation state")

	hilMock.readOrdersBackend(done, errChan, stopChan)
	hilMock.sendVehicleState(done, errChan)

	for {
		select {
		case err := <-errChan:
			close(done)
			return err
		case <-stopChan:
			close(done)
			return nil
		default:
		}
	}
}

func (hilMock *HilMock) readOrdersBackend(done <-chan struct{}, errChan chan<- error, stopChan chan<- struct{}) {
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				_, msg, err := hilMock.backConn.ReadMessage()
				stringMsg := string(msg)
				if err != nil {
					trace.Error().Err(err).Msg("Error reading message from back-end")
					errChan <- err
					return //FIXME
				}
				if stringMsg == FINISH_SIMULATION {
					trace.Info().Msg("Finsih simulation")
					stopChan <- struct{}{}
					return
				}

				dataType := binary.LittleEndian.Uint16(msg[0:2])
				switch dataType {
				case 2:
					var order FormOrder
					order.Read(msg[2:])
					trace.Info().Msg(fmt.Sprintf("Form order: %v", order))

				case 3:
					var order ControlOrder
					order.Read(msg[2:])
					trace.Info().Msg(fmt.Sprintf("Control order: %v", order))
				default:
					fmt.Println("Does NOT match any type")
				}
			}

		}
	}()
}

func (hilMock *HilMock) sendVehicleState(done <-chan struct{}, errChan chan<- error) {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				vehiclesState := []VehicleState{}
				vehicleState := RandomVehicleState()
				vehiclesState = append(vehiclesState, vehicleState)
				trace.Info().Msg(fmt.Sprint(vehiclesState))
				head := make([]byte, 2)
				binary.LittleEndian.PutUint16(head, VEHICLE_STATE_ID)

				msg := GetAllBytesFromVehicleState(vehiclesState)

				encodedMsg := append(head, msg...)

				err := hilMock.backConn.WriteMessage(websocket.BinaryMessage, encodedMsg)
				if err != nil {
					trace.Error().Err(err).Msg("Error sending message from back-end")
					errChan <- err
					return
				}
				//FIXME: Add default?
			}
		}
	}()
}
