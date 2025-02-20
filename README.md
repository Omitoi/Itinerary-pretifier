# Itinerary Prettifier

## Overview

Itinerary Prettifier is a command-line tool that transforms raw flight itinerary text files into customer-friendly formats. This tool is designed to streamline the workflow of back-office administrators at "Anywhere Holidays" by automating the conversion of complex and technical itinerary data into a readable format for customers.

## Features

- Reads raw itinerary text files and processes them into customer-friendly versions.
- Converts IATA and ICAO airport codes into full airport names using an airport lookup CSV.
- Formats ISO 8601 dates and times into human-readable formats.
- Trims excessive vertical whitespace to improve readability.
- Provides error handling for incorrect input formats and missing files.
- Supports dynamic airport lookup column orders (Bonus feature).
- Converts airport codes to city names if prefixed with `*` (Bonus feature).
- Outputs an email, ready to send html file when output has suffix .html (Bonus feature).

## Installation

### Prerequisites
- Go installed on your system (https://go.dev/doc/install)

### Clone the Repository
```sh
$ git clone https://github.com/Omitoi/Itinerary-pretifier
$ cd itinerary
```

## Usage

### Running the tool
```sh
$ go run . [-o]/[-r] ./input.txt ./output.txt ./airport-lookup.csv
```
  - o - Automatically overwrites your output without prompting.
  - r - Automatically rewrites your output name without prompting

### Help Menu
Display the usage instructions:
```sh
$ go run . -h
```

### Input Format
- The itinerary file should contain raw text with embedded airport codes (`#IATA` or `##ICAO`) and ISO 8601 timestamps in the following formats:
  - `D(YYYY-MM-DDTHH:MM±HH:MM)` → Converted to `DD-Mmm-YYYY`
  - `T12(YYYY-MM-DDTHH:MM±HH:MM)` → Converted to `HH:MM AM/PM (Offset)`
  - `T24(YYYY-MM-DDTHH:MM±HH:MM)` → Converted to `HH:MM (Offset)`
- Excessive blank lines should be reduced to a maximum of one.

### Airport Lookup Format
- The CSV file must have the following columns: `name, iso_country, municipality, icao_code, iata_code, coordinates`.
- If any column is missing or blank, an error will be thrown.
- Example row:
  ```csv
  London Heathrow Airport, GB, London, EGLL, LHR, 51.4706,-0.4619
  ```

### Example Input
```txt
Departure: #LAX
Arrival: ##EGLL
Date: D(2023-06-15T14:00-07:00)
Time: T12(2023-06-15T14:00-07:00)
```

### Example Output
```txt
Departure: Los Angeles International Airport
Arrival: London Heathrow Airport
Date: 15-Jun-2023
Time: 2:00PM (-07:00)
```

## Error Handling
- If incorrect arguments are provided, the tool prints the usage instructions.
- If the input file does not exist, the program outputs: `Input not found`.
- If the airport lookup file is missing, it outputs: `Airport lookup not found`.
- If the airport lookup CSV is malformed, it outputs: `Airport lookup malformed`.
- If any date/time format is incorrect, it remains unchanged in the output.

## Bonus Features
- **City Name Conversion:** Converts airport codes to city names when prefixed with `*` (e.g., `*#LHR` → `London`).
- **Dynamic CSV Column Order Handling:** Allows for flexibility in the airport lookup CSV file structure.
- **Optional HTML output:** Ready to send output for emailing.

## Author
Omitoi | Petr Kubec
