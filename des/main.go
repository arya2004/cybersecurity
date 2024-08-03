package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)


func main() {
	var input, key []uint8
	var decrypt bool

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nMenu:")
		fmt.Println("1. Enter input array")
		fmt.Println("2. Enter key array")
		fmt.Println("3. Set decrypt mode")
		fmt.Println("4. Call DES function")
		fmt.Println("5. Exit")
		fmt.Print("Enter your choice: ")

		scanner.Scan()
		choice, _ := strconv.Atoi(scanner.Text())

		switch choice {
		case 1:
			fmt.Print("Enter input array (e.g., 10001011): ")
			scanner.Scan()
			inputString := scanner.Text()
			input = parseArray(inputString)
		case 2:
			fmt.Print("Enter key array (e.g., 1001100111): ")
			scanner.Scan()
			keyString := scanner.Text()
			key = parseArray(keyString)
		case 3:
			fmt.Print("Enter decrypt mode (true/false): ")
			scanner.Scan()
			decryptString := scanner.Text()
			decrypt = parseBool(decryptString)
		case 4:
			if len(input) == 0 || len(key) == 0 {
				fmt.Println("Input or key array is empty. Please enter both before calling DES function.")
			} else {
				fmt.Println("Result from DES function:", DES(input, key, decrypt))
			}
		case 5:
			fmt.Println("Exiting program.")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func parseArray(input string) []uint8 {
	var result []uint8
	for _, char := range input {
		num, err := strconv.Atoi(string(char))
		if err == nil {
			result = append(result, uint8(num))
		} else {
			fmt.Println("Invalid input. Please enter a valid binary string.")
			return []uint8{}
		}
	}
	return result
}

func parseBool(input string) bool {
	value, err := strconv.ParseBool(input)
	if err != nil {
		fmt.Println("Invalid input. Please enter true or false.")
		return false
	}
	return value
}



func formatArguments(input, key string) (inputBytes, keyBytes []uint8) {
	var aux int
	for i := 1; i <= len(input); i++ {
		aux, _ = strconv.Atoi(input[i-1 : i])
		inputBytes = append(inputBytes, uint8(aux))
	}

	for i := 1; i <= len(key); i++ {
		aux, _ = strconv.Atoi(key[i-1 : i])
		keyBytes = append(keyBytes, uint8(aux))
	}

	return inputBytes, keyBytes
}

func permutation(list []uint8, positions []uint8) (permutedList []uint8) {
	// The for loop is done with the length of positions, not the list, because of P8. In the case of P8, a list of 10 positions enters and exits with 8, if it was done with len(list), it would give index out of range.
	for i := 0; i < len(positions); i++ {
		permutedList = append(permutedList, list[positions[i]-1])
	}
	return permutedList
}

func generateKeys(key []uint8) (k1, k2 []uint8) {
	P10 := []uint8{3, 5, 2, 7, 4, 10, 1, 9, 8, 6}
	LS1 := []uint8{2, 3, 4, 5, 1}
	LS2 := []uint8{3, 4, 5, 1, 2}
	P8 := []uint8{6, 3, 7, 4, 8, 5, 10, 9}

	key = permutation(key, P10)
	halfKey1, halfKey2 := key[:5], key[5:10]
	halfKey1 = permutation(halfKey1, LS1)
	halfKey2 = permutation(halfKey2, LS1)
	key = append(halfKey1, halfKey2...)
	k1 = permutation(key, P8)
	halfKey1 = permutation(halfKey1, LS2)
	halfKey2 = permutation(halfKey2, LS2)
	key = append(halfKey1, halfKey2...)
	k2 = permutation(key, P8)

	return k1, k2
}

func sw(leftSide, rightSide []uint8) ([]uint8, []uint8) {
	return rightSide, leftSide
}

func xor(input, key []uint8) (xorResult []uint8) {
	for i := 0; i < len(input); i++ {
		xorResult = append(xorResult, input[i]^key[i])
	}
	return xorResult
}

func binToInt(binValue []uint8) (intValue uint8) {
	for i := 0; i < len(binValue); i++ {
		intValue += binValue[len(binValue)-i-1] * uint8(math.Pow(float64(2), float64(i)))
	}
	return intValue
}

func sBox(bitList []uint8, sMatrix [][]uint8) (outputValue []uint8) {
	row := binToInt([]uint8{bitList[0], bitList[3]})
	column := binToInt([]uint8{bitList[1], bitList[2]})

	var tempValue uint8
	tempValue = sMatrix[row][column]
	auxString := fmt.Sprintf("%b", tempValue)
	p, _ := strconv.ParseInt(auxString, 10, 8)
	outputValue = append(outputValue, uint8(p)/10)
	outputValue = append(outputValue, uint8(p)%10)

	return outputValue
}

func functionF(leftInput, rightInput, key []uint8) (fOutput []uint8) {
	E_P := []uint8{4, 1, 2, 3, 2, 3, 4, 1}
	P4 := []uint8{2, 4, 3, 1}
	S0 := [][]uint8{
		{1, 0, 3, 2},
		{3, 2, 1, 0},
		{0, 2, 1, 3},
		{3, 1, 3, 2},
	}
	S1 := [][]uint8{
		{0, 1, 2, 3},
		{2, 0, 1, 3},
		{3, 0, 1, 0},
		{2, 1, 0, 3},
	}

	tempValue := permutation(rightInput, E_P)
	tempValue = xor(tempValue, key)
	sBoxSide0, sBoxSide1 := tempValue[0:4], tempValue[4:8]
	sBoxSide0 = sBox(sBoxSide0, S0)
	sBoxSide1 = sBox(sBoxSide1, S1)
	partialOutput := append(sBoxSide0, sBoxSide1...)
	partialOutput = permutation(partialOutput, P4)
	partialOutput = xor(leftInput, partialOutput)

	fOutput = append(partialOutput, rightInput...)

	return fOutput
}

func DES(inputText, key []uint8, decrypt bool) (outputText []uint8) {
	var IP = []uint8{2, 6, 3, 1, 4, 8, 5, 7}
	var IP_1 = []uint8{4, 1, 3, 5, 7, 2, 8, 6}
	var key_k1, key_k2 []uint8

	if decrypt {
		key_k2, key_k1 = generateKeys(key)
	} else {
		key_k1, key_k2 = generateKeys(key)
	}
	outputText = permutation(inputText, IP)
	leftSide, rightSide := outputText[:4], outputText[4:8]
	fOutput := functionF(leftSide, rightSide, key_k1)
	leftSide, rightSide = fOutput[:4], fOutput[4:8]
	leftSide, rightSide = sw(leftSide, rightSide)
	outputText = functionF(leftSide, rightSide, key_k2)
	outputText = permutation(outputText, IP_1)

	return outputText
}
