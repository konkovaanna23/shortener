package file

import "os"

func SaveToFile(fileName string, data []byte) error {
	err := os.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ReadFromFile(fileName string) ([]byte, error) {
	result, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return result, err

}
