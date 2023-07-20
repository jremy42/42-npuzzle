package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// Visu is the main function for the visualisation
func convertBoard(board [][]uint8) [][]string {
	var convertedBoard [][]string
	for i := 0; i < len(board); i++ {
		var row []string
		for j := 0; j < len(board); j++ {
			if board[i][j] == 0 {
				row = append(row, " ")
				continue
			}
			row = append(row, strconv.Itoa(int(board[i][j])))
		}
		convertedBoard = append(convertedBoard, row)
	}
	return convertedBoard

}

func displayBoard(board [][]uint8, path []byte, elvalName string, times string, tries, openSetComplexity, workers, closeSetSplit, speedDisplay int) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	table := createTable(board)
	rec := table.GetRect()
	texte := ""
	//lenText := 0
	if len(path) == 0 {
		texte = fmt.Sprintln("this board is not solvable ")
		//lenText = len(texte) * 2
	} else {
		texte = fmt.Sprintf(
			`Success with : %v in %v !
	len of solution : %v
	time complexity / tries : %d
	space complexity : %d
	worker : %d
	close set split : %d
	`, elvalName, times, len(path), tries, openSetComplexity, workers, closeSetSplit)
		//lenText = 60
	}
	par := widgets.NewParagraph()
	par.Text = texte
	par.SetRect(0, rec.Max.Y, rec.Max.Y+30, rec.Max.Y+10)
	ui.Render(par)
	ui.Render(table)
	uiEvents := ui.PollEvents()

	for i := 0; i < len(path); i++ {
		select {
		case <-uiEvents:
			return
		default:
			switch path[i] {
			case 'U':
				_, board = moveUp(board)
			case 'D':
				_, board = moveDown(board)
			case 'L':
				_, board = moveLeft(board)
			case 'R':
				_, board = moveRight(board)
			}
			//fmt.Println(string(path[i]))
			time.Sleep(time.Duration(speedDisplay) * time.Millisecond)
			table.Rows = convertBoard(board)
			ui.Render(table)

		}
		//time.Sleep(5000 * time.Millisecond)
	}
	<-uiEvents
}

func playBoard(board [][]uint8) bool {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	table := createTable(board)

	ui.Render(table)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return false
		case "s":
			moveUp(board)
		case "w":
			moveDown(board)
		case "d":
			moveLeft(board)
		case "a":
			moveRight(board)
		}
		if isEqual(board, goal(len(board))) {
			return handleWinScenario()
		}
		table.Rows = convertBoard(board)
		ui.Render(table)
	}
}

func createTable(board [][]uint8) *widgets.Table {
	table := widgets.NewTable()
	table.Title = "n-puzzle"
	table.TitleStyle = ui.NewStyle(ui.ColorBlue)
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.RowSeparator = true
	table.BorderStyle = ui.NewStyle(ui.ColorGreen)
	table.SetRect(0, 0, len(board)*6, len(board)*2+1)
	table.FillRow = true
	table.TextAlignment = ui.AlignCenter
	table.Rows = convertBoard(board)
	return table
}

func handleWinScenario() bool {
	ui.Clear()
	p := createWinParagraph()
	ui.Render(p)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "n", "<C-c>":
			return false
		case "y":
			return true
		}
	}
}

func createPressAnyKeyParagraph() (p *widgets.Paragraph) {
	p = widgets.NewParagraph()
	p.Text = "Pres any key to exit"
	return
}

func createWinParagraph() *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.Text = "You won! Do you want to restart? (y/n)"
	p.SetRect(0, 0, 25, 5)
	p.TextStyle = ui.NewStyle(ui.ColorGreen)
	p.BorderStyle = ui.NewStyle(ui.ColorGreen)
	return p
}
