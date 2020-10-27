package baomi

import "errors"

var baomidengji = map[string]int{
	"gongkai": 1,
	"mimi":    2,
	"jimi":    3,
	"juemi":   4,
}

//IsContain Determine whether users can access resources
func IsContain(userBaomi string, resourceBaomi string) (bool, error) {
	if baomidengji[userBaomi] == 0 {
		return false, errors.New("Baomidengji is not found")
	}
	if baomidengji[resourceBaomi] == 0 {
		return false, errors.New("Baomidengji is not found")
	}
	return baomidengji[userBaomi] >= baomidengji[resourceBaomi], nil
}
