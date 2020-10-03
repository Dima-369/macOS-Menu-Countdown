A Go application launched from a shell to display a menu timer for macOS.

It uses [zenith](https://github.com/ncruces/zenity) to display the 'Timer is finished' dialog and
a notification and [menuet](https://github.com/caseymrm/menuet) to display the menu timer.

---

Inspired by https://github.com/kristopherjohnson/MenubarCountdown which I previously liked a lot but
I wanted to invoke the countdown from a shell.

Multiple timers also supported (I just launch a new tab in iTerm) and the timer keeps 
counting once it reaches 00:00:00 which is useful when one misses a countdown and wants
to see how much time passed.


## Usage

```bash
> countdown
 ,15       (15 seconds)
  25       (25 minutes)
  25,      (25 minutes)
  25,20    (25 minutes and 20 seconds)
  1,25,120 (1 hour, 25 minutes and 120 seconds)
```

