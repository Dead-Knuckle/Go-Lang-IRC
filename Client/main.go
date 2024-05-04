package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	serverAddress = "localhost:6667"
)

type Message struct {
	Msg      string `json:"msg"`
	Username string `json:"username"`
}

func main() {
	p := tea.NewProgram(initialModel())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	viewport     viewport.Model
	messages     []string
	textarea     textarea.Model
	senderStyle  lipgloss.Style
	receiveStyle lipgloss.Style
	serverStyle  lipgloss.Style
	conn         net.Conn
	msgChan      chan string // Channel to receive new messages
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 0

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(50, 10)
	vp.SetContent(`If you see this we have a problem :D`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Failed to connect to the server:", err)
		os.Exit(1)
	}

	msgChan := make(chan string) // Create a channel for receiving messages

	m := model{
		textarea:     ta,
		messages:     []string{},
		viewport:     vp,
		senderStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#bde0fe")),
		receiveStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#b5e48c")),
		serverStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#EAC435")),
		conn:         conn,
		msgChan:      msgChan,
	}

	go receiveMessages(conn, msgChan) // Start receiving messages

	return m
}

func receiveMessages(conn net.Conn, msgChan chan string) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		var receivedMessage Message
		if err := json.Unmarshal([]byte(message), &receivedMessage); err != nil {
			fmt.Println("Error decoding message:", err)
			continue
		}

		if receivedMessage.Username == "SERVER" && receivedMessage.Msg == "HEARTBEAT" {
			// Respond to the server's heartbeat
			_, err := fmt.Fprintf(conn, "HEARTBEAT\n")
			if err != nil {
				fmt.Println("Error sending heartbeat response:", err)
			}
			continue // Skip adding the server's heartbeat message to the msgChan
		}

		msgChan <- message // Send other messages to the channel
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	// Check if there are new messages in the channel
	select {
	case message := <-m.msgChan:

		var receivedMessage Message
		if err := json.Unmarshal([]byte(message), &receivedMessage); err != nil {
			fmt.Println("Error decoding message:", err)
		} else {
			if receivedMessage.Username == "SERVER" && receivedMessage.Msg == "HEARTBEAT" {
				m.messages = append(m.messages, m.serverStyle.Render("["+receivedMessage.Username+"]: ")+receivedMessage.Msg)
				_, err := fmt.Fprintf(m.conn, "%s", "HEARTBEAT")
				if err != nil {
					fmt.Println("Error sending heartbeat:", err)
				}
			} else if receivedMessage.Username == "SERVER" {
				m.messages = append(m.messages, m.serverStyle.Render("["+receivedMessage.Username+"]: ")+receivedMessage.Msg)
			} else {
				m.messages = append(m.messages, m.receiveStyle.Render(receivedMessage.Username+": ")+receivedMessage.Msg)
			}
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.viewport.GotoBottom()
		}
	default:
		// No new messages, continue with normal update
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Fprintf(m.conn, "QUIT")
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			message := m.textarea.Value()

			_, err := fmt.Fprintf(m.conn, "%s\n", message)
			if err != nil {
				fmt.Println("Error sending message:", err)
				return m, tea.Quit
			}

			if message == "QUIT" {
				fmt.Println("Quitting...")
				fmt.Fprintf(m.conn, "QUIT")
				return m, tea.Quit
			}

			m.messages = append(m.messages, m.senderStyle.Render("You: ")+message)
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom() // Scroll to the end after adding a new message
		}
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}
