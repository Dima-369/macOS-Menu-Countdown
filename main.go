package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/caseymrm/menuet"
	"github.com/ncruces/zenity"
)

// timerFinishedAudioFile specifies the audio file which is played once the timer is finished.
//
// This file should be in the same directory as this executable.
//
// The file is played by invoking 'ffplay'.
const timerFinishedAudioFile = "you-can-heal.mp3"

const (
	secondsInMinute = 60
	secondsInHour   = 60 * 60
	hoursInDay      = 24
	timeStep        = time.Millisecond * 200
)

// this global variable is apparently required because one can't pass arguments to properQuitMenuItem().
var caffeinatePID = 0

type countdown struct {
	hours   int
	minutes int
	seconds int
}

func (c countdown) toString() string {
	return fmt.Sprintf("%0.2d:%0.2d:%0.2d", c.hours, c.minutes, c.seconds)
}

func (c countdown) isOverTime() bool {
	return c.hours <= 0 && c.minutes <= 0 && c.seconds <= 0
}

func (c *countdown) flipForOverTime() {
	if c.hours < 0 {
		c.hours = -c.hours
	}

	if c.minutes < 0 {
		c.minutes = -c.minutes
	}

	if c.seconds < 0 {
		c.seconds = -c.seconds
	}
}

func getRemainingTime(endTime time.Time) countdown {
	now := time.Now()
	difference := endTime.Sub(now)

	total := int64(difference.Seconds())
	hours := total / (secondsInHour) % hoursInDay
	minutes := (total / secondsInMinute) % secondsInMinute
	seconds := total % secondsInMinute

	return countdown{
		hours:   int(hours),
		minutes: int(minutes),
		seconds: int(seconds),
	}
}

func getSecondCountAsHumanString(c int) string {
	out := ""
	hours := c / (secondsInHour) % hoursInDay
	minutes := c / secondsInMinute % secondsInMinute
	seconds := c % secondsInMinute

	if hours == 1 {
		out += "1 hour"
	} else if hours > 1 {
		out += fmt.Sprintf("%d hours", hours)
	}

	if out != "" && minutes > 0 {
		out += ", "
	}

	if minutes == 1 {
		out += "1 minute"
	} else if minutes > 1 {
		out += fmt.Sprintf("%d minutes", minutes)
	}

	if out != "" && seconds > 0 {
		out += ", "
	}

	if seconds == 1 {
		out += "1 second"
	} else if seconds > 1 {
		out += fmt.Sprintf("%d seconds", seconds)
	}

	return out
}

func properQuitMenuItem() []menuet.MenuItem {
	return []menuet.MenuItem{
		{
			Text: "Proper Quit",
			Clicked: func() {
				exitAndKillCaffeinate(0)
			},
		},
	}
}

func countDown(startTime time.Time, timerName string, totalCount int) {
	menuet.App().Label = fmt.Sprintf("%d", caffeinatePID)
	menuet.App().Children = properQuitMenuItem

	countDown := time.Duration(totalCount) * time.Second
	doneOn := startTime.Add(countDown)

	isOverTime := false

	for {
		remaining := getRemainingTime(doneOn)
		menuString := ""

		if isOverTime {
			remaining.flipForOverTime()

			if remaining.seconds >= 1 {
				menuString = "-" + remaining.toString()
			} else {
				// to just display 00:00:00
				menuString = remaining.toString()
			}
		} else {
			menuString = remaining.toString()
		}

		title := menuString
		if timerName != "" {
			title = timerName + " " + title
		}

		menuet.App().SetMenuState(&menuet.MenuState{
			Title: title,
		})

		if remaining.isOverTime() && !isOverTime {
			isOverTime = true

			go timerIsUp(totalCount)
		}

		time.Sleep(timeStep)
	}
}

func playFinishedSound() {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}

	path := filepath.Dir(exe)

	// #nosec
	err = exec.Command("afplay", path+"/"+timerFinishedAudioFile).Run()
	if err != nil {
		panic(err)
	}
}

func timerIsUp(totalCount int) {
	killCaffeinate()

	forHuman := getSecondCountAsHumanString(totalCount)

	text := ""
	if strings.HasSuffix(forHuman, "s") {
		text = fmt.Sprintf("%s have passed.", forHuman)
	} else {
		text = fmt.Sprintf("%s has passed.", forHuman)
	}

	err := zenity.Notify(text,
		zenity.Title("Timer is finished"),
		zenity.Icon(zenity.InfoIcon))
	if err != nil {
		panic(err)
	}

	go playFinishedSound()

	_, err = zenity.Info(text,
		zenity.Title("Timer is finished"),
		zenity.Icon(zenity.InfoIcon))
	if err != nil {
		panic(err)
	}

	os.Exit(0)
}

func safeAtoi(s string) int {
	if s == "" {
		return 0
	}

	parsed, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return parsed
}

func parseStringCountToSeconds(s string) int {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "m", "")
	s = strings.ReplaceAll(s, "h", "")
	s = strings.ReplaceAll(s, "s", "")

	if strings.Contains(s, ",") {
		const (
			inMinutesSecondsFormat      = 2
			inHoursMinutesSecondsFormat = 3
		)

		parts := strings.Split(s, ",")
		switch len(parts) {
		case inMinutesSecondsFormat:
			return safeAtoi(parts[0])*secondsInMinute + safeAtoi(parts[1])
		case inHoursMinutesSecondsFormat:
			return safeAtoi(parts[0])*secondsInHour + safeAtoi(parts[1])*secondsInMinute + safeAtoi(parts[2])
		}
	} else {
		// just minutes
		return safeAtoi(s) * secondsInMinute
	}

	println(fmt.Sprintf("Problematic time format: %s\n", s))
	printUsage()
	os.Exit(1)

	// the return value here is really not important
	return -1
}

func printUsage() {
	println("Usage:\n" +
		"  countdown {time option} {optional timer name}\n\n" +
		"Valid time options are:\n" +
		" ,15       (15 seconds)\n" +
		"  25       (25 minutes)\n" +
		"  25,      (25 minutes)\n" +
		"  25,20    (25 minutes and 20 seconds)\n" +
		"  1,25,120 (1 hour, 25 minutes and 120 seconds)")
}

// waitForStdinToQuit queries stdin for an Enter to abort the program.
//
// Using a signal notifier like: signal.Notify(c, os.Interrupt, syscall.SIGTERM)
// causes an internal crash with the menu bar C code it seems and is not fixable.
func waitForStdinToQuit(startTime time.Time, totalSeconds int) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Hit Enter to cancel > ")

	_, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	doneOn := startTime.Add(time.Second * time.Duration(totalSeconds))
	remaining := getRemainingTime(doneOn)

	if remaining.isOverTime() {
		remaining.flipForOverTime()
		fmt.Printf("\n%s over time...\n", remaining.toString())
	} else {
		fmt.Printf("\n%s left...\n", remaining.toString())
	}

	exitAndKillCaffeinate(0)
}

func killCaffeinate() {
	// #nosec
	cmd := exec.Command("kill", strconv.Itoa(caffeinatePID))

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	// we do not check for errors here because the timer might have already been killed
	_ = cmd.Wait()
}

func exitAndKillCaffeinate(exitCode int) {
	killCaffeinate()
	os.Exit(exitCode)
}

// preventSystemSleepViaCaffeinate makes sure that the system does not
// idle sleep to keep the timer running correctly.
//
// This still allows the display to turn off.
func preventSystemSleepViaCaffeinate() {
	cmd := exec.Command("caffeinate", "-i")

	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	pid := cmd.Process.Pid

	go func() {
		// when the timer is up, the caffeinate process is killed, so we do not check for errors here
		_ = cmd.Wait()
	}()

	caffeinatePID = pid
}

func main() {
	count := ""

	const hasArg = 2
	if len(os.Args) >= hasArg {
		count = os.Args[1]
	} else {
		printUsage()
		os.Exit(1)
	}

	startTime := time.Now()
	countInSeconds := parseStringCountToSeconds(count)

	preventSystemSleepViaCaffeinate()

	go waitForStdinToQuit(startTime, countInSeconds)

	timerName := ""

	const hasTimerName = 3
	if len(os.Args) >= hasTimerName {
		timerName = os.Args[2]
	}

	go countDown(startTime, timerName, countInSeconds)

	menuet.App().RunApplication()
}
