package modules

import (
	"fmt"
	"time"

	"github.com/devnull-twitch/go-bot/pkg/tmi"
)

type timedMessageMod struct {
	currentIndex  int
	messages      []string
	lastMessage   time.Time
	targetChannel string
}

func TimedMessageMod(messages []string, channel string) tmi.Module {
	return &timedMessageMod{
		currentIndex:  0,
		lastMessage:   time.Now(),
		messages:      messages,
		targetChannel: fmt.Sprintf("#%s", channel),
	}
}

func (m *timedMessageMod) ExternalTrigger(client *tmi.Client) <-chan *tmi.ModuleArgs {
	if len(m.messages) <= 0 {
		return nil
	}

	genChan := tmi.CreateTimeTrigger(1 * time.Minute)
	filledChan := make(chan *tmi.ModuleArgs)
	go func() {
		for {
			// wait for time trigger
			<-genChan

			event := &tmi.ModuleArgs{
				Parameter: map[string]string{"msg": m.messages[m.currentIndex]},
			}

			m.currentIndex++
			if m.currentIndex >= len(m.messages) {
				m.currentIndex = 0
			}

			filledChan <- event
		}
	}()

	return filledChan
}
func (m *timedMessageMod) MessageTrigger(client *tmi.Client, incoming *tmi.IncomingMessage) *tmi.ModuleArgs {
	m.lastMessage = time.Now()
	return nil
}
func (m *timedMessageMod) Handler(client *tmi.Client, args tmi.ModuleArgs) *tmi.OutgoingMessage {
	return &tmi.OutgoingMessage{
		Message: args.Parameter["msg"],
		Channel: m.targetChannel,
	}
}
