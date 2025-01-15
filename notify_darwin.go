package nativenotify

import (
	"fmt"
	"strconv"
	"strings"

	darwinnotify "git.sr.ht/~jackmordaunt/go-notify-darwin"
)

func setup(cfg Config) error {
	darwinnotify.Init(cfg.Darwin.Categories...)

	darwinnotify.SetCallback(func(args darwinnotify.CallbackArgs) {
		// Extract the callback id that is prepended to the action id.
		id, actionName, _ := strings.Cut(args.Action, "-")

		fn, ok := callbacksTake(&callbacks, id)
		if !ok || fn == nil {
			return
		}

		data := map[string]string{}

		if args.UserText != "" {
			data[actionName] = args.UserText
		}

		// Map the user data map to the key-value slice style.
		for k, v := range args.UserData {
			data[k] = v
		}

		data["category"] = args.Category
		data["action"] = args.Action

		fn(args.Err, id, data)
	})

	return nil
}

func push(n Notification) (err error) {
	id := nextID.Add(1)

	defer func() {
		if err == nil {
			callbacksPut(&callbacks, strconv.FormatInt(id, 10), n.Callback)
		}
	}()

	var (
		buttons = make([]darwinnotify.Action, 0, len(n.ButtonActions))
		inputs  = make([]darwinnotify.TextInputAction, 0, len(n.TextActions))
	)

	for _, button := range n.ButtonActions {
		buttons = append(buttons, darwinnotify.Action{
			ID:    fmt.Sprintf("%d-%s", id, button.ID),
			Title: button.LabelText,
		})
	}

	for _, input := range n.TextActions {
		inputs = append(inputs, darwinnotify.TextInputAction{
			ID:          fmt.Sprintf("%d-%s", id, input.ID),
			Title:       input.Title,
			Placeholder: input.PlaceholderHint,
			ButtonTitle: input.ButtonLabel,
		})
	}

	darwinnotify.Notify(darwinnotify.Notification{
		Title:            n.Title,
		Body:             n.Body,
		Attachments:      []string{n.Icon},
		Actions:          buttons,
		TextInputActions: inputs,
	})

	return nil
}
