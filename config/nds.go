package config

import (
	"fmt"
	"strconv"
)

func toDecimal(N int, num float64) string {
	newNum := (int(num * 100))
	totalZeros := N - 2
	toString := strconv.Itoa(newNum)
	zeroString := ""
	result := ""

	for i := 0; i <= totalZeros; i++ {
		zeroString = "0" + zeroString
	}
	result = zeroString + toString
	fmt.Println("hasil =", result)
	return result
}

func seriesBall(height int) []float64 {
	toFloat := float64(height)
	var result []float64
	result = append(result, toFloat)
	for toFloat > 0.5 {
		toFloat = toFloat / 2
		result = append(result, toFloat)
	}
	fmt.Println("ketinggian =", result)
	return result
}

func rentBill(checkin int, checkout int) int {
	result := 0
	totalHours := checkout - checkin

	if totalHours <= 1 {
		result = 350000
	} else if totalHours <= 2 {
		result = 500000
	} else if totalHours > 2 {
		if totalHours <= 8 {
			result = 500000 + (((totalHours - 2) * 2) * 75000)
		} else if totalHours > 8 {
			result = 500000 + ((6 * 2) * 75000) + (((totalHours - 8) * 2) * 100000)
		}
	}
	fmt.Println("total =", result)
	return result
}

func seriesAdd(num1 int, num2 int, N int) []int {
	var result []int
	result = append(result, num1)
	result = append(result, num2)
	for i := 0; i < N-2; i++ {
		result = append(result, result[i]+result[i+1])
	}
	fmt.Println("hasil =", result)
	return result
}

func loanSim(amount int, tenor int, interest float64) float64 {
	var result float64
	result = float64(amount) + (float64(amount*tenor) * (interest))
	fmt.Println("total =", result)
	return result
}
