package main

// Test by example, test case #1
func ExampleCheckMagazine() {
	checkMagazine([]string{"give", "me", "one", "grand", "today", "night"}, []string{"give", "one", "grand", "today"})
	// Output:
	// Yes
}

// Test by example, test case #2
func ExampleCheckMagazine_second() {
	checkMagazine([]string{"two", "times", "three", "is", "not", "four"}, []string{"two", "times", "two", "is", "four"})
	// Output:
	// No
}

// Test by example, test case #3
func ExampleCheckMagazine_third() {
	checkMagazine([]string{"ive", "got", "a", "lovely", "bunch", "of", "coconuts"}, []string{"ive", "got", "some", "coconuts"})
	// Output:
	// No
}
