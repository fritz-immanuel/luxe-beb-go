package templatehtml

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	preload []string
)

func TemplateDefault() string {
	file, err := os.Open("html/DefaultTemplate.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	str, err := ioutil.ReadAll(file)
	return string(str)
}

func TemplateInvoice() string {
	file, err := os.Open("html/BillingTemplateNonBS.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	str, err := ioutil.ReadAll(file)
	return string(str)
}

func ReservationInvoiceTemplate() string {
	file, err := os.Open("html/ReservationInvoice.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	str, err := ioutil.ReadAll(file)
	return string(str)
}
