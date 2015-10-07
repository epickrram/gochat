package main

type ModelInputEvent interface {
	Execute(state *State)
}

func GetEventChannelToModel(renderer Renderer) chan ModelInputEvent {
	eventChannel := make(chan ModelInputEvent)
	state := NewState(renderer)

	go func() {
		for {
			event := <-eventChannel
			dispatchEvent(state, event)
		}
	}()

	return eventChannel
}

func dispatchEvent(state *State, event ModelInputEvent) {
	event.Execute(state)
}
