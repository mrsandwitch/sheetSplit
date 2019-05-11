package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

type doc struct {
	file   *xlsx.File
	sheet1 *xlsx.Sheet
	sheet2 *xlsx.Sheet
}

func (d *doc) init(sheetName1 string, sheetName2 string) error {
	d.file = xlsx.NewFile()
	sheet1, err := d.file.AddSheet(sheetName1)
	if err != nil {
		return err
	}
	sheet2, err := d.file.AddSheet(sheetName2)
	if err != nil {
		return err
	}
	d.sheet1 = sheet1
	d.sheet2 = sheet2
	return nil
}

func (d *doc) addRow(s []string, sheetName string) {
	var cell *xlsx.Cell
	var row *xlsx.Row
	if sheetName == "sheet1" {
		row = d.sheet1.AddRow()
	} else {
		row = d.sheet2.AddRow()
	}

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

// Date formatter
const (
	layoutISO    = "2006-01-02"
	layoutCustom = "2006-01-02 15:04:05"
)

func main() {
	files, err := ioutil.ReadDir("./")

	if err != nil {
		log.Fatal(err)
	}
	re := regexp.MustCompile("orderReport_all_(.*)_createtime.csv")
	var fileName string
	var queryDate string

	for _, f := range files {
		match := re.FindStringSubmatch(f.Name())

		if len(match) > 1 {
			fileName = f.Name()
			queryDate = match[1]
			break
		}
	}

	d, err := time.Parse(layoutISO, queryDate)
	if err != nil {
		panic(err)
	}
	dX := d.AddDate(0, 0, -1)
	outName := fmt.Sprintf("ECAC_orderinfo_%d%02d%02dcreatetime.xlsx", dX.Year(), dX.Month(), dX.Day())
	sheetName1 := fmt.Sprintf("%02d%02d", dX.Month(), dX.Day())
	sheetName2 := fmt.Sprintf("%02d%02d", d.Month(), d.Day())

	//- For debug and testing
	//fileName = "./temp.csv"
	//fileName = "./orderReport_all_2019-05-11_createtime.csv"
	//outName := "ECAC_orderinfo_20190510createtime.xlsx"
	//outName := "./test.xlsx"

	csvFile, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	r := csv.NewReader(csvFile)
	r.Comma = ','
	r.FieldsPerRecord = -1
	lines, err := r.ReadAll()
	if err != nil {
		log.Fatalf("error reading all lines: %v", err)
	}

	// Initial xlsx object
	out := doc{}
	err = out.init(sheetName1, sheetName2)
	if err != nil {
		panic(err)
	}

	header := []string{}

	for _, s := range lines {
		if len(s) < 2 {
			continue
		}
		if len(header) == 0 {
			header = s
			out.addRow(header, "sheet1")
			out.addRow(header, "sheet2")
			continue
		}

		dateStr := strings.Trim(s[4], "\"")
		orderTime, err := time.Parse(layoutCustom, dateStr)
		dateStr = strings.Trim(s[11], "\"")
		createTime, err := time.Parse(layoutCustom, dateStr)

		d1, err := time.Parse(layoutISO, "2019-01-01")
		if err != nil {
			fmt.Println(err)
			continue
		}
		d2, err := time.Parse(layoutISO, queryDate) // orderReport_all_2019-05-11_createtime => 2019-05-11
		if err != nil {
			fmt.Println(err)
			continue
		}
		d3 := d2.AddDate(0, 0, -1) // 2019-05-10 =  019-05-11 minus one day

		//select * where E > date '2019-01-01' and L >= date '2019-05-10' and L < date '2019-05-11'
		if orderTime.After(d1) && (createTime.After(d3) || createTime.Equal(d3)) && createTime.Before(d2) {
			out.addRow(s, "sheet1")
		} else if createTime.After(d2) || createTime.Equal(d2) {
			//select * where L >= date '2019-05-11'
			out.addRow(s, "sheet2")
		}
	}

	// Save to new xlsx file
	err = out.save(outName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully output to file %s\n", outName)
	fmt.Printf("Sheet1 has %d rows\n", len(out.sheet1.Rows))
	fmt.Printf("Sheet2 has %d rows\n", len(out.sheet2.Rows))
}
