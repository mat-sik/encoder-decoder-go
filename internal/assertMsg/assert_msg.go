package assertMsg

import "fmt"

func GetAssertMsg[T any, U any](input T, expected U) string {
	return fmt.Sprintf("Result should be %v for '%v'", expected, input)
}
