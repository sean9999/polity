package main

import (
	"context"
	"fmt"
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/udp4"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type model struct {
	isConnected bool
	status      string
	input       string
	output      string
	quitting    bool
	nick        string
	self        *polity.Principal[*udp4.Network]
	peer        *polity.Peer[*udp4.Network]
	verbosity   uint
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
type connectMessage struct{}
type commandMessage string
type quitMessage struct{}

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

		case "c":
			return m, func() tea.Msg {
				// Simulate connecting to polityd
				connected := connectToPolityd()
				if connected {
					return connectMessage{}
				}
				return errorStyle.Render("Failed to connect to polityd")
			}

		case "r": // Example of issuing a read command
			if !m.isConnected {
				return m, nil
			}
			return m, func() tea.Msg {
				res, err := issueCommand("read")
				if err != nil {
					return errorStyle.Render("Error: " + err.Error())
				}
				return commandMessage(res)
			}

		case "s": // Example of issuing a status command
			if !m.isConnected {
				return m, nil
			}
			return m, func() tea.Msg {
				//res, err := issueCommand("status")
				//if err != nil {
				//	return errorStyle.Render("Error: " + err.Error())
				//}
				res := ""
				for pub, info := range m.self.Peers.Entries() {
					res += fmt.Sprintf("%s - %s\n", pub.Nickname(), info.Addr.String())
				}
				return commandMessage(res)
			}
		}

	case connectMessage:
		m.isConnected = true
		m.status = "Connected to polityd"
		m.output = "Successfully connected to the node."
		return m, nil

	case commandMessage:
		m.output = string(msg)
		return m, nil

		//case errorStyle:
		//	m.output = string(msg)
		//	return m, nil
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
	status := fmt.Sprintf("Status: %s\n", textStyle.Render(m.status))
	commands := textStyle.Render("[c] Connect  [s] Status  [r] Read  [q] Quit") + "\n"
	output := fmt.Sprintf("Output:\n%s\n\n", m.output)
	return header + divider + status + divider + commands + divider + output
}

// Connect to the polityd process
func connectToPolityd() bool {
	// Replace with actual connection logic
	time.Sleep(1 * time.Second) // Simulate connection delay
	return true
}

// Issue a command to polityd
func issueCommand(cmd string) (string, error) {
	// Replace with actual command API logic
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	switch cmd {
	case "status":

		return "Node is running", nil
	case "read":
		return "Command executed successfully", nil
	default:
		return "", fmt.Errorf("unknown command: %s", cmd)
	}

	<-ctx.Done() // Simulate executing
	return "Command executed", nil
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

	//if err := p.Start(); err != nil {
	//	log.Fatalf("Error starting TUI: %v", err)
	//	os.Exit(1)
	//}
}
