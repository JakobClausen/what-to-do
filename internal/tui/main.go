package tui

import (
	list "what-to-do/internal/tui/list"
)

func InitialModel(width int) list.ListModel {

	// return column.CreateColumnModel(width)
	return list.NewListModel()
}
