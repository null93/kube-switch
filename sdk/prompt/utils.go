package prompt

import (
	"strings"
)

func Min ( a, b int ) int {
	if a > b {
		return b
	}
	return a
}

func Max ( a, b int ) int {
	if a > b {
		return a
	}
	return b
}

func FindIndex ( items [] string, search string ) int {
	for i, item := range items {
		if item == search {
			return i
		}
	}
	return 0
}

func Filter ( items [] string, search string ) [] string {
	filtered := [] string {}
	for _, item := range items {
		if strings.Contains ( strings.ToLower ( item ), strings.ToLower ( search ) ) {
			filtered = append ( filtered, item )
		}
	}
	return filtered
}

func GetIndex ( items [] string, index int, defaultValue string ) string {
	if index <= len ( items ) - 1 {
		return items [ index ]
	}
	return defaultValue
}
