package main

import (
	"log"

	"strconv"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// Visu is the main function for the visualisation
func convertBoard(board [][]int) [][]string {
	var convertedBoard [][]string
	for i := 0; i < len(board); i++ {
		var row []string
		for j := 0; j < len(board); j++ {
			if board[i][j] == 0 {
				row = append(row, " ")
				continue
			}
			row = append(row, strconv.Itoa(board[i][j]))
		}
		convertedBoard = append(convertedBoard, row)
	}
	return convertedBoard

}

func PrintBoard(Board [][]int) bool {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	table := widgets.NewTable()
	table.Title = "n-puzzle"
	table.TitleStyle = ui.NewStyle(ui.ColorBlue)
	table.Rows = convertBoard(Board)
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.RowSeparator = true
	table.BorderStyle = ui.NewStyle(ui.ColorGreen)
	table.SetRect(0, 0, len(Board)*8, len(Board)*4)
	table.FillRow = true
	table.TextAlignment = ui.AlignCenter

	ui.Render(table)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return false
		case "s":
			moveUp(Board)
		case "w":
			moveDown(Board)
		case "d":
			moveLeft(Board)

		case "a":
			moveRight(Board)
		}
		if isEqual(Board, goal(len(Board))) {
			ui.Clear()
			p := widgets.NewParagraph()
			p.Text = "You won ! do you want to restart ? (y/n)"
			p.SetRect(0, 0, 25, 5)
			p.TextStyle = ui.NewStyle(ui.ColorGreen)
			p.BorderStyle = ui.NewStyle(ui.ColorGreen)
			ui.Render(p)
			for {
				e := <-uiEvents
				switch e.ID {
				case "n", "<C-c>":
					return false
				case "y":
					return true
				}
			}
		} else {
			table.Rows = convertBoard(Board)
			ui.Render(table)
		}

	}
}
