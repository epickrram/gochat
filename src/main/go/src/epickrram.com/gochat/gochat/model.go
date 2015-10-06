package main

import (
	"math"
)

const (
	CONTACT_WINDOW int = 0
	CHAT_WINDOW    int = 1
)

const (
	AVAILABLE   int = 2
	UNAVAILABLE int = 3
)

type State struct {
	tabs        []Tab
	composition []string
	contacts    map[string]string
	viewState   int
}

func NewState() *State {
	state := &State{make([]Tab, 0), make([]string, 1), make(map[string]string), CHAT_WINDOW}
	state.composition[0] = ""
	return state
}

func (state *State) OnContactEvent(contactId, displayName string, status int) {
	if status == AVAILABLE {
		state.contacts[contactId] = displayName
	} else if status == UNAVAILABLE {
		delete(state.contacts, contactId)
	}
}

func (state *State) OnMessageDelivery(conversationId, source, message string) {
	conversationTab := findConversationTab(state.tabs, conversationId)
	if conversationTab == nil {
		newTab := Tab{conversationId, false, NewContentBuffer(100)}
		state.tabs = append(state.tabs, newTab)
		conversationTab = &newTab
	}

	conversationTab.contentBuffer.addLine(message)
}

func (state *State) SendKey(c string) {
	switch state.viewState {
	case CONTACT_WINDOW:
		break
	case CHAT_WINDOW:
		if c == "\n" {
			state.composition = append(state.composition, "")
		} else {
			currentLine := state.composition[len(state.composition)-1]
			state.composition[len(state.composition)-1] = currentLine + c
		}
		break
	}
}

func (state *State) DisplayChatWindow() {
	state.viewState = CHAT_WINDOW
}

func (state *State) DisplayContactWindow() {
	state.viewState = CONTACT_WINDOW
}

type ContentBuffer struct {
	content []string
	pointer int
}

type Tab struct {
	name          string
	active        bool
	contentBuffer *ContentBuffer
}

func (contentBuffer *ContentBuffer) getContent() []string {
	length := int(math.Min(float64(len(contentBuffer.content)), float64(contentBuffer.pointer)))
	view := make([]string, length)
	counter := 0
	for counter < length {
		sourceIndex := contentBuffer.pointer + counter - length
		view[counter] = contentBuffer.content[sourceIndex%len(contentBuffer.content)]
		counter++
	}

	return view
}

func (contentBuffer *ContentBuffer) addLine(line string) {
	contentBuffer.content[contentBuffer.pointer%len(contentBuffer.content)] = line
	contentBuffer.pointer++
}

func NewContentBuffer(length int) *ContentBuffer {
	return &ContentBuffer{make([]string, length), 0}
}

func (contentBuffer *ContentBuffer) modContentIndex(index int) int {
	return index % len(contentBuffer.content)
}

func findConversationTab(tabs []Tab, conversationId string) *Tab {
	for _, tab := range tabs {
		if tab.name == conversationId {
			return &tab
		}
	}

	return nil
}

func findActiveTab(tabs []Tab) *Tab {
	for _, tab := range tabs {
		if tab.active {
			return &tab
		}
	}

	return nil
}
