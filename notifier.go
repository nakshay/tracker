package main

import (
	"github.com/0xAX/notificator"
)

func notify(title, body, iconPath string, severity int) {

	notifier := notificator.New(notificator.Options{
		DefaultIcon: "icon/default.png",
		AppName:     "Tracker",
	})

	switch severity {

	case 0:
		notifier.Push(title, body, iconPath, notificator.UR_CRITICAL)

	case 1:
		notifier.Push(title, body, iconPath, notificator.UR_NORMAL)
	}
}
