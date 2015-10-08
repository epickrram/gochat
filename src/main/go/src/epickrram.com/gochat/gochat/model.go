package main

import (
	"math"
	"sort"
)

const (
	CONTACT_WINDOW int = 0
	CHAT_WINDOW    int = 1
)

const (
	AVAILABLE   int = 2
	UNAVAILABLE int = 3
)

type Contact struct {
	id   string
	name string
}

type Contacts struct {
	contacts      ContactSet
	selectedIndex int
}

type ContactSet []Contact

func (slice ContactSet) Len() int {
	return len(slice)
}

func (slice ContactSet) Less(i, j int) bool {
	return slice[i].id < slice[j].id
}

func (slice ContactSet) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func removeContact(contactId string, contacts *Contacts) {
	indexToRemove := -1
	for idx, contact := range contacts.contacts {
		if contact.id == contactId {
			indexToRemove = idx
			break
		}
	}

	if indexToRemove != -1 {
		contacts.contacts = append(contacts.contacts[:indexToRemove], contacts.contacts[indexToRemove+1:]...)
		sort.Sort(contacts.contacts)
	}
}

func addContact(contactId, displayName string, contacts *Contacts) {
	shouldAdd := true
	for _, contact := range contacts.contacts {
		if contact.id == contactId {
			shouldAdd = false
			break
		}
	}
	if shouldAdd {
		contacts.contacts = append(contacts.contacts, Contact{contactId, displayName})
		sort.Sort(contacts.contacts)
	}
}

type State struct {
	tabs        []Tab
	composition []string
	contacts    Contacts
	viewState   int
	renderer    Renderer
}

func NewState(renderer Renderer) *State {
	state := &State{make([]Tab, 0), make([]string, 1), Contacts{make([]Contact, 0), -1}, CONTACT_WINDOW, renderer}
	state.composition[0] = ""
	return state
}

func (state *State) OnContactEvent(contactId, displayName string, status int) {
	if status == AVAILABLE {
		addContact(contactId, displayName, &state.contacts)
	} else if status == UNAVAILABLE {
		removeContact(contactId, &state.contacts)
	}

	state.renderer.RenderState(state)
}

func (state *State) OnMessageDelivery(conversationId, source, message string) {
	conversationTab := findConversationTab(state.tabs, conversationId)
	if conversationTab == nil {
		newTab := Tab{conversationId, false, NewContentBuffer(100)}
		state.tabs = append(state.tabs, newTab)
		conversationTab = &newTab
	}

	conversationTab.contentBuffer.addLine(message)

	state.renderer.RenderState(state)
}

func (state *State) SendKey(c rune) {
	switch state.viewState {
	case CONTACT_WINDOW:
		break
	case CHAT_WINDOW:
		if c == '\n' {
			state.composition = append(state.composition, "")
		} else {
			currentLine := state.composition[len(state.composition)-1]
			state.composition[len(state.composition)-1] = currentLine + string(c)
		}
		break
	}

	state.renderer.RenderState(state)
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

func (contentBuffer *ContentBuffer) getLastContent(max int) []string {
	content := contentBuffer.getContent()
	if max < len(content) {
		return content[len(content) - max:]
	}

	return content
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

	if len(tabs) != 0 {
		tabs[0].active = true
		return &tabs[0]
	}

	return nil
}
