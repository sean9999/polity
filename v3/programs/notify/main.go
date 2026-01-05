package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/0xAX/notificator"
)

var notify *notificator.Notificator

func darwinNotify(title string, text string, icon string) error {

	notify = notificator.New(notificator.Options{
		DefaultIcon: "icon/default.png",
		AppName:     "Polity",
	})

	return notify.Push(title, text, icon, notificator.UR_CRITICAL)

}
func Notify(summary string, body string, isUrgent bool, iconPath string) error {

	switch runtime.GOOS {

	case "darwin":

		return darwinNotify(summary, body, iconPath)

	case "linux":

		var cmd *exec.Cmd
		if isUrgent {
			cmd = exec.Command("notify-send", "-i", iconPath, summary, body, "-u", "critical")
		} else {
			cmd = exec.Command("notify-send", "-i", iconPath, summary, body)
		}
		cmd.Env = os.Environ()
		return cmd.Run()

	case "android", "dragonfly", "freebsd", "nacl", "netbsd", "openbsd", "plan9", "solaris", "windows":
		fallthrough
	default:
		return errors.New("not implemented")

	}

}

func main() {
	err := Notify("summary", "body is the body of the fruit of the loom", true, "/usr/share/icons/Humanity/emblems/48/emblem-OK.svg")
	fmt.Println(err)
}
