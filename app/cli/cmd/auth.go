package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/do"
	"gssm/config"
	"gssm/db"
	"gssm/utils"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("43"))
	khakiStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("228"))
	forgot       = false
)

type authModel struct {
	inj *do.Injector

	masterKeyInput textinput.Model

	masterKey string
	email     string
	password  string
}

func initialModel(inj *do.Injector) authModel {
	ti := textinput.New()
	ti.Placeholder = "Master key"
	ti.Focus()
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'
	ti.PromptStyle = focusedStyle
	ti.TextStyle = focusedStyle

	return authModel{
		masterKeyInput: ti,
		inj:            inj,
	}
}

func (m authModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m authModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, m.handleSubmit()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	m.masterKeyInput, cmd = m.masterKeyInput.Update(msg)

	return m, cmd
}

func (m authModel) View() string {
	var b strings.Builder

	b.WriteRune('\n')
	b.WriteString("Type master key (used only in web-app interface)")
	b.WriteRune('\n')
	b.WriteString(m.masterKeyInput.View())

	return b.String()
}

func (m authModel) handleSubmit() tea.Cmd {
	m.masterKey = m.masterKeyInput.Value()

	cliConfig := config.ReadCliConfigToStruct()

	hashResult := utils.HashPassword(m.masterKey)

	cliConfig.MasterKeyHash = hex.EncodeToString(hashResult.Hash)
	cliConfig.MasterKeySalt = hex.EncodeToString(hashResult.Salt)

	err := config.DumpCliConfig(&cliConfig)
	if err != nil {
		fmt.Println(err.Error())
		return tea.Quit
	}

	return tea.Quit
}

func init() {
	var authCmd = &cobra.Command{
		Use:   "auth",
		Short: "Auth!!",
		PreRun: func(cmd *cobra.Command, args []string) {
			cliConfig := config.ReadCliConfigToStruct()

			if !forgot && cliConfig.MasterKeyHash != "" && cliConfig.MasterKeySalt != "" {
				fmt.Println(khakiStyle.Render("Master key is already set, exiting...\n"))

				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			config.ReadConfigFromHomeDirToViper()
			config.ReadCliConfigFromHomeDirToViper()
			inj := do.New()
			do.Provide(inj, db.NewDatabase)
			do.Provide(inj, db.NewUserSource)
			p := tea.NewProgram(initialModel(inj))
			if _, err := p.Run(); err != nil {
				fmt.Printf("error %s", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.AddCommand(authCmd)
	authCmd.Flags().BoolVarP(&forgot, "forgot", "f", false, "use this flag if you have forgotten master key and need to reset it")
}
