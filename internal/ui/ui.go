package ui

import (
	"atlas.sand/internal/sim"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TickMsg time.Time

var (
	sandStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
	waterStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("33"))
	wallStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	fireStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	saltStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	emptyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("0"))
)

type Model struct {
	sim          *sim.Simulation
	selectedType sim.ParticleType
	paused       bool
	width        int
	height       int
	mouseX       int
	mouseY       int
	isMouseDown  bool
}

func NewModel() Model {
	return Model{
		selectedType: sim.Sand,
		paused:       false,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tick(), tea.EnterAltScreen)
}

func tick() tea.Cmd {
	return tea.Every(time.Millisecond*50, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "1":
			m.selectedType = sim.Sand
		case "2":
			m.selectedType = sim.Wall
		case "3":
			m.selectedType = sim.Water
		case "4":
			m.selectedType = sim.Fire
		case "5":
			m.selectedType = sim.Salt
		case " ":
			m.paused = !m.paused
		case "c":
			if m.sim != nil {
				m.sim.Clear()
			}
		}

	case tea.MouseMsg:
		m.mouseX = msg.X
		m.mouseY = msg.Y - 3 // Adjust for header
		
		if msg.Type == tea.MouseLeft {
			m.isMouseDown = true
		}
		if msg.Type == tea.MouseRelease {
			m.isMouseDown = false
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Re-init sim with new size
		// Header is 3 lines, Footer is 2 lines
		simW := m.width - 2
		simH := m.height - 6
		if simW > 0 && simH > 0 {
			newSim := sim.NewSimulation(simW, simH)
			if m.sim != nil {
				// Copy old state if possible or just clear
			}
			m.sim = newSim
		}

	case TickMsg:
		if m.isMouseDown && m.sim != nil {
			// Drop 3x3 brush
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					m.sim.SetCell(m.mouseX+dx, m.mouseY+dy, m.selectedType)
				}
			}
		}
		if !m.paused && m.sim != nil {
			m.sim.Update()
		}
		return m, tick()
	}

	return m, nil
}

func (m Model) View() string {
	if m.sim == nil {
		return "Initializing simulation..."
	}

	var s strings.Builder

	// Header
	s.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render(" Atlas Sand "))
	s.WriteString(" | ")
	s.WriteString(m.renderSelected())
	if m.paused {
		s.WriteString(" (PAUSED)")
	}
	s.WriteString("\n\n")

	// Grid
	for y := 0; y < m.sim.Height; y++ {
		s.WriteString(" ")
		for x := 0; x < m.sim.Width; x++ {
			cell := m.sim.Grid[y][x]
			s.WriteString(m.renderCell(cell))
		}
		s.WriteString("\n")
	}

	// Footer
	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(" 1-5: Materials | Click: Drop | Space: Pause | c: Clear | q: Quit"))

	return lipgloss.NewStyle().Padding(1, 1).Render(s.String())
}

func (m Model) renderSelected() string {
	switch m.selectedType {
	case sim.Sand:
		return sandStyle.Render("Sand [1]")
	case sim.Wall:
		return wallStyle.Render("Wall [2]")
	case sim.Water:
		return waterStyle.Render("Water [3]")
	case sim.Fire:
		return fireStyle.Render("Fire [4]")
	case sim.Salt:
		return saltStyle.Render("Salt [5]")
	default:
		return "Unknown"
	}
}

func (m Model) renderCell(c sim.Cell) string {
	char := " "
	style := emptyStyle

	switch c.Type {
	case sim.Wall:
		char = "█"
		style = wallStyle
	case sim.Sand:
		char = "░"
		style = sandStyle
	case sim.Water:
		char = "≈"
		style = waterStyle
	case sim.Fire:
		char = "!"
		style = fireStyle
	case sim.Salt:
		char = "▒"
		style = saltStyle
	}

	return style.Render(char)
}
