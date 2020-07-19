package main

import (
	"api"
)

func main() {
	var iterations = []int32{49927, 50391, 50849, 51947, 52848, 53258, 54100, 54611, 55531, 55991, 56471, 57010, 57503, 58257, 59061, 59805, 60539, 61354, 62204, 62969, 63771, 64559, 65422, 67092, 67105}
	for _, uid := range iterations {
		api.TrainFromLastIteration(uid)
	}

}
