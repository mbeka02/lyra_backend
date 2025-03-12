package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // return random int btwn min and max
}

func RandString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		character := alphabet[rand.Intn(k)]
		sb.WriteByte(character)

	}
	return sb.String()
}

func RandName() string {
	name := RandString(15)
	return name
}

func RandEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandString(10))
}

func RandRole() string {
	roles := []string{"patient", "specialist"}
	return roles[rand.Intn(len(roles))]
}

func RandPhoneNumber() string {
	countryCodes := []int{253, 254, 255, 256, 257, 258, 259}
	countryCode := countryCodes[rand.Intn(len(countryCodes))]
	// Local number digits vary by country, I'm sticking with 9
	localNumber := rand.Int63n(1000000000)
	return fmt.Sprintf("+%d%09d", countryCode, localNumber)
}

func RandDateOfBirth(minAge, maxAge int) time.Time {
	now := time.Now()
	minYear := now.Year() - maxAge
	maxYear := now.Year() - minAge

	year := minYear + rand.Intn(maxYear-minYear+1)
	month := time.Month(rand.Intn(12) + 1) // 1-12

	// Get the correct number of days for the month/year (handles leap years)
	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	day := rand.Intn(daysInMonth) + 1 // 1-[28,29,30,31]

	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func RandDateOfBirthString(minAge, maxAge int, format string) string {
	if format == "" {
		format = "2006-01-02" // Default ISO format
	}
	return RandDateOfBirth(minAge, maxAge).Format(format)
}
