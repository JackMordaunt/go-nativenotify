package nativenotify

import (
	"fmt"
	"strconv"
	"strings"

	windowsnotify "git.sr.ht/~jackmordaunt/go-toast/v2"
)

func setup(cfg Config) error {
	windowsnotify.SetAppData(cfg.Windows)
	windowsnotify.SetActivationCallback(func(args string, userdata []windowsnotify.UserData) {
		parts := strings.Split(args, "-")

		id, _ := take(&parts)

		actionIDEncoded, _ := take(&parts)
		actionArgsEncoded, _ := take(&parts)

		actionID := decode(actionIDEncoded)
		actionArgs := decode(actionArgsEncoded)

		data := make(map[string]string)

		for _, ud := range userdata {
			if ud.Value != "" {
				data[ud.Key] = ud.Value
			}
		}

		if actionID != "" {
			data[actionID] = actionArgs
		}

		fn, ok := callbacksTake(&callbacks, id)
		if !ok || fn == nil {
			return
		}

		fn(nil, args, data)
	})
	return nil
}

// push the notification.
//
// Actions have the notification ID embedded in them so that we can retrieve the
// correct procedure when actvited by Windows.
//
// Action IDs and payloads are hex encoded to escape them and decoded to plaintext
// when received.
func push(n Notification) (err error) {
	id := nextID.Add(1)

	var (
		actions = make([]windowsnotify.Action, 0, len(n.ButtonActions))
		inputs  = make([]windowsnotify.Input, 0, len(n.TextActions))
	)

	for _, a := range n.ButtonActions {
		actions = append(actions, windowsnotify.Action{
			Content:   a.LabelText,
			Arguments: fmt.Sprintf("%d-%s-%s", id, encode(a.ID), encode(a.AppPayload)),
		})
	}

	for _, a := range n.TextActions {
		inputs = append(inputs, windowsnotify.Input{
			ID:          a.ID,
			Title:       a.Title,
			Placeholder: a.PlaceholderHint,
		})
		if a.ButtonLabel != "" {
			actions = append(actions, windowsnotify.Action{
				InputID:   a.ID,
				Content:   a.ButtonLabel,
				Arguments: fmt.Sprintf("%d", id),
			})
		}
	}

	tn := windowsnotify.Notification{
		Title:               n.Title,
		Body:                n.Body,
		Icon:                n.Icon,
		ActivationType:      windowsnotify.Foreground,
		ActivationArguments: fmt.Sprintf("%d-%s", id, encode(n.ID)),
		Actions:             actions,
		Inputs:              inputs,
	}

	// This closure ensures that the outer notification ID and payload is always passed to the callback.
	// The action data will appear in [userData].
	callbacksPut(&callbacks, strconv.FormatInt(id, 10), func(err error, args string, userData map[string]string) {
		if n.AppPayload != "" {
			userData["default"] = n.AppPayload
		}
		n.Callback(err, n.ID, userData)
	})

	return tn.Push()
}
