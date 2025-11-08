package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	intro()

	doneChan := make(chan bool)
	go readInput(os.Stdin, doneChan)

	<-doneChan

	close(doneChan)

	fmt.Println("Goodbye!")
}

func readInput(in io.Reader, doneChan chan bool) {
	// Creo un scanner asociado a os.Stdin, que es la entrada estándar del sistema operativo.
	// Esta entrada representa un flujo de datos (generalmente el teclado) que el SO expone como un archivo especial.
	// El scanner leerá desde ese flujo, línea por línea, cualquier dato que el usuario escriba o que se redirija al programa
	scanner := bufio.NewScanner(in)

	for {
		res, done := checkNumber(scanner)

		if done {
			doneChan <- true
			return
		}

		fmt.Println(res)
		prompt()
	}
}

func checkNumber(scanner *bufio.Scanner) (string, bool) {
	scanner.Scan() // pause for input

	if strings.EqualFold(scanner.Text(), "q") {
		return "", true
	}

	numToCheck, erro := strconv.Atoi(scanner.Text())
	if erro != nil {
		return "Please enter a valid number", false
	}

	_, msg := isPrime(numToCheck)
	return msg, false

}

func intro() {
	fmt.Println("Is it prime?")
	fmt.Println("============")
	fmt.Println("Enter a number to check if it is prime. Enter q to quit.")
	prompt()
}

func prompt() {
	fmt.Print("-> ")
}

func isPrime(n int) (bool, string) {
	if n == 0 || n == 1 {
		return false, fmt.Sprintf("%d is not prime", n)
	}

	if n < 0 {
		return false, "Negative numbers are not prime"
	}

	for i := 2; i <= n/2; i++ {
		if n%i == 0 {
			return false, fmt.Sprintf("%d is not a prime number because it is divisible by %d", n, i)
		}
	}

	return true, fmt.Sprintf("%d is a prime number", n)

}
