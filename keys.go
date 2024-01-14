package main

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Up        key.Binding
	Down      key.Binding
	Quit      key.Binding
	Help      key.Binding
	Clear     key.Binding
	SelectAll key.Binding
	Tab       key.Binding
	Toggle    key.Binding

	Remove  key.Binding
	Restart key.Binding
	Kill    key.Binding
	Stop    key.Binding
	Start   key.Binding
	Pause   key.Binding
	Unpause key.Binding

	Select key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Help,
		k.Quit,
		k.Clear,
		k.SelectAll,
		k.Tab,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.Help,
			k.Quit,
			k.Clear,
			k.SelectAll,
			k.Tab,
		},
		{
			k.Up,
			k.Down,
			k.Select,
		},
		{
			k.Remove,
			k.Restart,
			k.Kill,
			k.Stop,
			k.Start,
			k.Pause,
			k.Unpause,
		},
	}
}

var keys = keyMap{

	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Clear: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "clear selection"),
	),
	SelectAll: key.NewBinding(
		key.WithKeys("shift+a"),
		key.WithHelp("shift+a", "select all"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch page"),
	),
	Toggle: key.NewBinding(
		key.WithKeys(" ", "enter"),
		key.WithHelp("space/enter", "toggle selection"),
	),
	Remove: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "remove container"),
	),
	Restart: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "restart container"),
	),
	Kill: key.NewBinding(
		key.WithKeys("shift+k"),
		key.WithHelp("shift+k", "kill container"),
	),
	Stop: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "stop container"),
	),
	Start: key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "start container"),
	),
	Pause: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "pause container"),
	),
	Unpause: key.NewBinding(
		key.WithKeys("shift+p"),
		key.WithHelp("shift+p", "unpause container"),
	),
}
