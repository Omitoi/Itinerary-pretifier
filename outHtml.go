package main

import (
	"regexp"
	"strconv"
	"strings"
)

func getOutputStringHTML(input string) string {
	input = "<!DOCTYPE html><html><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>Flight Itinerary</title><style>@media screen and (max-width: 600px) {.container {width: 100% !important;}.content {padding: 15px !important;}.button{width: 100% !important;display: block !important;}}</style></head><body style=\"margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #f4f4f4;\"><table role=\"presentation\" width=\"100%\" cellspacing=\"0\" cellpadding=\"0\" border=\"0\" style=\"background-color: #f4f4f4;\"><tr><td align=\"center\"><table role=\"presentation\" class=\"container\" width=\"600\" cellspacing=\"0\" cellpadding=\"0\" border=\"0\" style=\"max-width: 600px; background-color: #ffffff; margin: 20px auto; border: 1px solid #ddd; border-radius: 5px;\"><tr><td align=\"center\" style=\"padding: 20px; background-color: #007bff; color: #ffffff; font-size: 24px; font-weight: bold; border-top-left-radius: 5px; border-top-right-radius: 5px;\">Flight Itinerary</td></tr><tr><td class=\"content\" style=\"padding:10px 30px; text-align: left; font-size: 16px; color: #333333;\"><p>" +
		input + "<p style=\"text-align: center;\"><a href=\"https://www.example.com\" class=\"button\" style=\"background-color: #007bff; color: #ffffff; text-decoration: none; padding: 15px 30px; border-radius: 5px; display: inline-block; font-size: 18px;\">See your Itinerary</a></p><p>Thank you for travelling with us,</p><p>Anywhere Holidays Team</p></td></tr><tr><td align=\"center\" style=\"padding: 20px; background-color: #f4f4f4; color: #777777; font-size: 14px; border-bottom-left-radius: 5px; border-bottom-right-radius: 5px;\">&copy; 2025 Anywhere Holidays, Inc. All rights reserved. <br><p style=\"text-align: center;\">This email has been sent to you because you've reserved holidays with us</p></td></tr></table></td></tr></table></body></html>"
	input = placeICAONamesHTML(input)
	input = placeIATANamesHTML(input)
	input = placeTimesHTML(input)
	input = replaceLineBreaks(input)
	input = cleanUpDoubleWhiteSpaces(input)
	input = replaceLineBreaksHTML(input)
	return input
}

func replaceLineBreaksHTML(input string) string {
	//Run through the input
	reHTMLSpacing := regexp.MustCompile(`\n`)
	input = reHTMLSpacing.ReplaceAllString(input, "</p><p>")
	return input
}

func placeICAONamesHTML(input string) string {
	//Find pattern ##XXXX
	re := regexp.MustCompile(`##.{4}`)

	return re.ReplaceAllStringFunc(input, func(match string) string {
		airportCode := match[2:]

		//Match it with an airport
		for _, airport := range airports {
			if airport.ICAO_Code == airportCode {
				link := "<a href=\"https://www.google.com/maps/search/?api=1&query=" +
					strings.ReplaceAll(airport.Name, " ", "+") + "\" target=\"_blank\">" +
					airport.Name + "</a>"
				return link
			}
		}

		//Keep as is if there is not match
		return match
	})
}

func placeIATANamesHTML(input string) string {
	//Find pattern #XXX
	re := regexp.MustCompile(`#.{3}`)

	return re.ReplaceAllStringFunc(input, func(match string) string {
		airportCode := match[1:]

		//Match it with an airport
		for _, airport := range airports {
			if airport.IATA_Code == airportCode {
				link := "<a href=\"https://www.google.com/maps/search/?api=1&query=" +
					strings.ReplaceAll(airport.Name, " ", "+") + "\" target=\"_blank\">" +
					airport.Name + "</a>"
				return link
			}
		}

		//Keep as is if not match
		return match
	})
}

func placeTimesHTML(input string) string {
	//Detect pattern (D|T12|T24)(NNNN-NN-NNTNN:NN(-NN:00|+NN:00|Z)) *N - any number
	re := regexp.MustCompile(`(?:D|T12|T24)\(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}(?:[−+-]\d{2}:00|Z)\)`)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		//D format - converts into DD-mmm-YYYY
		if match[0] == 'D' {
			//Find values
			year := match[2:6]
			month := match[7:9]
			day := match[10:12]

			if len(match) > 22 {
				//Get offset and positivity
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
			return "<strong>" + day + " " + month + " " + year + "</strong>"

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
					return "<em>" + prefixZero(hours) + ":" + minutes + "PM (" + outputOffset + ":00)</em>"
				} else if hours == 24 || hours == 0 {
					hours = 12
				}
				return "<em>" + prefixZero(hours) + ":" + minutes + "AM (" + outputOffset + ":00)</em>"
			} else { //Zulu offset
				if hours > 13 || hours < 24 {
					hours -= 12
					if hours < 0 {
						hours = -hours
					}
					return "<em>" + prefixZero(hours) + ":" + minutes + "PM (+00:00)</em>"
				} else if hours == 24 || hours == 0 {
					hours = 12
				}
				return "<em>" + prefixZero(hours) + ":" + minutes + "AM (+00:00)</em>"
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

				return "<em>" + prefixZero(hours) + ":" + minutes + " (" + outputOffset + ":00)</em>"
			} else { //Zulu offset
				return "<em>" + prefixZero(hours) + ":" + minutes + " (+00:00)</em>"
			}
		}
		return match
	})
}
