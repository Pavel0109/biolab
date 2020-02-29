package main
import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

func main() {
	xmlFile, err := os.Open("data/csf_metabolites.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)

	level := 0 // level of tegs

	for {
		var buf string
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

			if tok.Name.Local == "name" {
				if err := decoder.DecodeElement(&buf, &tok); err != nil {
					fmt.Println("error", err)
				}
				fmt.Println(level, buf)
				level--
			}

		case xml.EndElement:
			level--
		}
	}
	//fmt.Println(result)
}
