package logger

// Contains key:value representation
// of data parameters provided to standard Logger
// functions.
type Data map[string]string

func (d Data) Len() int {
	return len(d)
}

func createDataFromSlice(values []string) Data {
	if len(values)%2 != 0 {
		panic("number of items in slice provided for Data must be even")
	}

	res := make(Data)
	lastKey := ""

	for idx, val := range values {
		if idx%2 == 0 {
			lastKey = val

			continue
		}

		res[lastKey] = val
	}

	return res
}
