package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"gioui.org/app"
	"git.sr.ht/~jackmordaunt/go-nativenotify"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	var (
		title   string
		body    string
		icon    string
		payload string
	)

	flag.StringVar(&title, "title", "title", "title text of notification")
	flag.StringVar(&body, "body", "body", "body text of notification")
	flag.StringVar(&icon, "icon", "./cmd/notify/puzzle.png", "path to icon file for notification")
	flag.StringVar(&payload, "payload", "example-payload", "data to be returned upon activation")
	flag.Parse()

	if abs, err := filepath.Abs(icon); err == nil {
		icon = abs
	}

	slog.Info("icon", "path", icon)

	if err := nativenotify.Setup(nativenotify.Config{
		Windows: nativenotify.WindowsConfig{
			AppID:    "notify-test",
			GUID:     "{B5E38D62-B912-486C-96E3-6FAD1082B73D}",
			IconPath: icon,
		},
		Linux: nativenotify.LinuxConfig{
			AppName: "notify-test",
			AppIcon: icon,
		},
		Darwin: nativenotify.DarwinConfig{},
	}); err != nil {
		slog.Error("setting up notification subsystem", "err", err)
	}

	if err := nativenotify.Push(nativenotify.Notification{
		ID:         "text-message",
		Title:      title,
		Body:       body,
		Icon:       icon,
		AppPayload: payload,
		TextActions: []nativenotify.TextAction{
			{
				ID:              "reply",
				Title:           "Reply",
				PlaceholderHint: "type here...",
				ButtonLabel:     "Send",
			},
		},
		ButtonActions: []nativenotify.ButtonAction{
			{
				ID:         "like",
				LabelText:  "Like",
				AppPayload: "@jack",
			},
		},
		Callback: func(err error, id string, userData map[string]string) {
			slog.Info("callback", "id", id, "userData", userData)
			if err != nil {
				slog.Error("callback error", "err", err)
			}
		},
	}); err != nil {
		slog.Error("pushing notification", "err", err)
	}

	fmt.Printf("callback invocation will log, quit when done (ctr-c)\n")
	wait()
}

func wait() {
	switch runtime.GOOS {
	case "windows":
		// Avoid the deadlock panic on Windows.
		<-time.NewTimer(1<<63 - 1).C
	case "darwin":
		app.Main()
	default:
		select {}
	}
}
