package main

import (
	"log"
	"os"
)

func  Generate(fileName string) (error)  {
	// Подготавливаем типы для записи

	// Парсим файл


	// На основе спарсенных данных собираем из темплейтов результат


	err := writeResult(fileName, resultData)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func writeResult(fileName string, resultData []byte) error {
	fileName = "validate_" + fileName

	resultFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer resultFile.Close()

	_, err = resultFile.Write(resultData)
	if err != nil {
		return err
	}

	return nil
}