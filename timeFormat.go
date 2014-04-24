package main

import (
    "fmt"
    "time"
)

func formatDuration(secondsElapsed int) string {
    duration := time.Duration(secondsElapsed) * time.Second
    days := int(duration.Hours() / 24)
    hours := int(duration.Hours()) - (days * 24)
    minutes := int(duration.Minutes()) - (days * 24 * 60) - (hours * 60)

    return fmt.Sprintf("%d days %d hours %d minutes", days, hours, minutes)
}
