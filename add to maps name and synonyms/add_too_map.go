package  main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type NameAndSynonyms struct {
	Accesid string
	Name string
	Synonyms []string
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
func addtomap()  {

}
func main()  {

	result := make(map[string] int)
	xmlFile, err := os.Open("data/hmdb_metabolites.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	isinsyn := 0
	level := 0 // level of tegs
	for {
		tok, tokenErr := decoder.Token()
		if tokenErr != nil && tokenErr != io.EOF {
			fmt.Println("error happend", tokenErr)
			break
		} else if tokenErr == io.EOF {
			break
		}

		switch tok := tok.(type) {
		case xml.StartElement:
			level++
			if tok.Name.Local == "name" && level == 3 {
				var buf string
				if err := decoder.DecodeElement(&buf, &tok); err != nil {
					fmt.Println("error", err)
				}
				level--
				result[buf]++

			}
			if tok.Name.Local == "synonyms" {
				isinsyn = 1
			}
			if tok.Name.Local == "synonym" && isinsyn == 1 && level == 4{
				var buf string
				if err := decoder.DecodeElement(&buf, &tok); err != nil {
					fmt.Println("error", err)
				}
				level--

				result[buf]++
			}
		case xml.EndElement:
			level--
			if tok.Name.Local == "synonyms" {
				isinsyn = 0
			}
		}
	}


	file, err := os.Create("name and synonyms/result_name_synonyms_2.csv")
	checkError("Cannot create file", err)
	defer file.Close()
	//fmt.Println(result)
	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()
	checkError("Cannot write to file", err)
	for i,j := range result{
		writer.Write([]string{i,strconv.Itoa(j)})
		checkError("Cannot write to file", err)

	}
}





