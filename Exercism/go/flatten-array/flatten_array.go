package flatten

// switch case on type
func Flatten(nested interface{}) []interface{} {
	output := []interface{}{}

	for _, el := range nested.([]interface{}) {
		switch el.(type) {
		case int:
			output = append(output, el)
		case []interface{}:
			output = append(output, Flatten(el)...)
		}
	}

	return output
}

// checking if type casting is successful
func Flatten_v1(nested interface{}) []interface{} {
	output := []interface{}{}
	input, ok := nested.([]interface{})
	if !ok {
		return output
	}
	for _, el := range input {
		val, ok := el.(int)
		if ok {
			output = append(output, val)
		} else {
			output = append(output, Flatten(el)...)
		}
	}

	return output
}
