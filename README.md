A Go application launched from a shell to display a menu timer for macOS.

It uses [zenith](https://github.com/ncruces/zenity) to display the 'Timer is finished' dialog and
a notification and [menuet](https://github.com/caseymrm/menuet) to display the menu timer.

---

Inspired by https://github.com/kristopherjohnson/MenubarCountdown which I previously used and liked,
but I wanted to invoke the countdown from a shell and customize it more.

Multiple timers are supported (I just launch a new tab in iTerm) and the timer keeps 
counting once it reaches 00:00:00 which is useful when one misses a countdown and wants
to see how much time passed.


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
```

I rarely use the options to set seconds or hours, so I just run it like `countdown 15` to set a timer for 15 minutes.

I also have the application renamed to `tim` because it is a lot shorter and invoke it like `tim 15`.


## Caveats

There is the issue with the [menuet](https://github.com/caseymrm/menuet) library that it does not have any 
functionality to remove the 'Start at Login' and 'Quit' menu items or register any click events on them.

This is especially problematic because clicking the 'Quit' menu item correctly quits the `countdown` application 
but still leaves the `caffeinate -i` process (started in a Goroutine) running.
`caffeinate -i` prevents system sleep and it is not good to keep it running without any reason.

Registering a signal to still catch the 'Quit' click apparently causes an internal panic with 
[menuet](https://github.com/caseymrm/menuet) (maybe because it is used there internally as well?),
so there is really no way to correctly handle the 'Quit' menu item, so better not click it!

![]()
