package main

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
	//Muuttujat määritellään. Laustavasti kovakoodattuina. Alapuolella komentoriviltä toimivina.
	var granulariteetti string
	var tietotyyppi string

	type vuosiLukema struct {
		vuosi string
		tulos float64
	}

	type kuukausiLukema struct {
		vuosi    string
		kuukausi string
		tulos    float64
	}

	type paivalukema struct {
		vuosi    string
		kuukausi string
		paivays  string
		tulos    float64
	}
	//Komentoriviargumentit tallennetaan muuttujiin
	//granulariteetti = os.Args[1]
	//tietotyyppi = os.Args[2]

	granulariteetti = "month" // nyt vain koodin testausmielessä vakioita, jatkossa komentoriviargumentteina
	tietotyyppi = "sum"

	//Haarautuminen tietotyypin perusteella
	switch tietotyyppi {
	case "sum":
		Summa("sum.csv", granulariteetti)
	case "cumulative":
		KumulatiivinenSumma("cumulative.csv", granulariteetti)
	case "both":
		//Summa("sum.csv", granulariteetti) + Summa("cumulative.csv", granulariteetti)
		//lasketaan molemmat
	}
}

func KumulatiivinenSumma(tiedostonNimi string, granulariteetti string) float64 {
	var valiSumma float64
	var aikaLeimanPituus int

	//Luodaan hashmap aikamerkkijonoihin perustuen
	tallennetutDatat := make(map[string]float64)

	// Tiedosto avataan
	csvfile, err := os.Open(tiedostonNimi)
	if err != nil {
		log.Fatalln("cumulative.csv tiedostoa ei löytynyt", err)
	}
	// Parseroidaan tiedosto
	r := csv.NewReader(csvfile)

	// Luetaan ensimmäinen rivi, saadaan ensimmäinen aikaleima
	rivi, err := r.Read()
	ensimmainenRivi := string([]rune(rivi[0]))
	fmt.Println(ensimmainenRivi)

	//ekaAikaMerkkiJono tallentaa ekan rivin aika-arvon
	ekaAikaMerkkijono := string(ensimmainenRivi[0:4])
	ekaMjono := strings.Split(rivi[0], ";")

	//Muutetaan merkkijonon lukema-stringi float64:ksi
	ekaLukema, _ := strconv.ParseFloat(ekaMjono[1], 64)

	//Luetaan toinen rivi:
	rivi, err = r.Read()
	toinenRivi := string([]rune(rivi[0]))
	fmt.Println(toinenRivi)

	//ToinenAikaMerkkiJono tallentaa rivin aika-arvon
	toinenAikaMerkkijono := string(toinenRivi[0:4])
	toinenMjono := strings.Split(rivi[0], ";")

	//Muutetaan toisen merkkijonon lukemastringi float64:ksi
	toinenLukema, _ := strconv.ParseFloat(toinenMjono[1], 64)

	//Todellinen lukema saadaan vähentämällä perättäiset lukemat, koska lukemat ovat kumulatiivisia.
	if ekaAikaMerkkijono == toinenAikaMerkkijono {
		valiSumma = + (toinenLukema - ekaLukema)

		//Vaihdetaan muuttujat, että seuraava lukema voidaan laskea
		ekaAikaMerkkijono = toinenAikaMerkkijono
		ekaLukema = toinenLukema
	}
	// Tehdään granulariteettijako: vuoden, kuukauden tai päivän mukaan.
	// Käytetään muuttujana CSV:n aikaleiman pituutta. Aikaleimat 2017-01-01 tyyppisesti
	if granulariteetti == "year" {
		aikaLeimanPituus = 4
	} else if granulariteetti == "month" {
		aikaLeimanPituus = 7
	} else if granulariteetti == "day" {
		aikaLeimanPituus = 10
	}
	for {
		// Luetaan uusi rivi tiedostosta
		rivi, err := r.Read()
		//Tiedoston päättyminen huomioidaan päättämällä for rakenne
		if err == io.EOF {
			break
		}
		//Huomioidaan mahdollinen virhe tiedoston lukemisessa
		if err != nil {
			log.Fatal(err)
		}
		toinenMjono := strings.Split(rivi[0], ";")
		//Muutetaan merkkijonon lukemastringi float64:ksi
		toinenLukema, _ := strconv.ParseFloat(toinenMjono[1], 64)
		toinenRivi := string([]rune(rivi[0]))
		toinenAikaMerkkijono := string(toinenRivi[0:aikaLeimanPituus])
		// Lisätään lukema välisummaan, ensin tarkistetaan lukeman etumerkki
		// fmt.Println("Toinen lukema nyt: ", toinenLukema) //OK

		if toinenAikaMerkkijono != ekaAikaMerkkijono {
			fmt.Println("Täällä ") //OK
			valiSumma += (toinenLukema - ekaLukema)
			tallennetutDatat[ekaAikaMerkkijono] = valiSumma
			fmt.Println("aikamerkkijono nyt ", ekaAikaMerkkijono)
			fmt.Println("Välisumma nyt ", tallennetutDatat[ekaAikaMerkkijono])
			ekaAikaMerkkijono = toinenAikaMerkkijono
			ekaLukema = toinenLukema
			valiSumma = 0
		}
		//Käydään läpi csv rivi riviltä
		if rivi[0][0:aikaLeimanPituus] == ekaAikaMerkkijono {
			// Erotetaan luetusta rivistä mittarin lukema
			mjono := strings.Split(rivi[0], ";")
			//Virhetarkistus
			if err != nil {
				fmt.Println(err)
			}
			//Muutetaan merkkijonon lukemastringi float64:ksi
			ekaLukema, _ := strconv.ParseFloat(mjono[0], 64)
			// Tarkistetaan onko mittarin lukema nollattu
			erotus := toinenLukema - ekaLukema
			fmt.Println("erotus", erotus)
			if erotus > 10 {
				fmt.Println("ekaAikaMerkkijono", ekaAikaMerkkijono)
				//Nollauksen jälkeen ensimmäisen päivän lukema on sama kuin välisumma normaalisti
				valiSumma = valiSumma + ekaLukema
			} else {
				// Lisätään lukema välisummaan
				valiSumma += ekaLukema
			}
			//Tallennetaan välisumma Map:iin aikaleiman perusteella
			tallennetutDatat[ekaAikaMerkkijono] = valiSumma

		} else {
			//Jos luettu aikaleima on eri, kuin aiempi aikaleima, vaihdetaan uusi aikaleima
			//tallennetutDatat[aikaMerkkijono] = valiSumma
			//aikaMerkkijono = rivi[0][0:aikaLeimanPituus]
			//// fmt.Println("aikamerkkijono on nyt: ", aikaMerkkijono)
			//Nollataan välisumma, jotta uudella aikaleimalla laskenta alkaa alusta
			valiSumma = 0
			//// fmt.Println("viimeinen välisumma elsen jälkeen: ", valiSumma)
		}
		// Tulostetaan ruudulle lasketut arvot. Tämä ei vielä toimi.
		for ekaAikaMerkkijono, valiSumma := range tallennetutDatat {
			fmt.Println(ekaAikaMerkkijono, "  ", valiSumma)

			// Tässä koetulosteet, jotka jo osin toimivat
			fmt.Println("2017-02-01", tallennetutDatat["2017-02-01"])
			fmt.Println("2017-02-02", tallennetutDatat["2017-02-02"])
			fmt.Println("2017-02-03", tallennetutDatat["2017-02-03"])
			fmt.Println("2017-01", tallennetutDatat["2017-01"])
			fmt.Println("2017-02", tallennetutDatat["2017-01"])
			fmt.Println("2017-03", tallennetutDatat["2017-01"])
			fmt.Println("2017-04", tallennetutDatat["2017-04"])
			fmt.Println("2017-05", tallennetutDatat["2017-05"])
			fmt.Println("2017-06", tallennetutDatat["2017-06"])
			fmt.Println("2017-07", tallennetutDatat["2017-07"])
			fmt.Println("2017-08", tallennetutDatat["2017-08"])
			fmt.Println("2017-09", tallennetutDatat["2017-09"])
			fmt.Println("2017-10", tallennetutDatat["2017-10"])
			fmt.Println("2017-11", tallennetutDatat["2017-11"])
			fmt.Println("2017-12", tallennetutDatat["2017-12"])
			fmt.Println("2018-01", tallennetutDatat["2018-01"])
			fmt.Println("2018-02", tallennetutDatat["2018-01"])
			fmt.Println("2018-03", tallennetutDatat["2018-01"])
			fmt.Println("2018-04", tallennetutDatat["2018-04"])
			fmt.Println("2018-05", tallennetutDatat["2018-05"])
			fmt.Println("2018-06", tallennetutDatat["2018-06"])
			fmt.Println("2018-07", tallennetutDatat["2018-07"])
			fmt.Println("2018-08", tallennetutDatat["2018-08"])
			fmt.Println("2018-09", tallennetutDatat["2018-09"])
			fmt.Println("2018-10", tallennetutDatat["2018-10"])
			fmt.Println("2018-11", tallennetutDatat["2018-11"])
			fmt.Println("2018-12", tallennetutDatat["2018-12"])
			fmt.Println("2017-01", tallennetutDatat["2019-01"])
			fmt.Println("2019-02", tallennetutDatat["2019-01"])
			fmt.Println("2019-03", tallennetutDatat["2019-01"])
			fmt.Println("2019-04", tallennetutDatat["2019-04"])
			fmt.Println("2019-05", tallennetutDatat["2019-05"])
			fmt.Println("2019-06", tallennetutDatat["2019-06"])
			fmt.Println("2019-07", tallennetutDatat["2019-07"])
			fmt.Println("2019-08", tallennetutDatat["2019-08"])
			fmt.Println("2019-09", tallennetutDatat["2019-09"])
			fmt.Println("2019-10", tallennetutDatat["2019-10"])
			fmt.Println("2019-11", tallennetutDatat["2019-11"])
			fmt.Println("2019-12", tallennetutDatat["2019-12"])
			fmt.Println("2017", tallennetutDatat["2017"])
			fmt.Println("2018", tallennetutDatat["2018"])
			fmt.Println("2019", tallennetutDatat["2019"])

		}
	}

	return valiSumma //turha return tässä vaiheessa...
}

func Summa(tiedostonNimi string, granulariteetti string) float64 {
	var valiSumma float64
	var aikaLeimanPituus int
	//Luodaan hashmap aikamerkkijonoihin perustuen
	tallennetutDatat := make(map[string]float64)
	// Tiedosto avataan
	csvfile, err := os.Open(tiedostonNimi)
	if err != nil {
		log.Fatalln("sum.csv tiedostoa ei löytynyt", err)
	}
	// Parseroidaan tiedosto NewReaderillä
	r := csv.NewReader(csvfile)
	// Luetaan ensimmäinen rivi, että saadaan ensimmäinen aikaleima
	// Oletan, että näissä CSV-tiedostoissa datat ovat aikajärjestyksessä.
	rivi, err := r.Read()
	ensimmainenRivi := string([]rune(rivi[0]))
	//aikaMerkkiJono tallentaa ensimmäisen rivin aika-arvon
	aikaMerkkijono := string(ensimmainenRivi[0:aikaLeimanPituus])
	//Summataan ensimmäinen lukema
	mjono := strings.Split(rivi[0], ";")
	//Muutetaan merkkijonon lukemastringi float64:ksi
	yksiLuku, _ := strconv.ParseFloat(mjono[1], 64)
	// Tarkistus, että lukeman etumerkki on oikein
	if yksiLuku < 0 {
		valiSumma = + -yksiLuku
	} else {
		valiSumma = + yksiLuku
	}
	// Tehdään granulariteettijako: vuoden, kuukauden tai päivän mukaan.
	// Käytetään tässä CSV:n aikaleiman pituutta.
	if granulariteetti == "year" {
		aikaLeimanPituus = 4
	} else if granulariteetti == "month" {
		aikaLeimanPituus = 7
	} else if granulariteetti == "day" {
		aikaLeimanPituus = 10
	}
	for {
		// Luetaan uusi rivi tiedostosta
		rivi, err := r.Read()
		//Tiedoston päättyminen huomioidaan päättämällä for rakenne
		if err == io.EOF {
			break
		}
		//Huomioidaan mahdollinen virhe tiedoston lukemisessa
		if err != nil {
			log.Fatal(err)
		}
		mjono := strings.Split(rivi[0], ";")
		//Muutetaan merkkijonon lukemastringi float64:ksi
		yksiLuku, _ := strconv.ParseFloat(mjono[1], 64)
		// Lisätään lukema välisummaan, ensin tarkistetaan lukeman etumerkki
		if yksiLuku < 0 {
			valiSumma += -yksiLuku
		} else {
			valiSumma += yksiLuku
		}
		//Käydään läpi csv rivi riviltä
		if rivi[0][0:aikaLeimanPituus] == aikaMerkkijono {
			// Erotetaan luetusta rivistä mittarin lukema
			mjono := strings.Split(rivi[0], ";")
			//Virhetarkistus
			if err != nil {
				fmt.Println(err)
			}
			//Muutetaan merkkijonon lukemastringi float64:ksi
			yksiLuku, _ := strconv.ParseFloat(mjono[0], 64)
			// Lisätään lukema välisummaan
			valiSumma += yksiLuku
			//Tallennetaan välisumma Map:iin aikaleiman perusteella
			tallennetutDatat[aikaMerkkijono] = valiSumma
		} else {
			//Jos luettu aikaleima on eri, kuin aiempi aikaleima, vaihdetaan uusi aikaleima
			tallennetutDatat[aikaMerkkijono] = valiSumma
			aikaMerkkijono = rivi[0][0:aikaLeimanPituus]
			//// fmt.Println("aikamerkkijono on nyt: ", aikaMerkkijono)
			//Nollataan välisumma, jotta uudella aikaleimalla laskenta alkaa alusta
			valiSumma = 0
			//// fmt.Println("viimeinen välisumma elsen jälkeen: ", valiSumma)
		}
	}

	// Tulostus ruudulle ei vielä toimi

	// Tässä koetulosteet, jotka osin jo toimivat:
	fmt.Println("2017-02-01", tallennetutDatat["2017-02-01"])
	fmt.Println("2017-02-02", tallennetutDatat["2017-02-02"])
	fmt.Println("2017-02-03", tallennetutDatat["2017-02-03"])
	fmt.Println("2017-01", tallennetutDatat["2017-01"])
	fmt.Println("2017-02", tallennetutDatat["2017-01"])
	fmt.Println("2017-03", tallennetutDatat["2017-01"])
	fmt.Println("2017-04", tallennetutDatat["2017-04"])
	fmt.Println("2017-05", tallennetutDatat["2017-05"])
	fmt.Println("2017-06", tallennetutDatat["2017-06"])
	fmt.Println("2017-07", tallennetutDatat["2017-07"])
	fmt.Println("2017-08", tallennetutDatat["2017-08"])
	fmt.Println("2017-09", tallennetutDatat["2017-09"])
	fmt.Println("2017-10", tallennetutDatat["2017-10"])
	fmt.Println("2017-11", tallennetutDatat["2017-11"])
	fmt.Println("2017-12", tallennetutDatat["2017-12"])
	fmt.Println("2018-01", tallennetutDatat["2018-01"])
	fmt.Println("2018-02", tallennetutDatat["2018-01"])
	fmt.Println("2018-03", tallennetutDatat["2018-01"])
	fmt.Println("2018-04", tallennetutDatat["2018-04"])
	fmt.Println("2018-05", tallennetutDatat["2018-05"])
	fmt.Println("2018-06", tallennetutDatat["2018-06"])
	fmt.Println("2018-07", tallennetutDatat["2018-07"])
	fmt.Println("2018-08", tallennetutDatat["2018-08"])
	fmt.Println("2018-09", tallennetutDatat["2018-09"])
	fmt.Println("2018-10", tallennetutDatat["2018-10"])
	fmt.Println("2018-11", tallennetutDatat["2018-11"])
	fmt.Println("2018-12", tallennetutDatat["2018-12"])
	fmt.Println("2017-01", tallennetutDatat["2019-01"])
	fmt.Println("2019-02", tallennetutDatat["2019-01"])
	fmt.Println("2019-03", tallennetutDatat["2019-01"])
	fmt.Println("2019-04", tallennetutDatat["2019-04"])
	fmt.Println("2019-05", tallennetutDatat["2019-05"])
	fmt.Println("2019-06", tallennetutDatat["2019-06"])
	fmt.Println("2019-07", tallennetutDatat["2019-07"])
	fmt.Println("2019-08", tallennetutDatat["2019-08"])
	fmt.Println("2019-09", tallennetutDatat["2019-09"])
	fmt.Println("2019-10", tallennetutDatat["2019-10"])
	fmt.Println("2019-11", tallennetutDatat["2019-11"])
	fmt.Println("2019-12", tallennetutDatat["2019-12"])
	fmt.Println("2017", tallennetutDatat["2017"])
	fmt.Println("2018", tallennetutDatat["2018"])
	fmt.Println("2019", tallennetutDatat["2019"])

	return valiSumma
}
