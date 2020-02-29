package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"encoding/csv"
)

type AccesAndSecond struct {
	First string
	Second []string
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func deleteDublicate(a string, b[]string) []string  {
	result := []string{a}
	for _, i := range b{
		if i != a{
			result = append(result,i)
		}
	}
	return result
}
func main() {

	//configurets := []string{"iupac_name", "average_molecular_weight", "status"}

	xmlFile, err := os.Open("data/hmdb_metabolites.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)

	//result:= new ([]map[string]string)
	//current := make(map[string]string)
	//var current  []string
	//var result [][]string
	var result []AccesAndSecond
	ResultLength := 0
	IsInSecond := 0
	for {
		//len := 0
		tok, tokenErr := decoder.Token()
		if tokenErr != nil && tokenErr != io.EOF {
			fmt.Println("error happend", tokenErr)
			break
		} else if tokenErr == io.EOF {
			break
		}

		switch tok := tok.(type) {
		case xml.StartElement:
			if tok.Name.Local == "metabolite" {
				var cs AccesAndSecond
				result = append(result, cs)
				ResultLength++
			}
			if tok.Name.Local == "accession" {
				var buf string
				if err := decoder.DecodeElement(&buf, &tok); err != nil {
					fmt.Println("error", err)
				}
				if IsInSecond == 0{
					result[ResultLength-1].First = buf
				}
				if IsInSecond == 1{
					result[ResultLength-1].Second = append(result[ResultLength-1].Second,buf)
				}
			}
			if tok.Name.Local == "secondary_accessions" {
				IsInSecond = 1
			}
		case xml.EndElement:
			if tok.Name.Local == "secondary_accessions" {
				IsInSecond = 0
			}
		}
	}
	fmt.Println(result)

	file, err := os.Create("accession and second accession/result.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	//err = writer.Write(configurets)

	checkError("Cannot write to file", err)
	for _, current := range result {
		if len(current.Second) > 0 {
			err := writer.Write(deleteDublicate(current.First,current.Second))
			checkError("Cannot write to file", err)
		}
	}

}
