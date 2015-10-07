package main

import (
	"reflect"
	"testing"
)

func getState() *State {
	return NewState(&NoOpRenderer{})
}

func TestShouldHandleRemovingUnknownContact(t *testing.T) {
	state := getState()
	state.OnContactEvent("contact-0", "First Last", UNAVAILABLE)
	if len(state.contacts.contacts) != 0 {
		t.Errorf("Contacts should be empty: %v", state.contacts)
	}
}

func TestShouldHandleAdditionOfContact(t *testing.T) {
	state := getState()
	state.OnContactEvent("contact-0", "First Last", AVAILABLE)

	if len(state.contacts.contacts) != 1 {
		t.Errorf("Contacts should be size 1: %v", state.contacts.contacts)
	}
}

func TestAddMessageFromNewConversation(t *testing.T) {
	state := getState()
	state.OnMessageDelivery("conv-0", "sender", "hello world")
	tabFound := false
	for _, tab := range state.tabs {
		if tab.name == "conv-0" {
			tabFound = true
			bufferContent := tab.contentBuffer.getContent()
			if !reflect.DeepEqual(bufferContent, []string{"hello world"}) {
				t.Errorf("Tab content incorrect: %v", bufferContent)
			}
		}
	}

	if !tabFound {
		t.Error("Could not find tab for conversation")
	}
}

func TestAddMessageFromExistingConversation(t *testing.T) {
	state := getState()
	state.OnMessageDelivery("conv-0", "sender", "hello world")
	state.OnMessageDelivery("conv-0", "sender2", "foobar")
	tabFound := false
	for _, tab := range state.tabs {
		if tab.name == "conv-0" {
			tabFound = true
			bufferContent := tab.contentBuffer.getContent()
			if !reflect.DeepEqual(bufferContent, []string{"hello world", "foobar"}) {
				t.Errorf("Tab content incorrect: %v", bufferContent)
			}
		}
	}

	if !tabFound {
		t.Error("Could not find tab for conversation")
	}
}

func TestSendKeysToStateWhenInChatMode(t *testing.T) {
	state := getState()

	state.DisplayChatWindow()
	for _, c := range "hello\nworld" {
		state.SendKey(string(c))
	}

	cmp := state.composition
	if len(cmp) != 2 || cmp[0] != "hello" || cmp[1] != "world" {
		t.Errorf("Composition not as expected: %v", cmp)
	}
}

func TestContentBuffer(t *testing.T) {
	cb := NewContentBuffer(10)
	cb.addLine("one")
	cb.addLine("two")

	content := cb.getContent()
	if len(content) != 2 {
		t.Errorf("Content length incorrect: %v", len(content))
	}
	if content[0] != "one" || content[1] != "two" {
		t.Errorf("Content incorrect: %v", content)
	}
}

func TestContentBufferWrapping(t *testing.T) {
	cb := NewContentBuffer(2)
	cb.addLine("one")
	cb.addLine("two")
	cb.addLine("three")

	content := cb.getContent()
	if len(content) != 2 {
		t.Errorf("Content length incorrect: %v", len(content))
	}
	if content[0] != "two" || content[1] != "three" {
		t.Errorf("Content incorrect: %v", content)
	}
}

func TestContentBufferWrappingExactly(t *testing.T) {
	cb := NewContentBuffer(2)
	cb.addLine("one")
	cb.addLine("two")

	content := cb.getContent()
	if len(content) != 2 {
		t.Errorf("Content length incorrect: %v", len(content))
	}
	if content[0] != "one" || content[1] != "two" {
		t.Errorf("Content incorrect: %v", content)
	}
}
