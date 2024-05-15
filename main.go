package main

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"math"
	"strconv"
)

type model struct {
	totalBudget     string
	downPayment     string
	interestRate    string
	monthlyPayments []string
	inputField      int
}

func initialModel() model {
	return model{
		totalBudget:     "",
		downPayment:     "",
		interestRate:    "",
		monthlyPayments: make([]string, 5),
		inputField:      0,
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return update(msg, m)
}

func (m model) View() string {
	return view(m)
}

func update(msg tea.Msg, m model) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.inputField > 0 {
				m.inputField--
			}
		case "down", "j":
			if m.inputField < 2 {
				m.inputField++
			}
		case "tab":
			if m.inputField < 2 {
				m.inputField++
			} else {
				m.inputField = 0
			}
		case "shift+tab":
			if m.inputField > 0 {
				m.inputField--
			} else {
				m.inputField = 2
			}
		case "enter":
			m.calculatePayments()
		case "backspace":
			switch m.inputField {
			case 0:
				if len(m.totalBudget) > 0 {
					m.totalBudget = m.totalBudget[:len(m.totalBudget)-1]
				}
			case 1:
				if len(m.downPayment) > 0 {
					m.downPayment = m.downPayment[:len(m.downPayment)-1]
				}
			case 2:
				if len(m.interestRate) > 0 {
					m.interestRate = m.interestRate[:len(m.interestRate)-1]
				}
			}
		default:
			switch m.inputField {
			case 0:
				m.totalBudget += msg.String()
			case 1:
				m.downPayment += msg.String()
			case 2:
				m.interestRate += msg.String()
			}
		}
	}

	return m, nil
}

func view(m model) string {
	return fmt.Sprintf(
		"Enter the total budget for the house (in USD): %s%s\n"+
			"Enter the percentage of the budget to be used as down payment: %s%s\n"+
			"Enter the current annual interest rate for a 30-year mortgage (in %%): %s%s\n\n"+
			"Estimated Monthly Payments:\n"+
			"Interest Rate | Monthly Payment\n"+
			"-----------------------------\n"+
			"%s\n%s\n%s\n%s\n%s\n",
		m.totalBudget, cursor(m.inputField == 0),
		m.downPayment, cursor(m.inputField == 1),
		m.interestRate, cursor(m.inputField == 2),
		m.monthlyPayments[0],
		m.monthlyPayments[1],
		m.monthlyPayments[2],
		m.monthlyPayments[3],
		m.monthlyPayments[4],
	)
}

func cursor(active bool) string {
	if active {
		return " █"
	}
	return ""
}

func (m *model) calculatePayments() {
	budget, err1 := strconv.ParseFloat(m.totalBudget, 64)
	downPaymentPercentage, err2 := strconv.ParseFloat(m.downPayment, 64)
	interestRate, err3 := strconv.ParseFloat(m.interestRate, 64)

	if err1 != nil || err2 != nil || err3 != nil || budget <= 0 || downPaymentPercentage < 0 || downPaymentPercentage > 100 || interestRate < 0 {
		for i := range m.monthlyPayments {
			m.monthlyPayments[i] = "Invalid input"
		}
		return
	}

	downPayment := budget * (downPaymentPercentage / 100)
	loanAmount := budget - downPayment

	rates := []float64{
		interestRate - 0.25,
		interestRate - 0.125,
		interestRate,
		interestRate + 0.125,
		interestRate + 0.25,
	}

	for i, rate := range rates {
		monthlyRate := rate / (12 * 100)
		term := 30 * 12
		monthlyPayment := loanAmount * (monthlyRate * math.Pow(1+monthlyRate, float64(term))) / (math.Pow(1+monthlyRate, float64(term)) - 1)
		if rate < 0 {
			m.monthlyPayments[i] = fmt.Sprintf("%.3f%% | Invalid rate", rate)
		} else {
			m.monthlyPayments[i] = fmt.Sprintf("%.3f%% | $%.2f + fees", rate, monthlyPayment)
		}
	}
}

