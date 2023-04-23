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
	grid     *automata.Grid
	noise    gaul.Rng
	offset   float64
	timer    timer.Model
	seed     int64
	speed    int
	rows     int
	cols     int
	wrapEnds bool
}

func StartTea(seed int64) {
	m := Model{
		grid:  automata.NewGrid(1, 1, seed),
		noise: gaul.NewRng(seed),
		seed:  seed,
		speed: defaultSpeed,
		timer: timer.NewWithInterval(timeout, GetDelay(defaultSpeed)),
	}
	m.noise.SetNoiseOctaves(2)
	m.noise.SetNoisePersistence(0.23)
	m.noise.SetNoiseLacunarity(0.3)
	m.noise.SetNoiseScaleX(0.01)
	m.noise.SetNoiseScaleY(0.01)
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
	return m.timer.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeySpace:
			m.grid.ToggleWrap()
		case tea.KeyRight:
			m.grid.IncrementSeed()
		case tea.KeyLeft:
			m.grid.DecrementSeed()
		case tea.KeyUp:
			if m.speed < speedLevels {
				m.speed++
				m.timer.Interval = GetDelay(m.speed)
			}
		case tea.KeyDown:
			if m.speed > 1 {
				m.speed--
				m.timer.Interval = GetDelay(m.speed)
			}
		}
	case tea.WindowSizeMsg:
		m.rows = msg.Height
		m.cols = msg.Width
		m.grid = automata.NewGrid(m.rows, m.cols, m.seed)
	}
	m.grid.Update()
	m.offset++
	m.noise.SetNoiseOffsetY(m.offset)
	m.timer, cmd = m.timer.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	rows := strings.Split(m.grid.ToString(), "\n")
	n := len(rows)
	var s string
	for i, row := range rows {
		for j, c := range row {
			noise := m.noise.Noise2D(float64(j), float64(i))
			color, _ := colorful.MakeColor(synthWaveGradientFull.Color(noise))
			s += lipgloss.NewStyle().Foreground(lipgloss.Color(color.Hex())).Render(string(c))
		}
		if i < n-1 {
			s += "\n"
		}
	}
	return s
}
