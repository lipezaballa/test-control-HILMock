package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	trace "github.com/rs/zerolog/log"
)

const START_SIMULATION = "start_simulation"
const FINISH_SIMULATION = "finish_simulation"

const VEHICLE_STATE_ID = 1

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, msgByte, err := hilMock.backConn.ReadMessage()
			if err != nil {
				trace.Error().Err(err).Msg("error receiving message in IDLE")
				cancel()
			} else {
				msg := string(msgByte)
				switch msg {
				case START_SIMULATION:

					errStarting := hilMock.backConn.WriteMessage(websocket.BinaryMessage, []byte(START_SIMULATION))
					if errStarting != nil {
						trace.Error().Err(errStarting).Msg("Error sending message of starting simultaion to backend")
						cancel()
						break
					}

					err := hilMock.startSimulationState(ctx, cancel)
					trace.Info().Msg("IDLE")

					if err != nil {
						cancel()
						return
					}
				}

			}
		}

	}
}

func (hilMock *HilMock) startSimulationState(ctx context.Context, cancel context.CancelFunc) error {
	errChan := make(chan error)
	done := make(chan struct{})
	stopChan := make(chan struct{})
	trace.Info().Msg("Simulation state")

	hilMock.readOrdersBackend(ctx, errChan, stopChan)
	hilMock.sendVehicleState(ctx, errChan, done)

	for {
		select {
		case err := <-errChan:
			cancel()
			return err
		case <-stopChan:
			close(done)
			return nil
		default:
		}
	}
}

func (hilMock *HilMock) readOrdersBackend(ctx context.Context, errChan chan<- error, stopChan chan<- struct{}) {
	go func() {
		for {
			select {
			case <-ctx.Done(): //FIXME: Make sure this part is run
				return
			default:
				_, msg, err := hilMock.backConn.ReadMessage()
				stringMsg := string(msg)
				if err != nil {
					trace.Error().Err(err).Msg("Error reading message from back-end")
					errChan <- err
					return
				}
				if stringMsg == FINISH_SIMULATION {
					trace.Info().Msg("Finish simulation")
					errStoping := hilMock.backConn.WriteMessage(websocket.BinaryMessage, []byte(FINISH_SIMULATION))
					if errStoping != nil {
						trace.Error().Err(errStoping).Msg("Error returning finish msg to back")
						errChan <- errStoping
					}
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
					trace.Warn().Msg("Does NOT match any type")
				}
			}

		}
	}()
}

func (hilMock *HilMock) sendVehicleState(ctx context.Context, errChan chan<- error, done <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-done:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				vehiclesState := []VehicleState{}
				vehicleState := RandomVehicleState()
				vehiclesState = append(vehiclesState, vehicleState) //FIXME: As an array?
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
