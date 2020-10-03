# macos-menu-countdown

Inspired by https://github.com/kristopherjohnson/MenubarCountdown which I really liked but it missed some features
and had features I simply never used, so I wanted to create my own menu timer application.

This application allows being run multiple times, so multiple timers can be run in parallel which was very important to me.

It also keeps counting once it reaches 00:00:00 which is useful when one misses a timer countdown and wants
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

