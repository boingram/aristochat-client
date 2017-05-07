package aristochat

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"strings"
)

type UI struct {
	client   *Client
	messages []string
}

var INPUT_WINDOW_HEIGHT = 3

func (ui *UI) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("messages", -1, -1, maxX, maxY-INPUT_WINDOW_HEIGHT); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Autoscroll = true
		v.Wrap = true
		v.SetCursor(-1, -1)
		return nil
	}
	if v, err := g.SetView("stdin", -1, maxY-INPUT_WINDOW_HEIGHT, maxX, maxY+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView("stdin"); err != nil {
			return err
		}
		v.Autoscroll = true
		v.Editable = true
		v.Wrap = true
	}
	return nil
}

func (ui *UI) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("stdin", gocui.KeyEnter, gocui.ModNone, ui.sendMessage); err != nil {
		return err
	}

	return nil
}

func (ui *UI) sendMessage(_ *gocui.Gui, v *gocui.View) error {
	msg := strings.TrimSpace(v.Buffer())
	v.Clear()
	v.SetCursor(0, 0)
	return ui.client.SendMessage(msg)
}

func quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}

func (ui *UI) listenForPayloads(g *gocui.Gui) error {
	ch := ui.client.Channel()
	for {
		payload := <-ch
		str := fmt.Sprintf("%s > %s", payload.Username, payload.Body)
		ui.messages = append(ui.messages, str)
		writeMessage(g, str)
	}
}

func writeMessage(g *gocui.Gui, msg string) {
	g.Execute(func(g *gocui.Gui) error {
		v, err := g.View("messages")
		if err != nil {
			return err
		}
		fmt.Fprintln(v, msg)
		return nil
	})
}

func StartUI(client *Client) error {
	ui := UI{
		client:   client,
		messages: []string{},
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(ui.layout)

	if err := ui.keybindings(g); err != nil {
		return err
	}

	go ui.client.Listen()
	go ui.listenForPayloads(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}
