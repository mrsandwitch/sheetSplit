package main

import (
	"bufio"
	//"encoding/csv"
	"fmt"

	// "io"
	"strings"

	"github.com/tealeg/xlsx"

	//"log"
	//"net/http"
	"os"
	//"strings"
)

// func downloadFile(filepath string, url string) error {
// 	// Create the file
// 	out, err := os.Create(filepath)
// 	if err != nil {
// 		return err
// 	}
// 	defer out.Close()

// 	// Get the data
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	// Check server response
// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("bad status: %s", resp.Status)
// 	}

// 	// Writer the body to file
// 	_, err = io.Copy(out, resp.Body)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func test() {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "I am a cell!"
	err = file.Save("MyXLSXFile.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

type doc struct {
	file  *xlsx.File
	sheet *xlsx.Sheet
}

func (d *doc) init() error {
	d.file = xlsx.NewFile()
	sheet, err := d.file.AddSheet("Sheet1")
	if err != nil {
		return err
	}
	d.sheet = sheet
	return nil
}

func (d *doc) addRow(s []string) {
	var cell *xlsx.Cell
	row := d.sheet.AddRow()
	for _, x := range s {
		cell = row.AddCell()
		cell.Value = x
	}
}

func (d *doc) save(fileName string) error {
	err := d.file.Save(fileName)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	fileName := "./temp.csv"
	outName := "./test.xlsx"

	csvFile, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	scanner := bufio.NewScanner(csvFile)
	header := []string{}
	out := doc{}

	err = out.init()
	if err != nil {
		panic(err)
	}

	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ",")
		if len(s) < 2 {
			continue
		}
		if len(header) == 0 {
			header = s
			continue
		}
		out.addRow(s)
		//ip, port := s[0], s[1]
		//fmt.Println(s)
	}

	err = out.save(outName)
	if err != nil {
		panic(err)
	}
	//fmt.Println(header)
}

// func main2() {
// 	// url := "https://buy.line.me/manager/admin/downloadReport?exportFile=1&key=orderReport:all_2019-05-11_createtime.csv"
// 	//fileName := "./orderReport_all_2019-05-11_createtime.csv"
// 	fileName := "./temp.csv"

// 	// err := downloadFile(fileName, url)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	csvFile, err := os.Open(fileName)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer csvFile.Close()

// 	reader := csv.NewReader(csvFile)
// 	reader.FieldsPerRecord = -1 // optional
// 	reader.TrimLeadingSpace = true

// 	lines, err := reader.ReadAll()
// 	if err != nil {
// 		panic(err)
// 	}

// 	i := 0
// 	for _, line := range lines {
// 		data := CsvLine{
// 			Column1: line[0],
// 			Column2: line[1],
// 		}
// 		if i > 10 {
// 			break
// 		}
// 		fmt.Println(data.Column1 + " " + data.Column2)
// 	}

// 	// r := csv.NewReader(csvFile)

// 	// for {
// 	// 	record, err := r.Read()
// 	// 	if err == io.EOF {
// 	// 		break
// 	// 	}
// 	// 	if err != nil {
// 	// 		log.Fatal(err)
// 	// 	}

// 	// 	fmt.Println(record)
// 	// }
// 	// for i :=0 ; i < 10; i++{
// 	// 	record, err := r.Read()
// 	// 	if err == io.EOF {
// 	// 		break
// 	// 	}
// 	// 	if err != nil {
// 	// 		log.Fatal(err)
// 	// 	}

// 	// 	fmt.Println(record)
// 	// }

// 	fmt.Println("hello")
// }
