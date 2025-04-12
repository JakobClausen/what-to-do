	// // List.

	// list = lipgloss.NewStyle().
	// 	Border(lipgloss.NormalBorder(), false, true, false, false).
	// 	BorderForeground(subtle).
	// 	MarginRight(2).
	// 	Height(8).
	// 	Width(columnWidth + 1)

	// listHeader = base.
	// 	BorderStyle(lipgloss.NormalBorder()).
	// 	BorderBottom(true).
	// 	BorderForeground(subtle).
	// 	MarginRight(2).
	// 	Render

	// listItem = base.PaddingLeft(2).Render

	// checkMark = lipgloss.NewStyle().SetString("✓").
	// 	Foreground(special).
	// 	PaddingRight(1).
	// 	String()

	// listDone = func(s string) string {
	// 	return checkMark + lipgloss.NewStyle().
	// 		Strikethrough(true).
	// 		Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
	// 		Render(s)
	// }

	// lists := lipgloss.JoinHorizontal(lipgloss.Top,
	// 	list.Render(
	// 		lipgloss.JoinVertical(lipgloss.Left,
	// 			listHeader("Citrus Fruits to Try"),
	// 			listDone("Grapefruit"),
	// 			listDone("Yuzu"),
	// 			listItem("Citron"),
	// 			listItem("Kumquat"),
	// 			listItem("Pomelo"),
	// 		),
	// 	),

	// 	list.Width(columnWidth).Render(
	// 		lipgloss.JoinVertical(lipgloss.Left,
	// 			listHeader("Actual Lip Gloss Vendors"),
	// 			listItem("Glossier"),
	// 			listItem("Claire‘s Boutique"),
	// 			listDone("Nyx"),
	// 			listItem("Mac"),
	// 			listDone("Milk"),
	// 		),
	// 	),
	// doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, lists, colors))
