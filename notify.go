// Package nativenotify offers a notification API that works across Windows,
// macOS and Linux.
//
// It provides a common subset, and is not trying to expose all underlying
// features. Not all fields are functional across all operating systems.
//
// Broadly, Windows likes to define buttons and text inputs separately, macOs
// defines them together, and Linux only supports buttons.
//
// Each operating system requires different setup data. This is handled by
// a call to [Setup] that accepts a fat-union containing all the data
// required by all operating systems.
//
// Generally each operating system will accept png icons as file paths, but
// other formats (eg webp) and path types (eg url) vary.
package nativenotify

import (
	"fmt"

	darwinnotify "git.sr.ht/~jackmordaunt/go-notify-darwin"
	windowsnotify "git.sr.ht/~jackmordaunt/go-toast/v2"
)

// Callback is executed when the user interacts with a given notification.
// [action] is the activated action.
// [value] is any associated value for that action.
type Callback func(action, value string)

// Notification describes the notification.
type Notification struct {
	// Callback is called upon activation.
	Callback Callback

	// Title text of the notification.
	Title string

	// Body text of the notification.
	Body string

	// Icon contains the path to an icon. Some operating systems accept url paths.
	// All OS accept .png, other formats vary.
	Icon string

	// TextActions are for getting textual input from the user.
	// Not supported on Linux.
	TextActions []TextAction

	// ButtonActions are for getting button input from the user.
	ButtonActions []ButtonAction

	// Windows defines windows specific options.
	Windows WindowsOption
}

type WindowsOption int

const (
	WindowsOptionIconSquareCrop WindowsOption = iota
	WindowsOptionIconCircleCrop
	WindowsOptionIconHero
)

// TextAction describes a text input.
// It might have an associated button.
// The user input text is passed to the callback on activation.
// Not all OS support this.
type TextAction struct {
	// ID names the action. The user text will appear in the user data keyed by this ID.
	ID string
	// Title describes the title text of this action.
	Title string
	// PlaceholderHint will appear as the placeholder text within the text input.
	// Not all OS support this.
	PlaceholderHint string
	// ButtonLabel describes the content of any related button that is associated with
	// the text input. Not all OS support this.
	ButtonLabel string
}

// ButtonAction describes a button input.
// [AppPayload] is passed to the callback on activation.
type ButtonAction struct {
	// ID names the action. The [AppPayload] will appear in the user data keyed by this ID.
	ID string
	// Value is provided to the callback if this action is activated.
	Value string
	// LabelText describes the text content of this button action.
	LabelText string
}

// Config is a fat-union of the various initialization data required by each
// operating system.
type Config struct {
	Windows WindowsConfig
	Linux   LinuxConfig
	Darwin  DarwinConfig
}

type WindowsConfig = windowsnotify.AppData

type LinuxConfig struct {
	AppName string
	AppIcon string
}

type DarwinConfig struct {
	Categories []darwinnotify.Category
}

// Setup initializes the notification subsystem.
func Setup(cfg Config) error {
	return setup(cfg)
}

// Push a notification to the operating system.
func Push(n Notification) error {
	for ii, a := range n.ButtonActions {
		if a.ID == "" {
			return fmt.Errorf("buttonaction %d requires ID", ii)
		}
	}
	for ii, a := range n.TextActions {
		if a.ID == "" {
			return fmt.Errorf("text action %d requires ID", ii)
		}
	}
	return push(n)
}
