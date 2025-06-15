package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/farmingengineers/harvest/cmd/input/filter"
)

var (
	cropsFile = flag.String("crops", "crops.txt", "path to crops file")
	outFile   = flag.String("out", "", "path to output file")
)

type model struct {
	crops         []string
	filteredCrops []string
	cropInput     textinput.Model
	quantityInput textinput.Model
	selectedCrop  string
	selectedIndex int    // -1 is a special value, use 'resolveSelectedIndex()' to get a value that can be used as an index into filteredC
	state         string // "crop" or "quantity"
	output        []string
}

func (m model) resolveSelectedIndex() int {
	if m.selectedIndex == -1 {
		if len(m.filteredCrops) > 1 {
			return 1
		}
		return 0
	}
	return m.selectedIndex
}

func initialModel(crops []string) model {
	ti := textinput.New()
	ti.Placeholder = "Type crop name..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	qi := textinput.New()
	qi.Placeholder = "Type quantity (e.g. 60#, 88.8 lb, 12 pt)..."
	qi.CharLimit = 156
	qi.Width = 50

	return model{
		crops:         crops,
		filteredCrops: crops,
		selectedIndex: -1,
		cropInput:     ti,
		quantityInput: qi,
		state:         "crop",
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+d":
			// Write output to file
			if *outFile != "" {
				f, err := os.Create(*outFile)
				if err != nil {
					fmt.Printf("Error creating output file: %v\n", err)
					os.Exit(1)
				}
				defer f.Close()
				for _, line := range m.output {
					fmt.Fprintln(f, line)
				}
			}
			return m, tea.Quit
		}

		if m.state == "crop" {
			switch msg.String() {
			case "enter":
				if len(m.filteredCrops) > 0 {
					m.selectedCrop = m.filteredCrops[m.resolveSelectedIndex()]
					m.state = "quantity"
					m.quantityInput.Focus()
					m.cropInput.Blur()
				} else if m.cropInput.Value() != "" {
					m.selectedCrop = m.cropInput.Value()
					m.state = "quantity"
					m.quantityInput.Focus()
					m.cropInput.Blur()
				}
			case "up":
				if si := m.resolveSelectedIndex(); si > 0 {
					m.selectedIndex = si - 1
				}
			case "down":
				if si := m.resolveSelectedIndex(); si < len(m.filteredCrops)-1 {
					m.selectedIndex = si + 1
				}
			}
		} else if m.state == "quantity" {
			switch msg.String() {
			case "enter":
				if m.quantityInput.Value() != "" {
					m.output = append(m.output, fmt.Sprintf("%s | %s", m.selectedCrop, m.quantityInput.Value()))
				}
				m.state = "crop"
				m.cropInput.Focus()
				m.quantityInput.Blur()
				m.quantityInput.Reset()
				m.cropInput.Reset()
				m.selectedCrop = ""
				m.selectedIndex = -1
			}
		}
	}

	if m.state == "crop" {
		m.cropInput, cmd = m.cropInput.Update(msg)
		query := m.cropInput.Value()
		m.filteredCrops = make([]string, 0, 6)
		m.filteredCrops = append(m.filteredCrops, query)
		m.filteredCrops = append(m.filteredCrops, filter.Crops(m.crops, query, 5)...)
		if m.selectedIndex >= len(m.filteredCrops) {
			m.selectedIndex = -1
		}
	} else {
		m.quantityInput, cmd = m.quantityInput.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	var s strings.Builder

	s.WriteString("Harvest Input\n\n")
	for _, line := range m.output {
		s.WriteString(line + "\n")
	}
	s.WriteString("\n-----------------------\n\n")

	if m.state == "crop" {
		s.WriteString(m.cropInput.View())
		s.WriteString("\n\n")
		if len(m.filteredCrops) > 0 {
			s.WriteString("Matches:\n")
			for i, crop := range m.filteredCrops {
				si := m.resolveSelectedIndex()
				if i == si {
					s.WriteString("> ")
				} else {
					s.WriteString("  ")
				}
				s.WriteString(crop)
				s.WriteString("\n")
			}
		}
	} else {
		s.WriteString(fmt.Sprintf("Selected crop: %s\n", m.selectedCrop))
		s.WriteString(m.quantityInput.View())
	}

	s.WriteString("\n\n")
	s.WriteString("Press Ctrl+D when done\n")

	return s.String()
}

func main() {
	flag.Parse()

	if *outFile == "" {
		fmt.Println("Error: --out flag is required")
		os.Exit(1)
	}

	// Read crops file
	file, err := os.Open(*cropsFile)
	if err != nil {
		fmt.Printf("Error opening crops file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var crops []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" && !strings.HasPrefix(line, "Week") {
			crops = append(crops, line)
		}
	}

	p := tea.NewProgram(initialModel(crops))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
