package tui

import (
	"fmt"
	"github.com/aldernero/gaul"
	"github.com/aldernero/tui-apps/rule30/pkg/automata"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"os"
	"strings"
	"time"
)

const (
	timeout      = 365 * 24 * time.Hour
	speedLevels  = 20
	defaultSpeed = 10
	minDelay     = 75
	maxDelay     = 1000
)

type Model struct {
	Grid     *automata.Grid
	Noise    gaul.Rng
	Offset   float64
	Timer    timer.Model
	Seed     int64
	Speed    int
	Rows     int
	Cols     int
	WrapEnds bool
	palette  int
}

func StartTea(seed int64) {
	m := Model{
		Grid:  automata.NewGrid(1, 1, seed),
		Noise: gaul.NewRng(seed),
		Seed:  seed,
		Speed: defaultSpeed,
		Timer: timer.NewWithInterval(timeout, GetDelay(defaultSpeed)),
	}
	m.Noise.SetNoiseOctaves(2)
	m.Noise.SetNoisePersistence(0.23)
	m.Noise.SetNoiseLacunarity(0.3)
	m.Noise.SetNoiseScaleX(0.01)
	m.Noise.SetNoiseScaleY(0.01)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func GetDelay(level int) time.Duration {
	l := float64(speedLevels - level)
	return time.Duration(gaul.Map(1, speedLevels, minDelay, maxDelay, l)) * time.Millisecond
}

func (m Model) Init() tea.Cmd {
	return m.Timer.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeySpace:
			m.Grid.ToggleWrap()
		case tea.KeyRight:
			m.Grid.IncrementSeed()
		case tea.KeyLeft:
			m.Grid.DecrementSeed()
		case tea.KeyUp:
			if m.Speed < speedLevels {
				m.Speed++
				m.Timer.Interval = GetDelay(m.Speed)
			}
		case tea.KeyDown:
			if m.Speed > 1 {
				m.Speed--
				m.Timer.Interval = GetDelay(m.Speed)
			}
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "i":
				m.Grid.ToggleInvert()
			case "p":
				m.palette++
				if m.palette >= len(palettes) {
					m.palette = 0
				}
			case "q":
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.Rows = msg.Height
		m.Cols = msg.Width
		m.Grid = automata.NewGrid(m.Rows, m.Cols, m.Seed)
	}
	m.Grid.Update()
	m.Offset++
	m.Noise.SetNoiseOffsetY(m.Offset)
	m.Timer, cmd = m.Timer.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	rows := strings.Split(m.Grid.ToString(), "\n")
	n := len(rows)
	var s string
	for i, row := range rows {
		for j, c := range row {
			noise := m.Noise.Noise2D(float64(j), float64(i))
			color, _ := colorful.MakeColor(palettes[m.palette].Color(noise))
			s += lipgloss.NewStyle().Foreground(lipgloss.Color(color.Hex())).Render(string(c))
		}
		if i < n-1 {
			s += "\n"
		}
	}
	return s
}
