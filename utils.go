package main

import (
	"log"
	"strconv"
	"strings"
)

func increaseID(ID string) (string, error) {
	parts := strings.Split(ID, "-")
	num, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Printf("Error parsing int: %s", err)
		return "", err
	}

	return parts[0] + "-" + strconv.Itoa(num+1), nil
}
