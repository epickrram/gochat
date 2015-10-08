package main

import (
	"github.com/nsf/termbox-go"
	"log"
)


/*
0 header
1 -hr-
2 | tab 1   | tab 2
3 -hr-
4


(height -6) -hr-


x4

(height - 1) -hr-
 */


const (
	HEADER_OFFSET int           = 0
	HEADER_HR int = 1
	TAB_OFFSET int = 2
	TAB_HR int = 3
	CONTENT_START int = 4
	TOP_BORDER        int    = 4
	COMPOSITION_PANEL_HEIGHT int    = 4
	HEADER_MESSAGE    string = "Terminal Chat"
	COMPOSITION_PANEL_HR int = -6
	COMPOSITION_PANEL_START int = -5
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
	availableHeight := height - TOP_BORDER - COMPOSITION_PANEL_HEIGHT - 2
	contentFooter := height - 7

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawHeader(width)

	drawHorizontalLine(0, HEADER_HR, width)
	drawHorizontalLine(0, height + COMPOSITION_PANEL_HR, width)

	switch state.viewState {
	case CONTACT_WINDOW:
		yOffset := CONTENT_START
		for _, contact := range state.contacts.contacts {
			if yOffset < contentFooter {
				writeStringAt(1, yOffset, contact.name)
				yOffset++
			}
		}
	case CHAT_WINDOW:
		drawTabs(width, state.tabs)
		drawHorizontalLine(0, TAB_HR, width)
		drawTabContent(findActiveTab(state.tabs), availableHeight, contentFooter)
		drawCompositionContent(state.composition, height)
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

func drawCompositionContent(messages []string, height int) {
	log.Printf("composition: %v", messages)
	yOffset := height - 1
	idx := len(messages) - 1
	limit := height + COMPOSITION_PANEL_START
	log.Printf("yOffset: %v, limit: %v", yOffset, limit)
	for idx >= 0 && yOffset > limit {
		writeStringAt(0, yOffset, messages[idx])
		idx--
		yOffset--
	}
}

func drawTabContent(tab *Tab, availableHeight, contentFooter int) {
	if tab == nil {
		log.Print("active tab is null")
		return
	}
	yOffset := CONTENT_START
	availableHeight = contentFooter - CONTENT_START
	messages := tab.contentBuffer.getLastContent(availableHeight)

//	log.Printf("yOffset starts at %v, availableHeight is %v", yOffset, availableHeight)

	if len(messages) < availableHeight {
		emptyLines := availableHeight - len(messages)
		yOffset += emptyLines
//		log.Printf("yOffset updated to %v", yOffset)
	}
//	log.Printf("Messages for tab %v: %v, at %v", tab.name, messages, yOffset)

	for _, message := range messages {
		if yOffset <= contentFooter {
			writeStringAt(0, yOffset, message)
			yOffset++
		}
	}
}

func drawTabs(width int, tabs []Tab) {
	tabWidth := width / len(tabs)
	xOffset := 0
	for _, tab := range tabs {
		writeStringAt(xOffset, TAB_OFFSET, trimString(tabWidth, "| "+tab.name))
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
