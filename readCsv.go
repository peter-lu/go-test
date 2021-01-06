package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("10005236_20210105.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			//结束日志
			break
		} else if err != nil {
			//异常日志
			fmt.Println("Error:", err)
			return
		}
		userId, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil || 0 == userId {
			//异常日志
			fmt.Println("Error:", err)
		}
		fmt.Println(userId) // record has the type []string
	}
}
