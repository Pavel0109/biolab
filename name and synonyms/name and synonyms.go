package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
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
func main() {
	//configurets := []string{"iupac_name", "average_molecular_weight", "status"}

	xmlFile, err := os.Open("data/hmdb_metabolites.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	isinsyn := 0
	var result []NameAndSynonyms
	ResultLength := 0
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
			if tok.Name.Local == "metabolite" {
				var cs NameAndSynonyms
				result = append(result, cs)
				ResultLength++
			}
			if tok.Name.Local == "name" && level == 3 {
				var buf string
				if err := decoder.DecodeElement(&buf, &tok); err != nil {
					fmt.Println("error", err)
				}
				level--
				result[ResultLength-1].Name = buf

			}
			if tok.Name.Local == "accession" && level == 3 {
				var buf string
				if err := decoder.DecodeElement(&buf, &tok); err != nil {
					fmt.Println("error", err)
				}
				level--
				result[ResultLength-1].Accesid = buf

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

				result[ResultLength-1].Synonyms = append(result[ResultLength-1].Synonyms, buf)
			}
		case xml.EndElement:
			level--
			if tok.Name.Local == "synonyms" {
				isinsyn = 0
			}
		}
	}


	file, err := os.Create("name and synonyms/result_name_synonyms.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	//err = writer.Write(configurets)

	checkError("Cannot write to file", err)
	for _, current := range result {
			err := writer.Write(append([]string{current.Accesid,current.Name},current.Synonyms...))
			checkError("Cannot write to file", err)

	}

}
