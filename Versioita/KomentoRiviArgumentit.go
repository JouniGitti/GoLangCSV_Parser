package Versioita

import (
	"fmt"
	"os"
)

func main() {

lueKomentorivinArgumentit()

}

func lueKomentorivinArgumentit() {

	granulariteetti := os.Args
	tietotyypinValinta := os.Args[1:]

	fmt.Println(granulariteetti)
	fmt.Println(tietotyypinValinta)

}