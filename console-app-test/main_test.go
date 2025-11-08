package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func Test_isPrime(t *testing.T) {
	result, msg := isPrime(0)
	if result {
		t.Errorf("with %d, as test parameter, got true, but expected false", 0)
	}

	if msg != "0 is not prime" {
		t.Error("wring message returned:", msg)
	}

	result, msg = isPrime(7)
	if !result {
		t.Errorf("with %d, as test parameter, got false, but expected false", 7)
	}

	if msg != "7 is a prime number" {
		t.Error("wring message returned:", msg)
	}
}

func Test_isPrimeV2(t *testing.T) {

	type PrimeTestV2 struct {
		name     string
		testNum  int
		expected bool
		msg      string
	}

	// slice of structs
	primeTestV2 := []PrimeTestV2{
		{"prime", 7, true, "7 is a prime number"},
		{"not prime", 10, false, "10 is not a prime number because it is divisible by 2"},
		{"zero", 0, false, "0 is not prime"},
		{"negative", -5, false, "Negative numbers are not prime"},
	}

	for _, e := range primeTestV2 {
		result, msg := isPrime(e.testNum)
		if e.expected && !result {
			t.Errorf("%s: expected true but go false", e.name)
		}

		if !e.expected && result {
			t.Errorf("%s: expected false but go true", e.name)
		}

		if e.msg != msg {
			t.Error("wring message returned:", msg)
		}
	}
}

func Test_prompt(t *testing.T) {
	oldOutToRestaured := os.Stdout // esto es para guardar el valor original de os.Stdout
	r, w, _ := os.Pipe()           // creamos una tubería
	os.Stdout = w                  // redirigimos os.Stdout a la tubería

	prompt()

	_ = w.Close()                 // cerramos el escritor de la tubería
	os.Stdout = oldOutToRestaured // restauramos el valor original de os.Stdout

	out, _ := io.ReadAll(r) // leemos lo que se escribió en la tubería
	if string(out) != "-> " {
		t.Errorf("prompt() = %q, want %q", string(out), "-> ")
	}
}

func Test_intro(t *testing.T) {
	// Guardo el valor original de os.Stdout.
	// os.Stdout es una variable global de tipo *os.File (no una función).
	// Este puntero apunta a una estructura en memoria que representa la salida estándar del sistema operativo,
	// la cual está asociada físicamente al dispositivo o archivo de salida (por ejemplo, la consola /dev/stdout).
	oldOut := os.Stdout
	read, write, _ := os.Pipe() // creamos una tubería estandard input(os.Stdin) y output(os.Stdout)
	os.Stdout = write           // redirigimos os.Stdout a la tubería (que lo que escribamos en os.Stdout irá a la tubería)

	intro()

	_ = write.Close()  // cerramos el escritor de la tubería
	os.Stdout = oldOut // restauramos el valor original de os.Stdout

	out, _ := io.ReadAll(read) // leemos lo que se escribió en la tubería

	if !strings.Contains((string(out)), "Is it prime?") {
		t.Errorf("intro text not correct; got %q", string(out))
	}
}

func Test_checkNumber(t *testing.T) {
	test := []struct {
		name     string
		input    string
		expected string
	}{
		{"valid prime", "7", "7 is a prime number"},
		{"valid not prime", "10", "10 is not a prime number because it is divisible by 2"},
		{"invalid input", "hello", "Please enter a valid number"},
		{"quit input", "q", ""},
	}

	for _, e := range test {
		scanner := bufio.NewScanner(strings.NewReader(e.input))
		result, _ := checkNumber(scanner)
		if result != e.expected {
			t.Errorf("%s: expected %q but got %q", e.name, e.expected, result)
		}
	}
}

func TestReadInput(t *testing.T) {
	doneChan := make(chan bool)

	// creo referencia a un bytes.Buffer para simular la entrada estándar
	var stdin bytes.Buffer
	stdin.Write([]byte("1\nq\n"))

	go readInput(&stdin, doneChan)

	<-doneChan

	close(doneChan)
}
