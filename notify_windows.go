package nativenotify

import (
	"strconv"
	"strings"

	windowsnotify "git.sr.ht/~jackmordaunt/go-toast/v2"
)

func setup(cfg Config) error {
	windowsnotify.SetAppData(cfg.Windows)
	windowsnotify.SetActivationCallback(func(args string, userdata []windowsnotify.UserData) {
		// Address the fact that the ID might be prepended to some argument
		// and extract it.
		id, args, _ := strings.Cut(args, "-")

		fn, ok := callbacksTake(&callbacks, id)
		if !ok || fn == nil {
			return
		}

		data := make(map[string]string)

		for _, ud := range userdata {
			data[ud.Key] = ud.Value
		}

		fn(nil, args, data)
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
		actions = make([]windowsnotify.Action, 0, len(n.ButtonActions))
		inputs  = make([]windowsnotify.Input, 0, len(n.TextActions))
	)

	for _, a := range n.ButtonActions {
		actions = append(actions, windowsnotify.Action{
			Content:   a.LabelText,
			Arguments: a.AppPayload,
		})
	}

	for _, a := range n.TextActions {
		actionID := strconv.FormatInt(nextID.Add(1), 10)
		inputs = append(inputs, windowsnotify.Input{
			ID:          actionID,
			Title:       a.Title,
			Placeholder: a.PlaceholderHint,
		})
		if a.ButtonLabel != "" {
			actions = append(actions, windowsnotify.Action{
				InputID: actionID,
				Content: a.ButtonLabel,
			})
		}
	}

	tn := windowsnotify.Notification{
		Title:               n.Title,
		Body:                n.Body,
		Icon:                n.Icon,
		ActivationType:      windowsnotify.Foreground,
		ActivationArguments: strconv.FormatInt(id, 10),
		Actions:             actions,
		Inputs:              inputs,
	}

	return tn.Push()
}
