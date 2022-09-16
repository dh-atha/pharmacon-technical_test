package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println(nominal(12))
	fmt.Println(nominal(2012))
	fmt.Println(nominal(999999))
	fmt.Println(nominal(110110))
	fmt.Println(nominal(10001))
	fmt.Println(nominal(111000))
	fmt.Println(nominal(110000))
}

var data = map[string]string{
	"1":  "satu",
	"2":  "dua",
	"3":  "tiga",
	"4":  "empat",
	"5":  "lima",
	"6":  "enam",
	"7":  "tujuh",
	"8":  "delapan",
	"9":  "sembilan",
	"10": "sepuluh",
	"11": "sebelas",
	"12": "dua belas",
	"13": "tiga belas",
	"14": "empat belas",
	"15": "lima belas",
	"16": "enam belas",
	"17": "tujuh belas",
	"18": "delapan belas",
	"19": "sembilan belas",
}

func nominal(amount int) (result string) {
	result = strconv.Itoa(amount)
	// pisahkan 3 digit pertama (ribuan) dan 3 digit kedua (satuan)
	var ribuan, satuan string
	if len(result) > 3 {
		ribuan = result[:len(result)-3]
		satuan = result[len(result)-3:]
	} else {
		satuan = result
	}

	result = ""
	// convert ribuan
	if len(ribuan) == 1 { // case seribu atau dua ribu atau seterusnya
		if data[string(ribuan[0])] == "satu" {
			result += "seribu "
		} else {
			result += fmt.Sprint(data[string(ribuan[0])], " ribu ")
		}
	} else if len(ribuan) >= 2 { // case ratusan atau puluhan ribu
		result += fmt.Sprint(convertNominal(ribuan), " ribu ")
	}

	// convert satuan
	if len(satuan) == 1 { // case satu atau dua atau seterusnya
		result = data[string(satuan[0])]
	} else if len(satuan) >= 2 { // case ratusan atau puluhan
		result += fmt.Sprint(convertNominal(satuan))
	}
	return
}

func convertNominal(input string) (result string) {
	for len(input) != 0 {
		if data[string(input[0])] == "" {
			input = input[1:]
			continue
		}
		if len(input) == 3 {
			if data[string(input[0])] == "satu" {
				result += "seratus "
			} else {
				result += fmt.Sprint(data[string(input[0])], " ratus ")
			}
		} else if len(input) == 2 {
			if data[string(input[0])] == "satu" {
				result += data[input]
				break
			} else {
				result += fmt.Sprint(data[string(input[0])], " puluh ")
			}
		} else {
			result += data[string(input[0])]
		}
		input = input[1:]
	}
	return
}
