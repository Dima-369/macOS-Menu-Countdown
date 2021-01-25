A Go application launched from shell to display a menu timer for macOS.

It uses [zenith](https://github.com/ncruces/zenity) to display the 'Timer is finished' dialog and
notification plus [menuet](https://github.com/caseymrm/menuet) to display the menu timer.

---

**It's probably better to implement this directly in Swift to hide the 'Start at Login' and 'Quit'
menu items.**

## Features

* Multiple timers are supported (just launch a new tab in iTerm because launching this from shell is blocking)
* Menu timer keeps counting once it reaches 00:00:00 (useful when the initial countdown was missed)
  * This is displayed with a minus like -3,33 in the menu
* When quit by Enter from the shell, it displays the remaining time (don't use Ctrl+C)
  * Registering any signals does not work with [menuet](https://github.com/caseymrm/menuet) apparently

---

Inspired by https://github.com/kristopherjohnson/MenubarCountdown which I previously used and liked,
but I wanted to invoke the countdown from shell and customize it more.

Digging into the Objective C was not worth the effort for me, so I just redid this in Go.


## Installation

Once the timer is finished, an audo file named `you-can-heal.mp3` in the same directory as the application is played via `afplay`.

The `you-can-heal.mp3` audio is cut from the end of the [Medical Medium's Brain Meditation](https://podcasts.apple.com/us/podcast/003-brain-meditation-letting-go-of-fear/id1133835109?i=1000492358953).

Install by `go get github.com/Gira-X/macos-menu-countdown`

This builds the application in `$GOBIN` named `macos-menu-countdown` which is really not a nice name,
so feel free to rename it to `countdown` or `tim` (which I use personally).

Then move the `you-can-heal.mp3` audio file from this repository into `$GOBIN`.


## Usage

```bash
> countdown
 ,15       (15 seconds)
  25       (25 minutes)
  25,      (25 minutes)
  25,20    (25 minutes and 20 seconds)
  1,25,120 (1 hour, 25 minutes and 120 seconds)

> countdown 15
Hit Enter to cancel >

00:14:59 left...
```

I rarely use the options to set seconds or hours, so I usually run it like `countdown 15` to set a timer for 15 minutes.

## Caveats

There is the issue with the [menuet](https://github.com/caseymrm/menuet) library that it does not have any 
functionality to remove the 'Start at Login' and 'Quit' menu items or register any click events on them.

This is especially problematic because clicking the 'Quit' menu item correctly quits the `countdown` application 
but still leaves the `caffeinate -i` process (started in a Goroutine) running.
`caffeinate -i` prevents system sleep and it is not good to keep it running without any reason.

Registering a signal to catch the 'Quit' click causes an internal panic with 
[menuet](https://github.com/caseymrm/menuet).

**So do not click the 'Quit' menu item!**

If you still did, call `pkill caffeinate` to kill the process correctly.

