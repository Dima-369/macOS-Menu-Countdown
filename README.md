A Go application launched from shell to display a menu timer for macOS.

It uses [zenith](https://github.com/ncruces/zenity) to display the 'Timer is finished' dialog and
notification plus [menuet](https://github.com/caseymrm/menuet) to display the menu timer.

*Features*

* Multiple timers are supported (just launch a new tab in iTerm because this is blocking)
* Menu timer keeps counting once it reaches 00:00:00 (useful when the initial countdown was missed)
* When quit by Enter from the shell it displays the remaining time
  * Registering any signals does not work with [menuet](https://github.com/caseymrm/menuet) apparently

![]()

---

Inspired by https://github.com/kristopherjohnson/MenubarCountdown which I previously used and liked,
but I wanted to invoke the countdown from shell and customize it more.

Digging into the 


## Installation

`go get github.com/Gira-X/macos-menu-countdown`


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

I rarely use the options to set seconds or hours, so I just run it like `countdown 15` to set a timer for 15 minutes.

I also have the application renamed to `tim` because it is nicer to invoke it like `tim 15`.


## Caveats

There is the issue with the [menuet](https://github.com/caseymrm/menuet) library that it does not have any 
functionality to remove the 'Start at Login' and 'Quit' menu items or register any click events on them.

This is especially problematic because clicking the 'Quit' menu item correctly quits the `countdown` application 
but still leaves the `caffeinate -i` process (started in a Goroutine) running.
`caffeinate -i` prevents system sleep and it is not good to keep it running without any reason.

Registering a signal to still catch the 'Quit' click apparently causes an internal panic with 
[menuet](https://github.com/caseymrm/menuet) (maybe because it is used there internally as well?),
so there is really no way to correctly handle the 'Quit' menu item, so better not click it!

![Problematic Quti Menu Item](https://raw.githubusercontent.com/Gira-X/macos-menu-countdown/master/readme-images/menu.png)
