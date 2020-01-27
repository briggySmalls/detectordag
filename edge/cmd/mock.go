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
	edge "github.com/briggysmalls/detectordag/edge/internal"
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mockCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mockCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type dashboard struct {
	powerState bool
	messenger  edge.Messenger
	paragraph  *widgets.Paragraph
	list       *widgets.List
}

func run(cmd *cobra.Command, args []string) {
	// Create messenger
	messenger := edge.NewMessenger()
	if err := messenger.Connect("amqp://guest:guest@localhost:5672/"); err != nil {
		log.Fatalf("Failed to connect to AMQP: %v", err)
	}

	// Initialise terminal
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// Create a prompt
	paragraph := widgets.NewParagraph()
	paragraph.Title = "Mock dag-edge"

	// Create a list of controls
	list := widgets.NewList()
	list.Title = "Key commands"
	list.Rows = []string{
		"[p] power",
		"[q] quit",
	}
	list.TextStyle = ui.NewStyle(ui.ColorYellow)
	list.WrapText = false

	// Create a grid layout
	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	// Add the widgets to the grid
	grid.Set(
		ui.NewRow(
			1,
			ui.NewCol(0.8, paragraph),
			ui.NewCol(0.2, list)))

	// Draw the UI
	ui.Render(grid)

	// Listen for keyboard events
	d := dashboard{
		messenger: messenger,
		list:      list,
		paragraph: paragraph,
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
	if err := d.messenger.PowerStatusChanged(d.powerState); err != nil {
		d.paragraph.Text += fmt.Sprintf("Error sending power status: %v\n", err)
	} else {
		d.paragraph.Text += fmt.Sprintf("Power status message sent: %v\n", d.powerState)
	}
	// Update the ui
	ui.Render(d.paragraph)
}
