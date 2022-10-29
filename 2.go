package main

import "fmt"

func main() {

	// Shorthand declaration of array
	arr := [4]string{"chack", "gfg", "323321", "testforloop"}

	// Accessing the elements of
	// the array Using for loop
	fmt.Println("Elements of the array:")

	for i := 0; i < 3; i++ {
		fmt.Println(arr[i])
	}

	// passing the address of struct variable
	// emp8 is a pointer to the Employee struct
	emp8 := &Employee{"Sam", "Anderson", 55, 6000}

	// (*emp8).firstName is the syntax to access
	// the firstName field of the emp8 struct
	if (*emp8).age > 20 {

		fmt.Println("First Name:", (*emp8).firstName)
		fmt.Println("Age:", (*emp8).age)
	}

}

type Employee struct {
	firstName, lastName string
	age, salary         int
}
