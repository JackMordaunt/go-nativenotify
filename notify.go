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
	darwinnotify "git.sr.ht/~jackmordaunt/go-notify-darwin"
	windowsnotify "git.sr.ht/~jackmordaunt/go-toast/v2"
)

// Callback is executed when the user interacts with a given notification.
//
// [userData] contains user input data, which is either text input, the
// app payload for a button, or the app payload from the parent notification.
type Callback func(err error, id string, userData map[string]string)

// Notification describes the notification.
type Notification struct {
	ID string

	// Title text of the notification.
	Title string

	// Body text of the notification.
	Body string

	// Icon contains the path to an icon. Some operating systems accept url paths.
	// All OS accept .png, other formats vary.
	Icon string

	// AppPayload is passed to the callback upon activation.
	AppPayload string

	// TextActions are for getting textual input from the user.
	// Not supported on Linux.
	TextActions []TextAction

	// ButtonActions are for getting button input from the user.
	ButtonActions []ButtonAction

	// Callback is called upon activation.
	Callback Callback
}

// TextAction describes a text input.
// It might have an associated button.
// The user input text is passed to the callback on activation.
type TextAction struct {
	ID              string
	Title           string
	PlaceholderHint string
	ButtonLabel     string
}

// ButtonAction describes a button input.
// [AppPayload] is passed to the callback on activation.
type ButtonAction struct {
	ID         string
	LabelText  string
	AppPayload string
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
	return push(n)
}
