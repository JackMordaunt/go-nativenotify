package nativenotify

import (
	"fmt"
	"strconv"
	"sync/atomic"

	"git.sr.ht/~whereswaldon/shout"
	"github.com/godbus/dbus/v5"
)

var notifier atomic.Pointer[shout.Notifier]

func setup(cfg Config) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return fmt.Errorf("getting dbus session: %w", err)
	}

	n, err := shout.NewNotifier(
		conn,
		cfg.Linux.AppName,
		cfg.Linux.AppIcon,
		func(id, action string, platformData map[string]dbus.Variant, target, response dbus.Variant, err error) {
			fn, ok := callbacksTake(&callbacks, id)
			if !ok || fn == nil {
				return
			}
			fn(action, target.String())
		},
	)
	if err != nil {
		return fmt.Errorf("building notifier: %w", err)
	}

	notifier.Store(&n)

	return nil
}

func push(n Notification) (err error) {
	notifier := notifier.Load()

	if notifier == nil {
		return fmt.Errorf("notifier is nil, call setup to initialize")
	}

	id := nextID.Add(1)

	buttons := []shout.Button{}

	for _, a := range n.ButtonActions {
		buttons = append(buttons, shout.Button{
			Action: a.ID,
			Label:  a.LabelText,
			Target: a.Value,
		})
	}

	if err := (*notifier).Send(fmt.Sprintf("%d", id), shout.Notification{
		Title:               n.Title,
		Body:                n.Body,
		ReplaceID:           "",
		Markup:              false,
		IconPath:            n.Icon,
		Priority:            shout.Normal,
		DefaultAction:       "default",
		DefaultActionLabel:  "",
		DefaultActionTarget: dbus.Variant{},
		Buttons:             buttons,
		ExpirationTimeout:   0,
	}); err != nil {
		return fmt.Errorf("sending notification: %w", err)
	}

	callbacksPut(&callbacks, strconv.FormatInt(id, 10), n.Callback)

	return nil
}
