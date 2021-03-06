package commandevaluators

import (
	"errors"
	"log"
	"strings"

	"github.com/byuoitav/av-api/base"
	"github.com/byuoitav/av-api/dbo"
	"github.com/byuoitav/configuration-database-microservice/accessors"
)

// PowerOn is struct that implements the CommandEvaluation struct
type PowerOnDefault struct {
}

// Evaluate fulfills the CommmandEvaluation evaluate requirement.
func (p *PowerOnDefault) Evaluate(room base.PublicRoom) (actions []base.ActionStructure, err error) {
	log.Printf("Evaluating for PowerOn command.")
	var devices []accessors.Device
	if strings.EqualFold(room.Power, "on") {
		log.Printf("Room-wide PowerOn request received. Retrieving all devices.")
		//Get all devices.
		devices, err = dbo.GetDevicesByRoom(room.Building, room.Room)
		if err != nil {
			return
		}
		log.Printf("Setting power 'on' state for all output devices.")
		// Currently we only check for output devices.
		for i := range devices {
			if devices[i].Output {
				log.Printf("Adding device %+v", devices[i].Name)

				actions = append(actions, base.ActionStructure{
					Action:              "PowerOn",
					Device:              devices[i],
					GeneratingEvaluator: "PowerOnDefault",
					DeviceSpecific:      false,
				})
			}
		}
	}

	// Now we go through and check if power 'on' was set for any other device.
	log.Printf("Evaluating displays for power on command.")
	for _, device := range room.Displays {
		actions, err = p.evaluateDevice(device.Device, actions, devices, room.Room, room.Building)
		if err != nil {
			return
		}
	}

	for _, device := range room.AudioDevices {
		log.Printf("Evaluating audio devices for command power on. ")
		actions, err = p.evaluateDevice(device.Device, actions, devices, room.Room, room.Building)
		if err != nil {
			return
		}
	}

	log.Printf("%v actions generated.", len(actions))
	log.Printf("Evaluation complete.")

	return
}

// Validate fulfills the Fulfill requirement on the command interface
func (p *PowerOnDefault) Validate(action base.ActionStructure) (err error) {

	log.Printf("Validating action for comand PowerOn")

	ok, _ := CheckCommands(action.Device.Commands, "PowerOn")
	if !ok || !strings.EqualFold(action.Action, "PowerOn") {
		log.Printf("ERROR. %s is an invalid command for %s", action.Action, action.Device.Name)
		return errors.New(action.Action + " is an invalid command for" + action.Device.Name)
	}

	log.Printf("Done.")
	return
}

// GetIncompatibleCommands keeps track of actions that are incompatable (on the same device)
func (p *PowerOnDefault) GetIncompatibleCommands() (incompatableActions []string) {
	incompatableActions = []string{
		"standby",
	}

	return
}

// Evaluate devices just pulls out the process we do with the audio-devices and displays into one function.
func (p *PowerOnDefault) evaluateDevice(device base.Device,
	actions []base.ActionStructure,
	devices []accessors.Device,
	room string,
	building string) ([]base.ActionStructure, error) {

	// Check if we even need to start anything
	if strings.EqualFold(device.Power, "on") {
		// check if we already added it
		index := checkActionListForDevice(actions, device.Name, room, building)
		if index == -1 {
			// Get the device, check the list of already retreived devices first, if not there,
			// hit the DB up for it.
			dev, err := getDevice(devices, device.Name, room, building)
			if err != nil {
				return actions, err
			}

			actions = append(actions, base.ActionStructure{
				Action:              "PowerOn",
				Device:              dev,
				GeneratingEvaluator: "PowerOnDefault",
				DeviceSpecific:      true,
			})
		}
	}
	return actions, nil
}
