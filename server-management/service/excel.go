package service

import (
	"fmt"
	"log"
	db "server-management/db/sqlc"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func ImportServer(path string) ([][]string, error) {

	f, err := excelize.OpenFile(path)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Errorf("%w", err)
		}
	}()

	rows, err := f.GetRows("Result 1")

	if err != nil {
		return nil, err
	}

	return rows, nil
}

func ExportServer(servers []db.Server) (string, error) {
	f := excelize.NewFile()

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	index, err := f.NewSheet("Result 1")
	if err != nil {
		return "", err
	}

	// set value name column
	f.SetCellValue("Result 1", "A1", "id")
	f.SetCellValue("Result 1", "B1", "name")
	f.SetCellValue("Result 1", "C1", "ipv4")
	f.SetCellValue("Result 1", "D1", "status")
	f.SetCellValue("Result 1", "E1", "created_at")
	f.SetCellValue("Result 1", "F1", "updated_at")

	f.SetActiveSheet(index)

	// query of filter

	if err != nil {
		return "", nil
	}
	// fetch data into new file excel

	for i, data := range servers {
		var tmp = strconv.Itoa(i + 2)
		f.SetCellValue("Result 1", "A"+tmp, data.ID)
		f.SetCellValue("Result 1", "B"+tmp, data.Name)
		f.SetCellValue("Result 1", "C"+tmp, data.Ipv4)
		f.SetCellValue("Result 1", "D"+tmp, data.Status)
		f.SetCellValue("Result 1", "E"+tmp, data.CreatedAt.Time)
		f.SetCellValue("Result 1", "F"+tmp, data.UpdateAt.Time)
	}

	f.SetActiveSheet(index)

	if err := f.SaveAs("export.xlsx"); err != nil {
		fmt.Errorf("error %w", err)
	}
	return f.Path, nil
}
