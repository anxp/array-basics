package array_basics

// TypeScalar more details about generic types: https://go.dev/blog/intro-generics
type TypeScalar interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float64 | ~float32 | ~string
}

type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float64 | ~float32
}

// ElementIndex returns the index of the first occurrence of v in s, or -1 if not present.
func ElementIndex[E comparable](s []E, v E) int {
	for i, vs := range s {
		if v == vs {
			return i
		}
	}
	return -1
}

// InArray checks if element v is in slice s. Function acts as PHP in_array() function
func InArray[E comparable](s []E, v E) bool {
	return ElementIndex(s, v) >= 0
}

// ArrayMap returns a modified copy of the slice passed as an argument
// https://golangprojectstructure.com/functional-programming-with-slices/
func ArrayMap[Tin, Tout any](slice []Tin, callback func(value Tin, index int) Tout) []Tout {
	mappedSlice := make([]Tout, len(slice))

	for i, v := range slice {
		mappedSlice[i] = callback(v, i)
	}

	return mappedSlice
}

func ArrayUnique[T TypeScalar](slice []T) []T {
	size := len(slice)
	filteredResult := make([]T, 0, size)
	temp := map[T]struct{}{}

	for i := 0; i < size; i++ {
		if _, ok := temp[slice[i]]; ok == false {
			temp[slice[i]] = struct{}{}
			filteredResult = append(filteredResult, slice[i])
		}
	}

	return filteredResult
}

// ArrayIntersect returns a slice containing only COMMON values for all specified slices
func ArrayIntersect[T TypeScalar](slices ...[]T) []T {
	inSlicesCount := len(slices)
	temporaryFiltrationStorage := make([]map[T]int, inSlicesCount)

	for i := 0; i < inSlicesCount; i++ {
		tmpStorage := make(map[T]int, len(slices[i]))

		for j := 0; j < len(slices[i]); j++ {
			tmpStorage[slices[i][j]] = 1
		}

		temporaryFiltrationStorage[i] = tmpStorage
	}

	refMap := temporaryFiltrationStorage[0]

	for i := 1; i < inSlicesCount; i++ {
		currentMap := temporaryFiltrationStorage[i]

		for keyValue, _ := range currentMap {
			if _, present := refMap[keyValue]; present == true {
				refMap[keyValue]++
			}
		}
	}

	commonElements := make([]T, 0, len(refMap))

	for keyValue, counter := range refMap {
		// We collect only values which are present at least 1 time (*) in each input slice (COMMON VALUES),
		// so total number of such values is ALWAYS == number of input slices.
		// * If specific value presented twice in given slice, we count such entry as 1 due to unification by converting values to map-keys.
		if counter == inSlicesCount {
			commonElements = append(commonElements, keyValue)
		}
	}

	return commonElements
}

// ArraySubtract subtracts "small" set from "big" set,
// returns a slice that contains only those elements of the "big" slice that are NOT in the "small" slice.
// "big" and "small" are just names (for intuitiveness), actually, "small" slice can be bigger.
// Side effect #1: Order of elements in resulting array is NOT preserved;
// Side effect #2: "big" set will be UNIFIED.
func ArraySubtract[T TypeScalar](small, big []T) []T {
	invBig := make(map[T]struct{}, len(big))

	for _, bValue := range big {
		invBig[bValue] = struct{}{}
	}

	for _, sValue := range small {
		if _, present := invBig[sValue]; present == true {
			delete(invBig, sValue)
		}
	}

	result := make([]T, 0, len(invBig))

	for keyValue, _ := range invBig {
		result = append(result, keyValue)
	}

	return result
}

// IsArrayInArray checks whether an array is a subset of another, larger array.
// If ALL elements from first array ARE in second array, subset is complete, and function returns TRUE and EMPTY missed array;
// If there are elements from first array missed in second, function returns FALSE and array with missed elements.
func IsArrayInArray[T TypeScalar](small, big []T) (bool, []T) {
	bigArrayHashTable := make(map[T]struct{}, len(big))
	missedElements := make([]T, 0, len(small))

	for _, bV := range big {
		bigArrayHashTable[bV] = struct{}{}
	}

	for _, sV := range small {
		if _, present := bigArrayHashTable[sV]; present == false {
			missedElements = append(missedElements, sV)
		}
	}

	return len(missedElements) == 0, missedElements
}

// FindMedian by Torben Mogensen.
// It's not the fastest way to find the median, but it has a very interesting property:
// it doesn't change the input array while searching for the median, and it doesn't create a temporary copy of the input array.
// It becomes extremely powerful when the number of elements to consider starts to be large, and copying the input array may cause enormous overheads.
// http://ndevilla.free.fr/median/median/index.html
func FindMedian[N Numeric](data []N) N {
	elementsNum := len(data)

	// 1. First pass - find MIN and MAX values
	minV := data[0]
	maxV := data[0]

	for i := 1; i < elementsNum; i++ {
		if data[i] < minV {
			minV = data[i]
		}

		if data[i] > maxV {
			maxV = data[i]
		}
	}

	var less, greater, equal int
	var guess, maxUnderGuess, minAboveGuess N

	for {
		guess = (minV + maxV) / 2
		less, greater, equal = 0, 0, 0
		maxUnderGuess = minV
		minAboveGuess = maxV

		for i := 0; i < elementsNum; i++ {
			if data[i] < guess {
				less++
				if data[i] > maxUnderGuess {
					maxUnderGuess = data[i]
				}
			} else if data[i] > guess {
				greater++
				if data[i] < minAboveGuess {
					minAboveGuess = data[i]
				}
			} else {
				equal++
			}
		}

		// We are done if values distributed equally below and above median.
		// We also allow that the number of values can be less than half,
		// this means that there are values that belong neither to the "left" set nor to the "right" one,
		// they are exactly equal to the median.
		if less <= (elementsNum+1)/2 && greater <= (elementsNum+1)/2 {
			break
		} else if less > greater {
			maxV = maxUnderGuess
		} else {
			minV = minAboveGuess
		}
	}

	if less >= (elementsNum+1)/2 {
		return maxUnderGuess
	} else if less+equal >= (elementsNum+1)/2 {
		return guess
	} else {
		return minAboveGuess
	}
}
