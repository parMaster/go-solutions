package main

import (
	"fmt"
	"reflect"
)

type greeter interface {
	greet(string) string
}

type russian struct{}
type american struct{}

func (r *russian) greet(name string) string {
	return fmt.Sprintf("Привет, %s", name)
}

func (a *american) greet(name string) string {
	return fmt.Sprintf("Hello, %s", name)
}

func sayHello(g greeter, name string) {
	fmt.Println(g.greet(name))
}

func main2() {
	sayHello(&american{}, "Gusto")
	sayHello(&russian{}, "Gusto")
}

// ISP example

type animal interface {
	walker
	runner
}

type bird interface {
	walker
	flier
}

type walker interface {
	walk()
}

type runner interface {
	run()
}

type flier interface {
	fly()
}

type cat struct{}

type eagle struct{}

func (e *eagle) walk() {
	fmt.Println("eagle is walking")
}

func (e *eagle) run() {
	fmt.Println("eagle is running")
}

func (c *cat) walk() {
	fmt.Println("cat is walking")
}

func (e *eagle) fly() {
	fmt.Println("eagle is running")
}

func (c *cat) run() {
	fmt.Println("cat is running")
}

func walk(w walker) {
	w.walk()
}

func main() {
	var c animal = &cat{}
	var e bird = &eagle{}
	walk(c)
	walk(e)

	main2()
	main3()
}

// Empty Interfaces

func main3() {
	m := map[string]interface{}{}
	m["one"] = 1
	m["two"] = 2.0
	m["three"] = true

	for k, v := range m {
		switch v.(type) {
		case int:
			fmt.Printf("%s is an integer\n", k)
		case float64:
			fmt.Printf("%s is a float64\n", k)
		default:
			fmt.Printf("%s is %v\n", k, reflect.TypeOf(v))
		}
	}
}
