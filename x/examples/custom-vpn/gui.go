package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// App State
var (
	currentUser  UserInfo
	allServers   []Server
	activeServer *Server
	isConnected  bool
	statusLabel  *widget.Label
	connectBtn   *widget.Button
	contentArea  *fyne.Container
)

type drFrakeTheme struct {
	fyne.Theme
}

func (t *drFrakeTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 10, G: 15, B: 25, A: 255} // Midnight Navy
	case theme.ColorNameButton:
		return color.RGBA{R: 0, G: 161, B: 201, A: 255} // Professional Cyan
	case theme.ColorNamePrimary:
		return color.RGBA{R: 0, G: 215, B: 255, A: 255} // Electronic Blue
	case theme.ColorNameForeground:
		return color.White
	case theme.ColorNameOverlayBackground:
		return color.RGBA{R: 20, G: 25, B: 35, A: 230} // Semi-transparent sidebars
	default:
		return t.Theme.Color(name, variant)
	}
}

func setupGUI(myApp fyne.App) fyne.Window {
	myApp.Settings().SetTheme(&drFrakeTheme{Theme: theme.DefaultTheme()})
	win := myApp.NewWindow("Dr. Frake VPN - Business Edition")
	win.Resize(fyne.NewSize(800, 600))

	// Initial Data Load
	currentUser = GetUserInfo()
	allServers = FetchServerList()

	// Sidebar
	sidebar := createSidebar()
	contentArea = container.NewMax()

	// Default View
	showHomeView()

	mainLayout := container.NewHSplit(sidebar, contentArea)
	mainLayout.Offset = 0.2

	win.SetContent(container.NewMax(
		canvas.NewRectangle(color.RGBA{R: 5, G: 10, B: 20, A: 255}),
		mainLayout,
	))

	return win
}

func createSidebar() fyne.CanvasObject {
	logo := widget.NewLabelWithStyle("DR. FRAKE", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	homeBtn := widget.NewButtonWithIcon("Home", theme.HomeIcon(), showHomeView)
	locBtn := widget.NewButtonWithIcon("Locations", theme.NavigateNextIcon(), showLocationsView)
	priceBtn := widget.NewButtonWithIcon("Pricing", theme.SettingsIcon(), showPricingView)

	homeBtn.Alignment = widget.ButtonAlignLeading
	locBtn.Alignment = widget.ButtonAlignLeading
	priceBtn.Alignment = widget.ButtonAlignLeading

	avatar := widget.NewLabelWithStyle(currentUser.Email, fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
	planBadge := widget.NewLabelWithStyle(string(currentUser.Plan), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	return container.NewVBox(
		layout.NewSpacer(),
		logo,
		layout.NewSpacer(),
		homeBtn,
		locBtn,
		priceBtn,
		layout.NewSpacer(),
		container.NewVBox(avatar, planBadge),
		layout.NewSpacer(),
	)
}

func showHomeView() {
	title := canvas.NewText("SECURE CONNECTION", color.White)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	statusLabel = widget.NewLabel("Status: Disconnected")
	statusLabel.Alignment = fyne.TextAlignCenter

	serverLabel := widget.NewLabel("Selected: None")
	serverLabel.Alignment = fyne.TextAlignCenter
	if activeServer != nil {
		serverLabel.SetText(fmt.Sprintf("Selected: %s %s", activeServer.Flag, activeServer.Country))
	}

	connectBtn = widget.NewButton("CONNECT", nil)
	connectBtn.Importance = widget.HighImportance
	connectBtn.OnTapped = handleConnectToggle

	updateHomeUI()

	view := container.NewCenter(
		container.NewVBox(
			title,
			layout.NewSpacer(),
			serverLabel,
			statusLabel,
			layout.NewSpacer(),
			container.NewPadded(connectBtn),
		),
	)
	contentArea.Objects = []fyne.CanvasObject{view}
	contentArea.Refresh()
}

func showLocationsView() {
	list := widget.NewList(
		func() int { return len(allServers) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Flags"),
				widget.NewLabel("Country"),
				layout.NewSpacer(),
				widget.NewLabel("Latency"),
				widget.NewButton("Select", nil),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			s := allServers[id]
			hbox := item.(*fyne.Container)
			hbox.Objects[0].(*widget.Label).SetText(s.Flag)
			hbox.Objects[1].(*widget.Label).SetText(s.Country)
			hbox.Objects[3].(*widget.Label).SetText(fmt.Sprintf("%d ms", s.Latency))

			btn := hbox.Objects[4].(*widget.Button)
			if s.IsPremium && currentUser.Plan != PlanPremium {
				btn.SetText("PREMIUM")
				btn.OnTapped = showPricingView
			} else {
				btn.SetText("SELECT")
				btn.OnTapped = func() {
					activeServer = &allServers[id]
					showHomeView()
				}
			}
		},
	)

	view := container.NewBorder(
		widget.NewLabelWithStyle("GLOBAL SERVER LOCATIONS", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil, nil, nil,
		list,
	)
	contentArea.Objects = []fyne.CanvasObject{view}
	contentArea.Refresh()
}

func showPricingView() {
	title := widget.NewLabelWithStyle("CHOOSE YOUR PLAN", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	freeCard := container.NewVBox(
		widget.NewLabelWithStyle("FREE", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Basic Speed"),
		widget.NewLabel("2 Locations"),
		widget.NewButton("Current", nil),
	)

	premiumCard := container.NewBorder(
		nil,
		widget.NewButton("UPGRADE NOW", func() {
			currentUser.Plan = PlanPremium
			showHomeView()
		}),
		nil, nil,
		container.NewVBox(
			widget.NewLabelWithStyle("PREMIUM", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewLabel("Ultra High Speed"),
			widget.NewLabel("Global 10+ Locations"),
			widget.NewLabel("Dedicated Support"),
		),
	)

	view := container.NewCenter(
		container.NewVBox(
			title,
			layout.NewSpacer(),
			container.NewHBox(
				container.NewPadded(freeCard),
				container.NewPadded(premiumCard),
			),
		),
	)
	contentArea.Objects = []fyne.CanvasObject{view}
	contentArea.Refresh()
}

func handleConnectToggle() {
	if activeServer == nil {
		statusLabel.SetText("Please select a location first")
		return
	}

	if !isConnected {
		statusLabel.SetText("Connecting to " + activeServer.Country + "...")
		connectBtn.Disable()
		go func() {
			err := startVPN(activeServer.Config)
			if err != nil {
				isConnected = false
				statusLabel.SetText("Cloud Error: " + err.Error())
				connectBtn.Enable()
				return
			}
			isConnected = true
			updateHomeUI()
		}()
	} else {
		statusLabel.SetText("Disconnecting...")
		connectBtn.Disable()
		stopVPN()
		isConnected = false
		updateHomeUI()
	}
}

func updateHomeUI() {
	if isConnected {
		statusLabel.SetText("CONNECTED")
		connectBtn.SetText("DISCONNECT")
		connectBtn.Importance = widget.WarningImportance
	} else {
		statusLabel.SetText("DISCONNECTED")
		connectBtn.SetText("CONNECT")
		connectBtn.Importance = widget.HighImportance
	}
	connectBtn.Enable()
	if connectBtn.OnTapped == nil {
		connectBtn.OnTapped = handleConnectToggle
	}
}
