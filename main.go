package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/arr2036/yksofttoken/internal/token"
)

const appVersion = "1.0.0"

type ykSoftApp struct {
	app        fyne.App
	mainWindow fyne.Window
	token      *token.SoftToken
	tokenPath  string
	tokenDir   string

	// UI elements
	tokenSelect    *widget.Select
	otpDisplay     *widget.Entry
	regInfoDisplay *widget.Entry
	statusLabel    *widget.Label
	counterLabel   *widget.Label
	sessionLabel   *widget.Label
	generateBtn    *widget.Button
	copyBtn        *widget.Button
	copyRegBtn     *widget.Button
}

func main() {
	ykApp := &ykSoftApp{}
	ykApp.run()
}

func (y *ykSoftApp) run() {
	y.app = app.NewWithID("org.freeradius.yksoft")
	y.mainWindow = y.app.NewWindow("YKSoft Token")

	// Get default token directory
	var err error
	y.tokenDir, err = token.GetDefaultTokenDir()
	if err != nil {
		y.tokenDir = "."
	}

	// Ensure token directory exists
	os.MkdirAll(y.tokenDir, 0700)

	// Create UI
	y.createUI()

	// Load available tokens
	y.refreshTokenList()

	// Set window properties
	y.mainWindow.Resize(fyne.NewSize(500, 450))
	y.mainWindow.SetFixedSize(false)
	y.mainWindow.CenterOnScreen()
	y.mainWindow.ShowAndRun()
}

func (y *ykSoftApp) createUI() {
	// Token selection
	y.tokenSelect = widget.NewSelect([]string{}, y.onTokenSelected)
	y.tokenSelect.PlaceHolder = "Select or create a token..."

	newTokenBtn := widget.NewButtonWithIcon("New", theme.ContentAddIcon(), y.onNewToken)
	deleteTokenBtn := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), y.onDeleteToken)

	tokenRow := container.NewBorder(nil, nil, nil,
		container.NewHBox(newTokenBtn, deleteTokenBtn),
		y.tokenSelect,
	)

	// OTP display
	y.otpDisplay = widget.NewEntry()
	y.otpDisplay.SetPlaceHolder("OTP will appear here...")
	y.otpDisplay.Disable()

	y.generateBtn = widget.NewButtonWithIcon("Generate OTP", theme.MediaPlayIcon(), y.onGenerateOTP)
	y.generateBtn.Importance = widget.HighImportance
	y.generateBtn.Disable()

	y.copyBtn = widget.NewButtonWithIcon("Copy", theme.ContentCopyIcon(), y.onCopyOTP)
	y.copyBtn.Disable()

	otpButtons := container.NewHBox(y.generateBtn, y.copyBtn)

	// Registration info display
	y.regInfoDisplay = widget.NewMultiLineEntry()
	y.regInfoDisplay.SetPlaceHolder("Registration info will appear here...")
	y.regInfoDisplay.Disable()
	y.regInfoDisplay.Wrapping = fyne.TextWrapWord
	y.regInfoDisplay.SetMinRowsVisible(3)

	y.copyRegBtn = widget.NewButtonWithIcon("Copy Registration Info", theme.ContentCopyIcon(), y.onCopyRegInfo)
	y.copyRegBtn.Disable()

	// Status display
	y.statusLabel = widget.NewLabel("No token loaded")
	y.statusLabel.Alignment = fyne.TextAlignCenter

	y.counterLabel = widget.NewLabel("Counter: -")
	y.sessionLabel = widget.NewLabel("Session: -")

	statsRow := container.NewHBox(
		layout.NewSpacer(),
		y.counterLabel,
		widget.NewSeparator(),
		y.sessionLabel,
		layout.NewSpacer(),
	)

	// About/Help
	aboutBtn := widget.NewButtonWithIcon("About", theme.InfoIcon(), y.showAbout)

	// Layout
	content := container.NewVBox(
		widget.NewCard("Token", "", container.NewVBox(tokenRow)),
		widget.NewCard("One-Time Password", "", container.NewVBox(
			y.otpDisplay,
			otpButtons,
		)),
		widget.NewCard("Registration Information", "", container.NewVBox(
			y.regInfoDisplay,
			y.copyRegBtn,
		)),
		widget.NewCard("Status", "", container.NewVBox(
			y.statusLabel,
			statsRow,
		)),
		container.NewHBox(layout.NewSpacer(), aboutBtn),
	)

	scrollContent := container.NewVScroll(content)
	y.mainWindow.SetContent(scrollContent)
}

func (y *ykSoftApp) refreshTokenList() {
	tokens := []string{}

	entries, err := os.ReadDir(y.tokenDir)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
				tokens = append(tokens, entry.Name())
			}
		}
	}

	y.tokenSelect.Options = tokens
	if len(tokens) > 0 && y.tokenSelect.Selected == "" {
		y.tokenSelect.SetSelected(tokens[0])
	}
}

func (y *ykSoftApp) onTokenSelected(name string) {
	if name == "" {
		return
	}

	y.tokenPath = token.GetTokenPath(y.tokenDir, name)

	var err error
	y.token, err = token.Load(y.tokenPath)
	if err != nil {
		dialog.ShowError(fmt.Errorf("Failed to load token: %v", err), y.mainWindow)
		return
	}

	y.updateUI()
}

func (y *ykSoftApp) onNewToken() {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Token name (e.g., default)")

	dialog.ShowForm("New Token", "Create", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Name", entry),
		},
		func(confirmed bool) {
			if !confirmed || entry.Text == "" {
				return
			}

			name := strings.TrimSpace(entry.Text)
			path := token.GetTokenPath(y.tokenDir, name)

			// Check if token already exists
			if _, err := os.Stat(path); err == nil {
				dialog.ShowError(fmt.Errorf("Token '%s' already exists", name), y.mainWindow)
				return
			}

			// Create new token
			newToken, err := token.New()
			if err != nil {
				dialog.ShowError(fmt.Errorf("Failed to create token: %v", err), y.mainWindow)
				return
			}

			// Save token
			if err := newToken.Save(path); err != nil {
				dialog.ShowError(fmt.Errorf("Failed to save token: %v", err), y.mainWindow)
				return
			}

			y.token = newToken
			y.tokenPath = path
			y.refreshTokenList()
			y.tokenSelect.SetSelected(name)
			y.updateUI()

			// Show registration info for new token
			dialog.ShowInformation("Token Created",
				fmt.Sprintf("New token created!\n\nRegistration info:\n%s", newToken.RegistrationInfo()),
				y.mainWindow)
		},
		y.mainWindow,
	)
}

func (y *ykSoftApp) onDeleteToken() {
	if y.tokenPath == "" {
		return
	}

	name := filepath.Base(y.tokenPath)
	dialog.ShowConfirm("Delete Token",
		fmt.Sprintf("Are you sure you want to delete token '%s'?\n\nThis cannot be undone!", name),
		func(confirmed bool) {
			if !confirmed {
				return
			}

			if err := os.Remove(y.tokenPath); err != nil {
				dialog.ShowError(fmt.Errorf("Failed to delete token: %v", err), y.mainWindow)
				return
			}

			y.token = nil
			y.tokenPath = ""
			y.tokenSelect.ClearSelected()
			y.refreshTokenList()
			y.clearUI()
		},
		y.mainWindow,
	)
}

func (y *ykSoftApp) onGenerateOTP() {
	if y.token == nil {
		return
	}

	otp, err := y.token.GenerateOTP()
	if err != nil {
		dialog.ShowError(fmt.Errorf("Failed to generate OTP: %v", err), y.mainWindow)
		return
	}

	// Save updated token state
	if err := y.token.Save(y.tokenPath); err != nil {
		dialog.ShowError(fmt.Errorf("Failed to save token state: %v", err), y.mainWindow)
		return
	}

	y.otpDisplay.SetText(otp)
	y.copyBtn.Enable()
	y.updateUI()
}

func (y *ykSoftApp) onCopyOTP() {
	if y.otpDisplay.Text != "" {
		y.mainWindow.Clipboard().SetContent(y.otpDisplay.Text)
		y.statusLabel.SetText("OTP copied to clipboard!")
		go func() {
			time.Sleep(2 * time.Second)
			y.statusLabel.SetText("Ready")
		}()
	}
}

func (y *ykSoftApp) onCopyRegInfo() {
	if y.token != nil {
		y.mainWindow.Clipboard().SetContent(y.token.RegistrationInfo())
		y.statusLabel.SetText("Registration info copied to clipboard!")
		go func() {
			time.Sleep(2 * time.Second)
			y.statusLabel.SetText("Ready")
		}()
	}
}

func (y *ykSoftApp) updateUI() {
	if y.token == nil {
		y.clearUI()
		return
	}

	y.generateBtn.Enable()
	y.copyRegBtn.Enable()
	y.regInfoDisplay.SetText(y.token.RegistrationInfo())
	y.counterLabel.SetText(fmt.Sprintf("Counter: %d", y.token.Counter))
	y.sessionLabel.SetText(fmt.Sprintf("Session: %d", y.token.Session))
	y.statusLabel.SetText("Ready")
}

func (y *ykSoftApp) clearUI() {
	y.generateBtn.Disable()
	y.copyBtn.Disable()
	y.copyRegBtn.Disable()
	y.otpDisplay.SetText("")
	y.regInfoDisplay.SetText("")
	y.counterLabel.SetText("Counter: -")
	y.sessionLabel.SetText("Session: -")
	y.statusLabel.SetText("No token loaded")
}

func (y *ykSoftApp) showAbout() {
	dialog.ShowInformation("About YKSoft Token",
		fmt.Sprintf("YKSoft Token v%s\n\n"+
			"A software Yubikey token emulator.\n\n"+
			"Generates HOTP One Time Passcodes in Yubikey format.\n\n"+
			"Useful for testing and M2M VPN connections\n"+
			"that require 2FA.\n\n"+
			"© 2022-2024 Arran Cudbard-Bell",
			appVersion),
		y.mainWindow)
}
