package mscore

import "fmt"

func LogInfo(text string) {
	fmt.Println("Info=>", text)
}

func LogError(text string) {
	fmt.Println("Error=>", text)
}

func LogWarning(text string) {
	fmt.Println("Warning=>", text)
}
