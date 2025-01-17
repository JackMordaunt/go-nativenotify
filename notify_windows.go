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

		// If the action was a text input, or a selection, grab the value for it.
		// If the action was a button, the actionArgs will contain the value.
		for _, ud := range userdata {
			if ud.Key == actionID {
				actionArgs = ud.Value
			}
		}

		fn, ok := callbacksTake(&callbacks, id)
		if !ok || fn == nil {
			return
		}

		fn(actionID, actionArgs)
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
			Arguments: fmt.Sprintf("%d-%s-%s", id, encode(a.ID), encode(a.Value)),
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
				Arguments: fmt.Sprintf("%d-%s", id, encode(a.ID)),
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

	switch n.Windows {
	case WindowsOptionIconCircleCrop:
		tn.IconCrop = windowsnotify.CropStyleCircle
	case WindowsOptionIconSquareCrop:
		tn.IconCrop = windowsnotify.CropStyleSquare
	case WindowsOptionIconHero:
		tn.HeroIcon, tn.Icon = tn.Icon, ""
	}

	callbacksPut(&callbacks, strconv.FormatInt(id, 10), n.Callback)

	return tn.Push()
}
