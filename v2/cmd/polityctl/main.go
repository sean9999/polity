package main

import (
	"crypto/rand"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
	"github.com/sean9999/polity/v2/udp4"
	"log"
)

type model struct {
	//isConnected bool
	//status      string
	//input       string
	quitting   bool
	output     string
	nick       string
	self       *polity.Principal[*udp4.Network]
	selfAsPeer *polity.Peer[*udp4.Network]
	verbosity  uint
}

// Initialize the genericModel
//func initialModel() genericModel {
//	return genericModel{
//		isConnected: false,
//		status:      "Disconnected",
//		input:       "",
//		output:      "",
//		quitting:    false,
//	}
//}

// Bubble Tea styles for the TUI
var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	textStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))
	inputStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("110"))
)

// Message types for the TUI
type commandMessage string
type errorMessage string

// Init initializes the genericModel
func (m model) Init() tea.Cmd {
	// Initialize with no commands
	m.nick = "Ray Charles"
	return nil
}

// Update updates the genericModel
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit // Send quit command to Bubble Tea

		case "d":
			//	dump thyself
			return m, func() tea.Msg {
				e := m.self.Compose([]byte("dump thyself"), m.selfAsPeer, nil)
				e.Subject(subj.DumpThyself)
				i, err := m.self.Send(e)
				if err != nil {
					return errorMessage(fmt.Errorf("failed to send %d bytes. %w", i, err).Error())
				}
				return commandMessage("dump thyself")
			}

		case "D":
			//	tell everyone to dump
			return m, func() tea.Msg {
				e := m.self.Compose([]byte("tell everyone dump thyself"), m.selfAsPeer, nil)
				e.Subject(subj.CmdEveryoneDump)
				i, err := m.self.Send(e)
				if err != nil {
					return errorMessage(fmt.Errorf("failed to send %d bytes. %w", i, err).Error())
				}
				return commandMessage("tell everyone dump thyself")
			}

		case "b":
			//	broadcast (say hello to all your friends)
			return m, func() tea.Msg {
				e := m.self.Compose([]byte("broadcast"), m.selfAsPeer, nil)
				e.Subject(subj.CmdBroadcast)
				i, err := m.self.Send(e)
				if err != nil {
					return errorMessage(fmt.Errorf("failed to send %d bytes. %w", i, err).Error())
				}
				return commandMessage("broadcast")
			}

		case "f":
			//	ask all your friends who their friends are, and befriend them
			return m, func() tea.Msg {
				e := m.self.Compose([]byte("broadcast"), m.selfAsPeer, nil)
				e.Subject(subj.CmdMakeFriends)
				i, err := m.self.Send(e)
				if err != nil {
					return errorMessage(fmt.Errorf("failed to send %d bytes. %w", i, err).Error())
				}
				return commandMessage("friends of friends")
			}

		case "k":
			// kill yourself
			return m, func() tea.Msg {
				e := m.self.Compose([]byte("kill yourself"), m.selfAsPeer, nil)
				e.Subject(subj.KillYourself)
				e.Message.Sign(rand.Reader, m.self)
				i, err := m.self.Send(e)
				if err != nil {
					return errorMessage(fmt.Errorf("failed to send %d bytes. %w", i, err).Error())
				}
				return commandMessage("kill yourself")
			}

		case "s":
			// sleep (stay alive but stop responding to messages)
			return m, func() tea.Msg {
				e := m.self.Compose([]byte("go to sleep"), m.selfAsPeer, nil)
				e.Subject(subj.Sleep)
				i, err := m.self.Send(e)
				if err != nil {
					return errorMessage(fmt.Errorf("failed to send %d bytes. %w", i, err).Error())
				}
				return commandMessage("go to sleep")
			}
		}

	case commandMessage:
		m.output = string(msg)
		return m, nil

	case errorMessage:
		m.output = errorStyle.Render(string(msg))
		return m, nil
	}

	return m, nil
}

// View renders the genericModel
func (m model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	hdr := fmt.Sprintf("Nickname - %s", m.self.Nickname())

	header := titleStyle.Render(hdr) + "\n"
	divider := textStyle.Render("------------------------------") + "\n"
	commands := textStyle.Render(`
[d] dump thyself
[b] broadcast hello
[f] ask for friends
[k] kill yourself
[s] sleep
`)
	output := fmt.Sprintf("\n%s\n\n", m.output)
	return header + divider + divider + commands + divider + output
}

// Main function
func main() {
	m, err := parseFlargs()
	if err != nil {
		panic(err)
	}
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
