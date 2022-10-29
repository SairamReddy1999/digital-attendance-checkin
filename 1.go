package main

import "fmt"

func main() {

	var x int
	x = 3
	y := 4
	sum, prod := learnMultiple(x, y)        // Function returns two values.
	fmt.Println("sum:", sum, "prod:", prod) // Simple output.

	xf := 10.5

	yf := 11.2
	sumf, prodf := floatMultiple(xf, yf)      // Function returns two values.
	fmt.Println("sum:", sumf, "prod:", prodf) // Simple output.

	str := "GeeksforGeeks"
	str2 := "Geeks"
	length(str)
	res := check(str, str2)
	fmt.Println("check:", res)

}

func learnMultiple(x, y int) (sum, prod int) {
	return x + y, x * y // Return two values.
}

func floatMultiple(x, y float64) (sum, prod float64) {
	return x + y, x * y // Return two values.
}

func length(str string) {

	fmt.Printf("Length of the string is:%d",
		len(str))
}
func check(x string, y string) (res bool) {
	return x == y
}
