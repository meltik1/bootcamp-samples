package generator

import (
	"fmt"
	"log"
	"os"
	"strconv"
	_ "time"

	"github.com/devhands-io/bootcamp-samples/golang/vanilla/models"
)

func GeneratePhonesRow(index int) models.User {
	return models.User{
		ID:      index,
		Name:    fmt.Sprintf("robert_%d", index),
		Surname: fmt.Sprintf("paulson_%d", index),
	}
}

func WritePhoneRowIntoCsv(f *os.File, user models.User) error {
	_, err := fmt.Fprintf(f, "%d,%s,%s\n",
		user.ID,
		user.Name,
		user.Surname,
	)

	if err != nil {
		return fmt.Errorf("Error in writeIntocsv", err)
	}

	return nil
}

const rows = 1_000
const fileNamePhonesV1 = "user.csv"

func DO() {
	for i := 1; i <= 1; i++ {
		err := usersV1(i, "data/"+strconv.FormatInt(int64(i), 10)+"_")

		if err != nil {
			log.Fatal("Error while generating phones")
		}
	}

}

func usersV1(step int, filename string) error {
	file, err := openPhonesFileAndWriteHeader(filename + fileNamePhonesV1)
	if err != nil {
		log.Fatal("Error occured", err)
	}
	defer file.Close()

	i := int(0)
	for ; i < rows*step; i++ {
		row := GeneratePhonesRow(i)
		err := WritePhoneRowIntoCsv(file, row)
		if err != nil {
			return err
		}
	}
	return nil
}

func openPhonesFileAndWriteHeader(fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Error while opening file", err)
	}

	_, err = fmt.Fprint(file, "userId, name, surname\n")
	if err != nil {
		return nil, fmt.Errorf("Error while writing header", err)
	}

	return file, nil
}
