package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Amr-elwetaidy/pokedexcli/internal/pokeapi"
	"github.com/Amr-elwetaidy/pokedexcli/internal/pokecache"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

type model struct {
	viewport     viewport.Model
	textInput    textinput.Model
	replState    *replState
	commands     map[string]cliCommand
	history      []string
	historyIndex int
	ready        bool
	content      string // Store the full content history
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter command..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 80

	// We'll set the viewport size later in the update loop
	vp := viewport.New(0, 0)
	initialContent := titleStyle.Render("Welcome to the Pokedex!") + "\nType 'help' for a list of commands."
	vp.SetContent(initialContent)

	return model{
		textInput:    ti,
		viewport:     vp,
		replState:    initialReplState(),
		commands:     getCommands(),
		history:      []string{},
		historyIndex: 0,
		content:      initialContent,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 3
		m.ready = true
		// Pass the resize message to the viewport
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		switch msg.Type {
		// These keys are captured for app control and are not passed down.
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyUp, tea.KeyDown:
			if msg.Type == tea.KeyUp {
				if len(m.history) > 0 {
					if m.historyIndex > 0 {
						m.historyIndex--
					}
					m.textInput.SetValue(m.history[m.historyIndex])
					m.textInput.SetCursor(len(m.textInput.Value()))
				}
			}
			if msg.Type == tea.KeyDown {
				if len(m.history) > 0 {
					if m.historyIndex < len(m.history)-1 {
						m.historyIndex++
						m.textInput.SetValue(m.history[m.historyIndex])
						m.textInput.SetCursor(len(m.textInput.Value()))
					} else {
						m.historyIndex = len(m.history)
						m.textInput.Reset()
					}
				}
			}
			// Return here to prevent the key press from being sent to other components.
			return m, nil

		case tea.KeyEnter:
			// Handle command execution...
			input := m.textInput.Value()
			m.textInput.Reset()
			if input == "" {
				break
			}

			if len(m.history) == 0 || m.history[len(m.history)-1] != input {
				m.history = append(m.history, input)
			}
			m.historyIndex = len(m.history)

			words := cleanInput(input)
			if len(words) == 0 {
				break
			}
			commandName := words[0]
			args := words[1:]

			command, exists := m.commands[commandName]
			if !exists {
				// Append unknown command error to content history
				if m.content != "" && !strings.HasSuffix(m.content, "\n") {
					m.content += "\n"
				}
				m.content += fmt.Sprintf("> %s\n", input)
				m.content += errorStyle.Render(fmt.Sprintf("Unknown command: %s", commandName))
				m.viewport.SetContent(m.content)
				m.viewport.GotoBottom()
				break
			}

			// Append the command to the content history
			if m.content != "" && !strings.HasSuffix(m.content, "\n") {
				m.content += "\n"
			}
			m.content += fmt.Sprintf("> %s\n", input)

			output, err := command.callback(m.replState, args)
			if err != nil {
				if errors.Is(err, ErrExit) {
					return m, tea.Quit
				}
				m.content += errorStyle.Render(fmt.Sprintf("Error: %v", err))
			} else {
				m.content += output
			}

			m.viewport.SetContent(m.content)
			m.viewport.GotoBottom()
			// Return here to prevent the key press from being sent to other components.
			return m, nil

		// Default case for all other keys (typing)
		default:
			// Pass the key press ONLY to the text input.
			m.textInput, cmd = m.textInput.Update(msg)
			cmds = append(cmds, cmd)
		}

	// For all other messages (like mouse scrolls), pass them ONLY to the viewport.
	default:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}
	return fmt.Sprintf(
		"%s\n%s",
		m.viewport.View(),
		m.textInput.View(),
	)
}

func cleanInput(text string) []string {
	lowercase := strings.ToLower(text)
	words := strings.Fields(lowercase)
	return words
}

func initialReplState() *replState {
	return &replState{
		config: &pokeapi.Config{
			Next: strPtr("https://pokeapi.co/api/v2/location-area/"),
		},
		cache:   pokecache.NewCache(10 * time.Second),
		pokedex: make(map[string]pokeapi.Pokemon),
	}
}

func getCommands() map[string]cliCommand {
	commands := map[string]cliCommand{
		"map": {
			name:        "map",
			description: "Get the next page of locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Get the previous page of locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore <area_name>",
			description: "Lists the pokemon in a given location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <pokemon_name>",
			description: "Attempt to catch a pokemon and add it to your pokedex",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon_name>",
			description: "View details about a caught pokemon",
			callback:    commandInspect,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback: func(s *replState, args []string) (string, error) {
			return commandHelp(s, commands, args)
		},
	}
	return commands
}

func strPtr(s string) *string {
	return &s
}

func main() {
	// check if we're in a pipe
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	// if we're getting piped input, run in non-interactive mode
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		state := initialReplState()
		commands := getCommands()
		for scanner.Scan() {
			line := scanner.Text()
			words := cleanInput(line)
			if len(words) == 0 {
				continue
			}
			commandName := words[0]
			args := words[1:]
			if command, ok := commands[commandName]; ok {
				out, err := command.callback(state, args)
				if err != nil {
					if errors.Is(err, ErrExit) {
						return
					}
					fmt.Fprintln(os.Stderr, "error:", err)
					continue
				}
				if out != "" {
					fmt.Println(out)
				}
			} else {
				fmt.Fprintln(os.Stderr, "unknown command:", commandName)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		return
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
