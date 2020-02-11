/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/briggysmalls/detectordag/shared"
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/spf13/cobra"
)

// mockCmd represents the mock command
var mockCmd = &cobra.Command{
	Use:   "mock",
	Short: "Mock version of the edge application",
	Run:   run,
}

func init() {
	rootCmd.AddCommand(mockCmd)
}

type dashboard struct {
	powerState bool
	messenger  shared.SensingMessenger
	messages   *widgets.List
	commands   *widgets.List
}

func run(cmd *cobra.Command, args []string) {
	// Create messenger
	messenger := shared.NewSensingMessenger()
	if err := messenger.Connect("amqp://guest:guest@localhost:5672/"); err != nil {
		log.Fatalf("Failed to connect to AMQP: %v", err)
	}
	defer messenger.Close()

	// Initialise terminal
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// Create a prompt
	messageList := widgets.NewList()
	messageList.Title = "Mock dag-edge"
	messageList.Rows = []string{}

	// Create a list of controls
	commandList := widgets.NewList()
	commandList.Title = "Key commands"
	commandList.Rows = []string{
		"[p] power",
		"[q] quit",
	}
	commandList.TextStyle = ui.NewStyle(ui.ColorYellow)
	commandList.WrapText = false

	// Create a grid layout
	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	// Add the widgets to the grid
	grid.Set(
		ui.NewRow(
			1,
			ui.NewCol(0.8, messageList),
			ui.NewCol(0.2, commandList)))

	// Draw the UI
	ui.Render(grid)

	// Listen for keyboard events
	d := dashboard{
		messenger: messenger,
		commands:  commandList,
		messages:  messageList,
	}
	for e := range ui.PollEvents() {
		switch e.ID {
		case "p": // Toggle power status
			d.togglePowerStatus()
		case "q": // Quit
			return
		}
	}
}

func (d *dashboard) togglePowerStatus() {
	// Toggle the state
	d.powerState = !d.powerState
	// Send a new message
	var message string
	if err := d.messenger.PowerStatusChanged(d.powerState); err != nil {
		message = fmt.Sprintf("Error sending power status: %v", err)
	} else {
		message = fmt.Sprintf("Power status message sent: %v", d.powerState)
	}
	d.messages.Rows = append(d.messages.Rows, message)
	d.messages.ScrollBottom()
	// Update the ui
	ui.Render(d.messages)
}
