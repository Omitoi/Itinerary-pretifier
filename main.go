package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Airport struct {
	Name         string
	ISO_Country  string
	Municipality string
	ICAO_Code    string
	IATA_Code    string
	Coordinates  string
}

// Declare a slice of airport structs and variables for bonuses
var airports []Airport
var outputType string

const ( //Terminal color constants
	Red    = "\033[31m"
	Green  = "\033[32m"
	Blue   = "\033[34m"
	Yellow = "\033[33m"
	Reset  = "\033[0m"
)

// Define a map for month number to name mapping
var monthMap = map[string]string{
	"01": "Jan",
	"02": "Feb",
	"03": "Mar",
	"04": "Apr",
	"05": "May",
	"06": "Jun",
	"07": "Jul",
	"08": "Aug",
	"09": "Sep",
	"10": "Oct",
	"11": "Nov",
	"12": "Dec",
}

var options = map[int]string{
	1: "Overwrite",
	2: "Change Name",
	3: "Cancel",
}

func loadFile(path string) (string, error) {
	//Check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s", path)
	}

	//Read file
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("error reading file:  %w", err)
	}

	//Convert to string and send away
	return string(content), nil
}

func validString(input string) bool {
	for _, r := range input {
		if rune(r) > 127 {
			return false
		}
	}
	return true
}

func printHelp() {
	//Help is printed when there are no arguments or -h flag is used
	fmt.Printf("\n%sItinerary usage:%s\n", Blue, Reset)
	fmt.Printf("  go run . -h %s-- Show this help message%s\n", Red, Reset)
	fmt.Printf("  go run . %s[-o]/[-r]%s ./input.txt ./output.txt ./airport-lookup.csv %s-- Proper use of program%s\n", Yellow, Reset, Green, Reset)
	fmt.Printf("\n%sDescription:%s\n", Yellow, Reset)
	fmt.Printf("  -o %s- Automatically overwrites your output without prompting%s\n", Yellow, Reset)
	fmt.Printf("  -r %s- Automatically rewrites your output name without prompting%s\n", Yellow, Reset)
	fmt.Println("  This program prettifies your itinerary that you input.")
	fmt.Println("  Providing invalid input will result in an error message")
	println("")
}

func main() {
	//Define the help flag
	showHelp := flag.Bool("h", false, "Show help information")
	overwrite := flag.Bool("o", false, "Enable overwrite mode")
	rewrite := flag.Bool("r", false, "Enable rewrite mode")
	flag.Parse()

	//Show help if no arguments are found or -h flag is used
	if *showHelp || len(flag.Args()) == 0 {
		printHelp()
		return
	}

	if len(flag.Args()) < 3 {
		fmt.Println("Error: Not enough arguments provided.")
		printHelp()
		return
	}

	//Store arguments for convenience
	inputPath := flag.Args()[0]
	outputPath := flag.Args()[1]
	lookupPath := flag.Args()[2]

	//Load the input and check if it exists
	userInput, err := loadFile(inputPath)
	if err != nil {
		log.Println("Error:", err)
		fmt.Println("Input not found")
		return
	}

	//Check for type of output
	outputType = "txt"
	if strings.HasSuffix(outputPath, ".html") {
		outputType = "html"
	}

	//Load the airport-lookup.csv
	lookup, err := loadFile(lookupPath)
	if err != nil {
		log.Println("Error", err)
		fmt.Println("Airport lookup not found")
		return
	}

	//Reading the airport-lookup.csv
	r := csv.NewReader(strings.NewReader(lookup))
	records, err := r.ReadAll()

	if err != nil {
		fmt.Printf("\n%sError reading CSV: %v%s\nAirport lookup malformed\n", Red, err, Reset)

		return //Interrupt if there's something wrong with the file
	}

	//Correct lookup has 6 columns
	const expectedColumns = 6

	//Check for correct amount of columns in the header
	if len(records[0]) != expectedColumns {
		fmt.Printf("\n%sInvalid amount of columns - data malformed in %v %s\n", Yellow, lookupPath, Reset)
		os.Exit(1) //Exit with an error
	}

	var invalidLookupRows int
	var nameid *int
	var countryid *int
	var municid *int
	var icaoid *int
	var iataid *int
	var coordid *int

	for i, record := range records {
		//Adjust for non-standard airport lookup column order
		if i == 0 {
			for i, head := range record {
				switch head {
				case "name":
					nameid = &i
				case "iso_country":
					countryid = &i
				case "municipality":
					municid = &i
				case "icao_code":
					icaoid = &i
				case "iata_code":
					iataid = &i
				case "coordinates":
					coordid = &i
				}
			}

			continue
		}

		if len(record) != expectedColumns {
			fmt.Printf("\n%sSkipping row %d: Expected %d columns, got %d%s\n", Yellow, i+1, expectedColumns, len(record), Reset)
			continue
		}

		//Check for empty fields or UTF-8 exceeding characters in the lookup
		malformedRow := false
		for i, data := range record {
			if strings.TrimSpace(data) == "" {
				malformedRow = true
				break
			}
			if i != 5 {
				if !validString(data) {
					malformedRow = true
				}
			}
		}

		//Skip invalid data and count them
		if malformedRow {
			invalidLookupRows++
			continue
		}

		// Append valid data
		airport := Airport{
			Name:         record[*nameid],
			ISO_Country:  record[*countryid],
			Municipality: record[*municid],
			ICAO_Code:    record[*icaoid],
			IATA_Code:    record[*iataid],
			Coordinates:  record[*coordid],
		}

		airports = append(airports, airport)

	}

	//Inform the user of skipped records
	if invalidLookupRows > 0 {
		fmt.Printf("\n%sCould not read %d airport records. Exceeded UTF-8 characters%s\n", Yellow, invalidLookupRows, Reset)
	}

	//Check if the output already exists
	exist := false
	if _, err := os.Stat(outputPath); err == nil {
		exist = true
	} else if os.IsNotExist(err) {
		exist = false
	} else {
		//Some other error occurred (e.g., permission issues)
		fmt.Printf("Could not access file %v, option to overwrite is unavailable - %v\n", outputPath, err)
	}

	//Prompt the user to choose to Overwrite / Keep / Cancel
	var choice int
	var newOutputPath string
	originalOutputPath := outputPath
	iterate := 1
	if exist && !*overwrite {
		for {
			if !*rewrite {
				fmt.Printf("\n%s%v already exists%s\n\n", Red, outputPath, Reset)
				//Print the prompt
				fmt.Println("Choose an option:")
				fmt.Println("1 - Overwrite")
				if outputType == "txt" {
					fmt.Printf("2 - Change Name to:%v (%d).txt\n", originalOutputPath[:len(originalOutputPath)-4], iterate)
				} else if outputType == "html" {
					fmt.Printf("2 - Change Name to:%v (%d).html\n", originalOutputPath[:len(originalOutputPath)-5], iterate)
				}
				fmt.Println("3 - Cancel")

				//Scan for choice
				_, err := fmt.Scanf("%d\n", &choice)

				//Check for choice validity
				if err != nil || (choice < 1 || choice > 3) {
					fmt.Println("Invalid choice. Please enter 1, 2, or 3.")
					continue
				}
			} else {
				choice = 2
			}

			//Change output if user chose to
			if choice == 2 {
				if outputType == "txt" {
					if iterate > 1 {
						if iterate > 9 {
							outputPath = outputPath[:len(outputPath)-9] + ".txt"
						} else {
							outputPath = outputPath[:len(outputPath)-8] + ".txt"
						}
					}
					newOutputPath = outputPath[:len(outputPath)-4] + " (" + strconv.Itoa(iterate) + ").txt"
				} else if outputType == "html" {
					if iterate > 1 {
						if iterate > 9 {
							outputPath = outputPath[:len(outputPath)-10] + ".html"
						} else {
							outputPath = outputPath[:len(outputPath)-9] + ".html"
						}
					}
					newOutputPath = outputPath[:len(outputPath)-5] + " (" + strconv.Itoa(iterate) + ").html"
				}
				// Check if the new output path still exists
				if _, err := os.Stat(newOutputPath); err == nil {
					if !*rewrite {
						fmt.Println("File with the new name also exists. Please choose again.")
					}
					outputPath = newOutputPath
					iterate++
					continue
				} else if !os.IsNotExist(err) {
					fmt.Printf("Could not access file %v, option to overwrite is unavailable - %v\n", outputPath, err)
					continue
				}
				outputPath = newOutputPath
			}

			//When valid choice, continue
			fmt.Println("You chose to:", options[choice])
			if choice == 1 || choice == 2 {
				break
			}
			if choice == 3 {
				return
			}
		}

	}

	if outputType == "html" {
		//Format the input for the output as an html
		userInput = getOutputStringHTML(userInput)
	} else {
		//Format the input for the output
		userInput = getOutputString(userInput)
	}

	//Create output
	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Error creating file: ", err)
		return
	}
	defer file.Close()

	//Write into output
	_, err = file.WriteString(userInput)
	if err != nil {
		fmt.Println("Error writing to file: ", err)
		return
	}

	println(outputPath, " created succesfully")

	//Testing tools
	if len(os.Args) > 4 {
		if os.Args[4] == "test" {
			switch os.Args[5] {
			case "string":

				// Testing func validString
				print(airports[24].Name + " ")
				println(validString(airports[24].Name))

			case "input":

				// Testing an input
				println(userInput)

			}
			return
		}
	}
}
