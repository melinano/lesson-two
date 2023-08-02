package dev01

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
)

/*
=== Базовая задача ===

Создать программу, печатающую точное время с использованием NTP библиотеки. Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу, печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

// Variable to store the NTP server URL
var ntpServerURL = "ru.pool.ntp.org"

func PrintCurrentTime() {
	// get the exact time from the NTP-server
	time, err := ntp.Time(ntpServerURL)
	if err != nil {
		// Handle the error by printing to STDERR and returning a non-zero exit
		fmt.Fprintf(os.Stderr, "Error fetching NTP time: %v\n", err)
		os.Exit(1)
	}

	// Print the exact time obtained from NTP server
	fmt.Println(time.Format("2 Jan 2006 15:04:05.999999"))
}
