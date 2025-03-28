package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/cpainter1/PassLock/internal/database"
	"github.com/cpainter1/PassLock/internal/encryption"
	"log"
)

// ShowCreateVaultForm displays a form to create a new vault.
func ShowCreateVaultForm(win fyne.Window) {
	win.SetTitle("Create New Vault")
	win.Resize(fyne.NewSize(300, 250))

	// Entry fields for vault name and password
	vaultNameEntry := widget.NewEntry()
	vaultNameEntry.SetPlaceHolder("Enter vault name")

	vaultPasswordEntry := widget.NewPasswordEntry()
	vaultPasswordEntry.SetPlaceHolder("Enter PRIVATE master password")

	masterPasswordNote := widget.NewLabelWithStyle(
		"NOTE: This master password will be required to access your vault. DO NOT SHARE IT.",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true})
	masterPasswordNote.Wrapping = fyne.TextWrapWord // Wrapping

	// Create button
	createButton := widget.NewButtonWithIcon("Create", theme.ConfirmIcon(), func() {
		vaultName := vaultNameEntry.Text
		vaultPassword := vaultPasswordEntry.Text

		if vaultName == "" || vaultPassword == "" {
			log.Println("Vault name or password cannot be empty")
			return
		}

		keySalt, err := encryption.GenerateSalt(16)
		if err != nil {
			log.Printf("Error creating salt: %s", err)
			return
		}

		_, authKey, err := encryption.DeriveMasterKeys(vaultPassword, keySalt)
		log.Printf("Derived master key: %v", authKey)
		if err != nil {
			log.Fatalf("Failed to derive master keys: %v", err)
			return
		}

		err = database.CreateVault(vaultName, authKey, keySalt)
		if err != nil {
			log.Println("Failed to create vault:", err)
			return
		}

		// After creating, go back to vault selection UI
		ShowLoginUI(win)
	})
	createButton.Importance = widget.HighImportance

	// Cancel button to go back
	cancelButton := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		ShowLoginUI(win)
	})
	cancelButton.Importance = widget.DangerImportance

	// Layout
	form := container.NewVBox(
		widget.NewLabelWithStyle("Create a New Vault", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		vaultNameEntry,
		vaultPasswordEntry,
		masterPasswordNote,
		createButton,
		cancelButton,
	)

	win.SetContent(container.NewPadded(form))
}

// ShowAuthenticationForm displays the authentication form
func ShowAuthenticationForm(win fyne.Window, vaultName string) {
	win.SetTitle("Vault Authentication - " + vaultName)
	win.Resize(fyne.NewSize(300, 250))

	// Top labels
	infoLabel := widget.NewLabelWithStyle(
		"Enter your master password below",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true})

	// Entry field for master password
	vaultPasswordEntry := widget.NewPasswordEntry()
	vaultPasswordEntry.SetPlaceHolder("Enter PRIVATE master password")

	// Label for error message
	authResultLabel := widget.NewLabel("")

	// Authenticate button
	authenticateButton := widget.NewButtonWithIcon("Authenticate", theme.LoginIcon(), func() {
		// Generate salt
		salt, err := database.GetSaltFromVault(vaultName)
		if err != nil {
			log.Printf("Error retrieving salt: %s", err)
			authResultLabel.SetText("Failed to authenticate")
			return
		}

		// Obtain authKey from provided master password
		_, authKey, err := encryption.DeriveMasterKeys(vaultPasswordEntry.Text, salt)
		if err != nil {
			log.Printf("Failed to derive master key: %v", err)
		}

		authenticationResult, err := database.AuthenticateVault(vaultName, authKey)
		if err != nil {
			log.Printf("Failed to authenticate: %s", err)
			authResultLabel.SetText("Failed to authenticate")
			return
		}

		// Check if authenticated
		if !authenticationResult {
			log.Print(vaultPasswordEntry.Text)
			authResultLabel.SetText("Authentication denied. Please try again.")
			return
		} else {
			authResultLabel.SetText("Authentication succeeded")
			// TODO: Implement main view
		}
	})
	authenticateButton.Importance = widget.HighImportance

	// Back button
	backButton := widget.NewButtonWithIcon("Back", theme.CancelIcon(), func() {
		ShowLoginUI(win)
	})
	backButton.Importance = widget.DangerImportance

	// Layout
	form := container.NewVBox(
		infoLabel,
		vaultPasswordEntry,
		authResultLabel,
		authenticateButton,
		backButton,
	)

	win.SetContent(container.NewPadded(form))
}

// ShowLoginUI displays the main selection view.
func ShowLoginUI(win fyne.Window) {
	// Window setup
	win.SetTitle("Vault Selection")
	win.Resize(fyne.NewSize(300, 450))
	win.SetFixedSize(true)

	// Create UI elements
	title := widget.NewLabelWithStyle(
		"Welcome to PassLock!",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true})
	description := widget.NewLabelWithStyle(
		"Select a vault or create a new one",
		fyne.TextAlignCenter,
		fyne.TextStyle{})

	// List vaults
	vaults, err := database.ListVaults()
	if err != nil {
		log.Println("Unable to list vaults:", err)
	}

	// Create vault list
	vaultsList := widget.NewList(
		// List length
		func() int {
			return len(vaults)
		},
		// Create item
		func() fyne.CanvasObject {
			// Return new button for each vault
			return widget.NewButton("", nil)
		},
		// Update item
		func(i int, o fyne.CanvasObject) {
			button := o.(*widget.Button)
			button.SetText(vaults[i])
			button.OnTapped = func() {
				ShowAuthenticationForm(win, vaults[i])
			}
		})

	// Button to create a new vault
	createVaultButton := widget.NewButtonWithIcon(
		"Create New Vault",
		theme.ContentAddIcon(),
		func() {
			ShowCreateVaultForm(win)
		})
	createVaultButton.Importance = widget.HighImportance

	// Create "Vaults:" label
	vaultsLabel := widget.NewLabelWithStyle(
		"Vaults:",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true})

	// Layout
	topContent := container.NewVBox(
		title,
		description,
		createVaultButton,
		vaultsLabel,
	)
	// Make a border to contain all objects
	fullContent := container.NewBorder(
		topContent,
		nil, // No bottom needed
		nil, // No left needed
		nil, // No right needed
		container.NewPadded(vaultsList),
	)

	win.SetContent(fullContent)
}
