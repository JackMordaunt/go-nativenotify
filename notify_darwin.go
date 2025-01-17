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

		data := map[string]string{}

		if args.UserText != "" {
			data[actionID] = args.UserText
		} else if actionArgs != "" {
			data[actionID] = actionArgs
		}

		// Map the user data map to the key-value slice style.
		for k, v := range args.UserData {
			if k == "id" {
				continue
			}
			data[k] = v
		}

		fn, ok := callbacksTake(&callbacks, id)
		if !ok || fn == nil {
			return
		}

		fn(args.Err, actionID, data)
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

	userData["default"] = n.AppPayload
	userData["id"] = strconv.FormatInt(id, 10)

	for _, button := range n.ButtonActions {
		buttons = append(buttons, darwinnotify.Action{
			ID:    fmt.Sprintf("%s-%s", encode(button.ID), encode(button.AppPayload)),
			Title: button.LabelText,
		})
	}

	for _, input := range n.TextActions {
		inputs = append(inputs, darwinnotify.TextInputAction{
			ID:          fmt.Sprintf("%d-%s", id, encode(input.ID)),
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

	callbacksPut(&callbacks, strconv.FormatInt(id, 10), func(err error, id string, data map[string]string) {
		n.Callback(err, n.ID, data)
	})

	return nil
}
