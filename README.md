IP HIL: 127.0.0.2 #Por definir
IP backend: 127.0.0.1:8010 #Por definir
IP Frontend: 127.0.0.1 #Por definir

Message received and sent for starting and finishing the simulation: When receives it must sent it to communicate the backend that everything is correct. Both msg are sent as []byte. Front init this communication
const START_SIMULATION = "start_simulation"
const FINISH_SIMULATION = "finish_simulation"

This is the necessary communication received and sent by the HIL, everything in Little Endian:

Structs sent in []byte:

Prepared to send the identifier and several structs in the same msg

1. VehicleState: 2 first bytes identify the struct: VEHICLE_STATE_ID = 1
   type VehicleState struct {
   YDistance float64 `json:"yDistance"` //It goes between 22mm and 10 mm
   Current float64 `json:"current"`
   Duty byte `json:"duty"`
   Temperature float64 `json:"temperature"`
   }

Structs received in []byte:

Prepared to receive the identifier and several structs in the same msg

1. ControlOrder: 2 first bytes identify the struct: CONTROL_ORDER_ID = 3

type ControlOrder struct {
Id uint8 `json:"id"` #one byte
State bool `json:"state"` #only one byte
}

Ids:

-   0: Levitation
-   1: Propulsion
-   2: Brake
-   3: Custom 3
-   4: Custom 4
-   5: Custom 5
-   6: Custom 6

2. FormOrder: 2 first bytes identify the struct: FORM_ORDER_ID = 2
   NOT DEFINED BECAUSE AT THE MOMENT IT IS NOT GOING TO BE USED, IT IS FOR CUSTOM FORM WE CANCELED

type FormOrder struct {
Kind string `json:"kind"`
Payload float64 `json:"payload"`
}
