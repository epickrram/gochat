package main

import (
	"github.com/nsf/termbox-go"
	"log"
)

const (
	HEADER_OFFSET            = 0
	TOP_BORDER        int    = 4
	COMPOSITION_PANEL int    = 6
	HEADER_MESSAGE    string = "Terminal Chat"
)

type Renderer interface {
	RenderState(state *State)
	OnViewPortResize(width, height int)
	Sync()
}

type TermBoxRenderer struct {
	width  int
	height int
}

func (renderer *TermBoxRenderer) RenderState(state *State) {
	width = renderer.width
	height = renderer.height
	availableHeight := height - TOP_BORDER - COMPOSITION_PANEL - 2

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawHeader(width)

	drawHorizontalLine(0, HEADER_OFFSET+1, width)
	drawHorizontalLine(0, height-COMPOSITION_PANEL, width)

	switch state.viewState {
	case CONTACT_WINDOW:
		yOffset := HEADER_OFFSET + 2
		for _, contact := range state.contacts.contacts {
			if yOffset < availableHeight+(HEADER_OFFSET+2) {
				writeStringAt(1, yOffset, contact.name)
				yOffset++
			}
		}
	case CHAT_WINDOW:
		drawTabs(width, state.tabs)
		drawHorizontalLine(0, HEADER_OFFSET+3, width)
		drawTabContent(findActiveTab(state.tabs), availableHeight)
	}

	termbox.Flush()
}

func (renderer *TermBoxRenderer) OnViewPortResize(width, height int) {
	renderer.width = width
	renderer.height = height
}

func (renderer *TermBoxRenderer) Sync() {
	termbox.Sync()
}

func drawHorizontalLine(x, y, width int) {
	limit := x + width
	for x < limit {
		termbox.SetCell(x, y, '-', termbox.ColorDefault, termbox.ColorDefault)
		x++
	}
}

func drawTabContent(tab *Tab, availableHeight int) {
	if tab == nil {
		log.Print("active tab is null")
		return
	}
	yOffset := HEADER_OFFSET + 5
	messages := tab.contentBuffer.getLastContent(availableHeight)

	if len(messages) < availableHeight {
		yOffset += availableHeight - len(messages)
	}
	log.Printf("Messages for tab %v: %v, at %v", tab.name, messages, yOffset)

	for _, message := range tab.contentBuffer.getLastContent(availableHeight) {
		if yOffset < availableHeight - COMPOSITION_PANEL {
			writeStringAt(0, yOffset, message)
			yOffset++
		}
	}
}

func drawTabs(width int, tabs []Tab) {
	tabWidth := width / len(tabs)
	xOffset := 0
	for _, tab := range tabs {
		writeStringAt(xOffset, HEADER_OFFSET+2, trimString(tabWidth, "| "+tab.name))
		xOffset += width
	}
}

func trimString(width int, text string) string {
	if len(text) > width {
		return text[0:width]
	}
	return text
}

func drawHeader(width int) {
	xOffset := width/2 - len(HEADER_MESSAGE)/2
	writeStringAt(xOffset, HEADER_OFFSET, HEADER_MESSAGE)
}

func writeStringAt(xOffset, yOffset int, message string) {
	for _, c := range message {
		termbox.SetCell(xOffset, yOffset, c, termbox.ColorDefault, termbox.ColorDefault)
		xOffset += 1
	}
}

type NoOpRenderer struct {
}

func (renderer *NoOpRenderer) RenderState(state *State) {

}

func (renderer *NoOpRenderer) OnViewPortResize(width, height int) {

}

func (renderer *NoOpRenderer) Sync() {

}
