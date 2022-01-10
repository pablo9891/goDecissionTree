package datahandler

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	BooleanType   = iota
	NumericalType // Includes float and integer
	CharacterType // Includes chars and strings
	UndefinedType
)

type Field struct {
	Val interface{}
}

type Column struct {
	Name     string
	DataType int
	Fields   []*Field
	Idx      int
}

type Data struct {
	Columns      map[string]*Column
	NumberOfRows int
	NumberOfCol  int
}

func isNumeric(s interface{}) bool {
	_, err := strconv.ParseFloat(fmt.Sprintf("%v", s), 64)
	return err == nil
}

func isBoolean(s interface{}) bool {
	_, err := strconv.ParseBool(fmt.Sprintf("%v", s))
	return err == nil
}

func isCharacter(s interface{}) bool {
	return !(isNumeric(s) || isBoolean(s))
}

func tryBoolean(c *Column) error {
	i := 0
	for i < len(c.Fields) {
		if !isBoolean(c.Fields[i].Val) {
			return fmt.Errorf("found a not boolean")
		}
		i++
	}
	return nil
}

func tryNumeric(c *Column) error {
	i := 0
	for i < len(c.Fields) {
		if !isNumeric(c.Fields[i].Val) {
			return fmt.Errorf("found a not numeric")
		}
		i++
	}
	return nil
}

func tryCharacter(c *Column) error {
	i := 0
	for i < len(c.Fields) {
		if !isCharacter(c.Fields[i].Val) {
			return fmt.Errorf("found a not character")
		}
		i++
	}
	return nil
}

func getColumnType(t string, c *Column) (int, error) {
	result := UndefinedType
	var err error

	if t == "bool" {
		err = tryBoolean(c)
		if err == nil {
			result = BooleanType
		}
	} else if t == "float64" {
		err = tryNumeric(c)
		if err == nil {
			result = NumericalType
		}
	} else if t == "string" {
		err = tryCharacter(c)
		if err == nil {
			result = CharacterType
		}
	}

	return result, err
}

func inferType(c *Column) {
	types := [3]string{"bool", "float64", "string"}

	i := 0
	for i < len(types) {
		t := types[i]
		result, err := getColumnType(t, c)
		if err == nil {
			c.DataType = result
			break
		}
		i++
	}
}

func ReadCsvColumns(path string) (map[string]*Column, int, error) {
	columns := make(map[string]*Column)
	readFile, err := os.Open(path)
	colNames := make([]string, 0)

	if err != nil {
		return nil, -1, fmt.Errorf("failed to open file")
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	i := 0
	numberOfFoundColumns := 0
	for fileScanner.Scan() {
		line := fileScanner.Text()
		fields := strings.Split(line, ",")
		for idx, d := range fields {
			if i == 0 {
				colName := strings.TrimSpace(d)
				colNames = append(colNames, colName)
				columns[colName] = &Column{colName, UndefinedType, make([]*Field, 0), idx}
				numberOfFoundColumns++
			} else {
				if (idx + 1) > numberOfFoundColumns {
					fmt.Println((idx + 1), " ", numberOfFoundColumns)
					return nil, -1, fmt.Errorf("wrong format at line: %d", (i + 1))
				}
				columns[colNames[idx]].Fields = append(columns[colNames[idx]].Fields, &Field{strings.TrimSpace(d)})
			}
		}
		i++
	}
	return columns, i - 1, nil
}

func SplitData(data *Data, col *Column, val interface{}) *Data {
	var i, j int
	isFirstIteration := true
	j = 0
	newData := &Data{Columns: make(map[string]*Column), NumberOfRows: 0, NumberOfCol: data.NumberOfCol - 1}
	for columnName, column := range data.Columns {
		if columnName != col.Name {
			if _, ok := newData.Columns[columnName]; !ok {
				newData.Columns[columnName] = &Column{Name: data.Columns[columnName].Name, DataType: data.Columns[columnName].DataType, Fields: make([]*Field, 0), Idx: data.Columns[columnName].Idx}
			}
			i = 0
			for _, f := range column.Fields {
				if col.Fields[i].Val == val {
					newData.Columns[columnName].Fields = append(newData.Columns[columnName].Fields, &Field{Val: f.Val})
					j++
				}
				i++
			}
			if isFirstIteration {
				newData.NumberOfRows = j
				isFirstIteration = false
			}
		}
	}

	return newData
}

func Parse(path string) (*Data, error) {
	columns, rows, err := ReadCsvColumns(path)

	if err != nil {
		return nil, err
	}

	for _, c := range columns {
		inferType(c)
	}
	return &Data{columns, rows, len(columns)}, err
}
