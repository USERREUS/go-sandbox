//Напишите функцию, которая убирает дубликаты, сохраняя порядок слайса

package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	output := make([]string, 0)
	inputSet := make(map[string]int, len(input))
	for _, v := range input {
		if _, ok := inputSet[v]; !ok {
			output = append(output, v)

		}
		inputSet[v]++
	}

	return output
}

func main() {
	input := []string{
		"cat",
		"dog",
		"bird",
		"dog",
		"parrot",
		"cat",
	}
	fmt.Println(input)

	fmt.Println(RemoveDuplicates(input))
}