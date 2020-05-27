package Versioita

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var kokoSumma float64
	var vuosi string

	// Faili avataan
	csvfile, err := os.Open("sum.csv")
	if err != nil {
		log.Fatalln("Tällaista failia ei löytynyt", err)
	}

	// Parseroidaan faili
	r := csv.NewReader(csvfile)

	// Käydään läpi rivi riviltä
	for {
		// Luetaan jokainen rivi failista
		rivi, err := r.Read()
		//Failin päättyminen huomioidaan päättämällä for rakenne
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		//case: year
		vuosi := strings.Split(rivi[0], ";")
		if strings.Contains(rivi[0], "2018") {
			mjono := strings.Split(rivi[0], ";")
			fmt.Println(mjono[1])
			if err != nil {
				fmt.Println(err)

			}
			yksiLuku, _ := strconv.ParseFloat(mjono[1], 64)
			fmt.Println(yksiLuku)
			yhteensa += yksiLuku
		}

	}
	if yhteensa < 0 {
		fmt.Println("Yhteensä: ", -yhteensa)
	} else {
		fmt.Println("Yhteensä: ", yhteensa)
	}
}
