package util

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

func IsEmail(inputString string) (bool, error) {
	return regexp.MatchString(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`, inputString)
}

func IsPhone(inputString string) (bool, error) {
	return regexp.MatchString(`^[0-9]*$`, inputString)
}

func IsUsername(inputString string) (bool, error) {
	return regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]*$`, inputString)
}

func RandomFileName(length int) string {
	rand.Seed(time.Now().UnixNano())
	stringTime := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s%d", stringTime, rand.Intn(1000000))
}

func RemoveSpaces(inputString string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(inputString, "")
}
