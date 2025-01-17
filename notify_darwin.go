package nativenotify

import (
	"fmt"
	"strconv"
	"strings"

	darwinnotify "git.sr.ht/~jackmordaunt/go-notify-darwin"
)

const defaultAction = "com.apple.UNNotificationDefaultActionIdentifier"

func setup(cfg Config) error {
	darwinnotify.Init(cfg.Darwin.Categories...)

	darwinnotify.SetCallback(func(args darwinnotify.CallbackArgs) {
		id := args.UserData["id"]

		parts := strings.Split(args.Action, "-")

		actionIDEncoded, _ := take(&parts)
		actionArgsEncoded, _ := take(&parts)

		actionID := decode(actionIDEncoded)
		actionArgs := decode(actionArgsEncoded)

		if args.Action == defaultAction {
			actionID = "default"
		}

		if args.UserText != "" {
			actionArgs = args.UserText
		}

		fn, ok := callbacksTake(&callbacks, id)
		if !ok || fn == nil {
			return
		}

		fn(actionID, actionArgs)
	})

	return nil
}

func push(n Notification) (err error) {
	id := nextID.Add(1)

	var (
		buttons  = make([]darwinnotify.Action, 0, len(n.ButtonActions))
		inputs   = make([]darwinnotify.TextInputAction, 0, len(n.TextActions))
		userData = make(darwinnotify.UserData)
	)

	userData["id"] = strconv.FormatInt(id, 10)

	for _, button := range n.ButtonActions {
		buttons = append(buttons, darwinnotify.Action{
			ID:    fmt.Sprintf("%s-%s", encode(button.ID), encode(button.Value)),
			Title: button.LabelText,
		})
	}

	for _, input := range n.TextActions {
		inputs = append(inputs, darwinnotify.TextInputAction{
			ID:          encode(input.ID),
			Title:       input.Title,
			Placeholder: input.PlaceholderHint,
			ButtonTitle: input.ButtonLabel,
		})
	}

	var attachments []string

	if n.Icon != "" {
		attachments = []string{n.Icon}
	}

	darwinnotify.Notify(darwinnotify.Notification{
		Title:            n.Title,
		Body:             n.Body,
		Attachments:      attachments,
		Actions:          buttons,
		TextInputActions: inputs,
		UserData:         userData,
	})

	callbacksPut(&callbacks, strconv.FormatInt(id, 10), n.Callback)

	return nil
}
