package file

import (
	"fmt"
	"os"
)

func SaveToFile(fileName string, data []byte) error {
	err := os.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}
	fmt.Println("Команда сохранения в файл отработала ", fileName, " данные ", string(data))
	return nil
}

func ReadFromFile(fileName string) ([]byte, error) {
	result, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return result, err

}
