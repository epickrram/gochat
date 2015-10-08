package main

import "github.com/mattn/go-xmpp"

type KeyPressEvent struct {
	key rune
}

func (event *KeyPressEvent) Execute(state *State) {
	state.SendKey(event.key)
}

type PresenceUpdateEvent struct {
	presence xmpp.Presence
}

func (event *PresenceUpdateEvent) Execute(state *State) {
	presenceState := AVAILABLE
	if event.presence.Type != "" || event.presence.Show != "" {
		presenceState = UNAVAILABLE
	}
	state.OnContactEvent(event.presence.From, event.presence.From, presenceState)
}

type MessageReceivedEvent struct {
	message xmpp.Chat
}

func (event *MessageReceivedEvent) Execute(state *State) {
	state.OnMessageDelivery(event.message.Remote, event.message.Remote, event.message.Text)
}

type ResizeEvent struct {
	width  int
	height int
}

func (event *ResizeEvent) Execute(state *State) {
	state.renderer.OnViewPortResize(event.width, event.height)
	state.renderer.RenderState(state)
	state.renderer.Sync()
}

type SwitchViewEvent struct {
	viewType int
}

func (event *SwitchViewEvent) Execute(state *State) {
	state.viewState = event.viewType
	state.renderer.RenderState(state)
}
