package flatten

func Flatten(nested interface{}) []interface{} {
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
