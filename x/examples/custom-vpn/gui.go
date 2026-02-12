package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type drFrakeTheme struct {
	fyne.Theme
}

func (t *drFrakeTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 14, G: 30, B: 45, A: 255} // Dark primary from logo
	case theme.ColorNameButton:
		return color.RGBA{R: 0, G: 181, B: 226, A: 255} // Cyan from logo
	case theme.ColorNamePrimary:
		return color.RGBA{R: 0, G: 215, B: 255, A: 255} // Glowing Cyan
	case theme.ColorNameForeground:
		return color.White
	default:
		return t.Theme.Color(name, variant)
	}
}

func (t *drFrakeTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return 14
	}
	return t.Theme.Size(name)
}

func setupGUI(myApp fyne.App) fyne.Window {
	myApp.Settings().SetTheme(&drFrakeTheme{Theme: theme.DefaultTheme()})
	win := myApp.NewWindow("Dr. Frake VPN")
	win.Resize(fyne.NewSize(400, 600))

	// UI Elements
	logoLabel := widget.NewLabelWithStyle("Dr. Frake", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	logoLabel.Alignment = fyne.TextAlignCenter

	title := canvas.NewText("AI-ASSISTANT VPN", color.RGBA{R: 0, G: 215, B: 255, A: 255})
	title.TextSize = 20
	title.Alignment = fyne.TextAlignCenter

	statusLabel := widget.NewLabel("Status: Disconnected")
	statusLabel.Alignment = fyne.TextAlignCenter

	transportEntry := widget.NewMultiLineEntry()
	transportEntry.SetPlaceHolder("Enter Shadowsocks (ss://) Key")
	transportEntry.Wrapping = fyne.TextWrapBreak

	var isConnected bool
	connectBtn := widget.NewButton("CONNECT", nil)
	connectBtn.Importance = widget.HighImportance

	connectBtn.OnTapped = func() {
		if !isConnected {
			statusLabel.SetText("Status: Connecting...")
			connectBtn.Disable()
			err := startVPN(transportEntry.Text)
			if err != nil {
				statusLabel.SetText("Error: " + err.Error())
				connectBtn.Enable()
				return
			}
			isConnected = true
			statusLabel.SetText("Status: Connected")
			connectBtn.SetText("DISCONNECT")
			connectBtn.Enable()
		} else {
			statusLabel.SetText("Status: Disconnecting...")
			connectBtn.Disable()
			stopVPN()
			isConnected = false
			statusLabel.SetText("Status: Disconnected")
			connectBtn.SetText("CONNECT")
			connectBtn.Enable()
		}
	}

	content := container.NewVBox(
		layout.NewSpacer(),
		logoLabel,
		title,
		layout.NewSpacer(),
		container.NewPadded(transportEntry),
		layout.NewSpacer(),
		statusLabel,
		container.NewPadded(connectBtn),
		layout.NewSpacer(),
	)

	win.SetContent(container.NewMax(
		canvas.NewRectangle(color.RGBA{R: 10, G: 20, B: 30, A: 255}), // Deeper background
		content,
	))

	return win
}
