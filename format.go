package main

import (
	"regexp"
	"strconv"
	"strings"
)

func getOutputString(input string) string {
	input = placeICAONameCities(input)
	input = placeICAONames(input)
	input = placeIATANameCities(input)
	input = placeIATANames(input)
	input = placeTimes(input)
	input = replaceLineBreaks(input)
	input = cleanUpDoubleWhiteSpaces(input)
	return input
}

func validOffset(input int) bool {
	if input < -12 || input > 14 {
		return false
	}
	return true
}

func prefixZero(input int) string {
	if input < 10 && input > -10 {
		if input < 0 {
			input = -input
			return "-0" + strconv.Itoa(input)
		}
		return "0" + strconv.Itoa(input)
	}
	return strconv.Itoa(input)
}

func isAlphaNumeric(input rune) bool {
	return (input >= 'A' && input <= 'Z') || // Uppercase letters
		(input >= 'a' && input <= 'z') || // Lowercase letters
		(input >= '0' && input <= '9') // Numbers
}

func replaceLineBreaks(input string) string {
	//Stupid Windows
	reWhiteSpace := regexp.MustCompile(`\r\n`)
	input = reWhiteSpace.ReplaceAllString(input, "\n")
	//Run through the input
	return strings.Map(func(r rune) rune {

		//Look for a control character
		if r == '\v' || r == '\f' || r == '\r' {
			//Replace with \n
			return '\n'
		}
		return r
	}, input)
}

func cleanUpDoubleWhiteSpaces(input string) string {
	//Remove consecutive newlines
	reWhiteSpace := regexp.MustCompile(`\n\n+`)
	input = reWhiteSpace.ReplaceAllString(input, "\n\n")
	//Remove consecutive spaces
	/*
		reNewLine := regexp.MustCompile(` +`)
		input = reNewLine.ReplaceAllString(input, " ")
	*/
	//Remove spaces as a first character in a newline
	/*
		reNewLineSpace := regexp.MustCompile(`\n `)
		input = reNewLineSpace.ReplaceAllString(input, "\n")
	*/
	//Trim for a space in the first line and the end
	input = strings.TrimSpace(input)
	return input
}

func placeICAONames(input string) string {
	//Find pattern ##XXXX
	re := regexp.MustCompile(`##[A-Z]{4}`)

	return re.ReplaceAllStringFunc(input, func(match string) string {
		//Check previous or following character for failure exceptions
		position := strings.Index(input, match)
		if position > 0 && isAlphaNumeric(rune(input[position-1])) {
			return match
		}
		if position+6 < len(input) && isAlphaNumeric(rune(input[position+6])) {
			return match
		}

		airportCode := match[2:]

		//Match it with an airport
		for _, airport := range airports {
			if airport.ICAO_Code == airportCode {
				return airport.Name
			}
		}

		//Keep as is if there is not match
		return match
	})
}

func placeICAONameCities(input string) string {
	//Find pattern *##XXXX
	re := regexp.MustCompile(`\*##[A-Z]{4}`)

	return re.ReplaceAllStringFunc(input, func(match string) string {
		position := strings.Index(input, match)
		if position > 0 && isAlphaNumeric(rune(input[position-1])) {
			return match
		}
		if position+7 < len(input) && isAlphaNumeric(rune(input[position+7])) {
			return match
		}

		airportCode := match[3:]

		//Match it with an airport
		for _, airport := range airports {
			if airport.ICAO_Code == airportCode {
				return airport.Municipality
			}
		}

		//Keep as is if there is not match
		return match
	})
}

func placeIATANames(input string) string {
	//Find pattern #XXX
	re := regexp.MustCompile(`#[A-Z]{3}`)

	return re.ReplaceAllStringFunc(input, func(match string) string {
		//Check previous or following character for failure exceptions
		position := strings.Index(input, match)
		if position > 0 && (input[position-1] == '#' || isAlphaNumeric(rune(input[position-1]))) {
			return match
		}
		if position+4 < len(input) && isAlphaNumeric(rune(input[position+4])) {
			return match
		}

		airportCode := match[1:]

		//Match it with an airport
		for _, airport := range airports {
			if airport.IATA_Code == airportCode {
				return airport.Name
			}
		}

		//Keep as is if not match
		return match
	})
}

func placeIATANameCities(input string) string {
	//Find pattern *#XXX
	re := regexp.MustCompile(`\*#[A-Z]{3}(,|.|\s|\n|\t)`)

	return re.ReplaceAllStringFunc(input, func(match string) string {
		//Check previous or following character for failure exceptions
		position := strings.Index(input, match)
		if position > 0 && isAlphaNumeric(rune(input[position-1])) {
			return match
		}
		if position+5 < len(input) && isAlphaNumeric(rune(input[position+5])) {
			return match
		}

		airportCode := match[2:]

		//Match it with an airport
		for _, airport := range airports {
			if airport.IATA_Code == airportCode {
				return airport.Municipality
			}
		}

		//Keep as is if not match
		return match
	})
}

func placeTimes(input string) string {
	//Detect pattern (D|T12|T24)(NNNN-NN-NNTNN:NN(-NN:00|+NN:00|Z)) *N - any number
	re := regexp.MustCompile(`(?:D|T12|T24)\(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}(?:[−+-]\d{2}:00|Z)\)`)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		//D format - converts into DD-mmm-YYYY
		if match[0] == 'D' {
			//Find values
			year := match[2:6]
			month := match[7:9]
			day := match[10:12]

			//Detect if offset moves the day or not
			// Bonus feature?
			if len(match) > 22 {
				//Get offset, hours and positivity
				offset, _ := strconv.Atoi(match[19:21])
				sign := string(match[18])

				//Give it a minus if there is minus
				if sign == "-" {
					offset = -offset
				}

				//Check for valid offset
				if !validOffset(offset) {
					return match
				}
			}

			//Replace the month
			if monthName, exists := monthMap[month]; exists {
				month = monthName
			}
			return day + " " + month + " " + year

		}

		//T12 format - converts into HH-MM(AM/PM) (+OO:00)
		if match[0:3] == "T12" {
			hours, _ := strconv.Atoi(match[15:17])
			minutes := match[18:20]
			sign := string(match[20])

			//Detect not Zulu offset
			if len(match) > 22 {
				//Get offset
				offset, _ := strconv.Atoi(match[21:23])

				//Give it a minus if there is minus
				if sign == "-" {
					offset = -offset
				}

				//Check for valid offset
				if !validOffset(offset) {
					return match
				}
				outputOffset := prefixZero(offset)

				if hours > 13 && hours < 24 {
					hours -= 12
					if hours < 0 {
						hours = -hours
					}
					return prefixZero(hours) + ":" + minutes + "PM (" + outputOffset + ":00)"
				} else if hours == 24 || hours == 0 {
					hours = 12
				}
				return prefixZero(hours) + ":" + minutes + "AM (" + outputOffset + ":00)"
			} else { //Zulu offset
				if hours > 13 && hours < 24 {
					hours -= 12
					if hours < 0 {
						hours = -hours
					}
					return prefixZero(hours) + ":" + minutes + "PM (+00:00)"
				} else if hours == 24 || hours == 0 {
					hours = 12
				}
				return prefixZero(hours) + ":" + minutes + "AM (+00:00)"
			}
		}

		if match[0:3] == "T24" {
			//T24 format - converts into HH-MM (+OO:00)
			hours, _ := strconv.Atoi(match[15:17])
			minutes := match[18:20]
			sign := string(match[20])
			//Detect not Zulu offset
			if len(match) > 22 {
				//Get offset
				offset, _ := strconv.Atoi(match[21:23])

				//Give it a minus if there is minus
				if sign == "-" {
					offset = -offset
				}

				// Check for valid offset
				if !validOffset(offset) {
					return match
				}
				outputOffset := prefixZero(offset)

				return prefixZero(hours) + ":" + minutes + " (" + outputOffset + ":00)"
			} else { //Zulu offset
				return prefixZero(hours) + ":" + minutes + " (+00:00)"
			}
		}
		return match
	})
}
