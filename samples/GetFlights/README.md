### Sample from a /getFlights api request

~~~json
{
	"Airport": "AUH",
	"Direction": "ARR/DEP",
	"ScheduleFlightsFrom": "2023-09-05T09:22:11",
	"ScheduleFlightsTo": "2023-09-06T21:22:11",
	"NumberOfFlights": "1449",
	"Airline": "*",
	"Flight": "*",
	"Route": "*",
	"CustomFieldQuery": [],
	"Warnings": [],
	"Errors": [],
	"Flights": [
		{
			"Action": "STATUS",
			"FlightId": {
				"FlightKind": "Arrival",
				"FlightNumber": "7072",
				"ScheduledDate": "2023-09-05",
				"AirportCode": {
					"IATA": "AUH",
					"ICAO": "OMAA"
				},
				"AirlineDesignator": {
					"IATA": "EK",
					"ICAO": "UAE"
				}
			},
			"FlightState": {
				"ScheduledTime": "2023-09-05T09:24:00",
				"LinkedFlight": {
					"FlightId": {
						"FlightKind": "Departure",
						"FlightNumber": "7073",
						"ScheduledDate": "2023-09-05",
						"AirportCode": {
							"IATA": "AUH",
							"ICAO": "OMAA"
						},
						"AirlineDesignator": {
							"IATA": "EK",
							"ICAO": "UAE"
						}
					},
					"Values": {
						"FlightUniqueID": "DEP_882463",
						"ScheduledTime": "2023-09-05T10:23:00"
					}
				},
				"AircraftType": {
					"AircraftTypeId": {
						"AircraftTypeCode": {
							"IATA": "32A",
							"ICAO": "A320"
						}
					},
					"Values": {
						"Name": "Airbus A320 ( Sharklet )"
					}
				},
				"Aircraft": {
					"AircraftId": {
						"Registration": "HZAS34"
					}
				},
				"Route": {
					"CustomType": "International",
					"ViaPoints": [
						{
							"SequenceNumber": "0",
							"AirportCode": {
								"IATA": "AMS",
								"ICAO": "EHAM"
							}
						}
					]
				},
				"Values": {
					"FlightUniqueID": "ARR_872666",
					"SYS_ETA": "2023-09-08T09:24:00"
				},
				"StandSlots": [
					{
						"StartTime": "2023-09-05T09:24:00",
						"EndTime": "2023-09-05T10:23:00"
					}
				],
				"CarouselSlots": [],
				"GateSlots": [],
				"CheckInSlots": [],
				"ChuteSlots": []
			}
		},
		{
			"Action": "STATUS",
			"FlightId": {
				"FlightKind": "Arrival",
				"FlightNumber": "7102",
				"ScheduledDate": "2023-09-05",
				"AirportCode": {
					"IATA": "AUH",
					"ICAO": "OMAA"
				},
				"AirlineDesignator": {
					"IATA": "AF",
					"ICAO": "AFR"
				}
			},
			"FlightState": {
				"ScheduledTime": "2023-09-05T09:25:00",
				"LinkedFlight": {
					"FlightId": {
						"FlightKind": "Departure",
						"FlightNumber": "7103",
						"ScheduledDate": "2023-09-05",
						"AirportCode": {
							"IATA": "AUH",
							"ICAO": "OMAA"
						},
						"AirlineDesignator": {
							"IATA": "AF",
							"ICAO": "AFR"
						}
					},
					"Values": {
						"FlightUniqueID": "DEP_882451",
						"ScheduledTime": "2023-09-05T10:56:00"
					}
				},
				"AircraftType": {
					"AircraftTypeId": {
						"AircraftTypeCode": {
							"IATA": "AWH",
							"ICAO": "A139"
						}
					},
					"Values": {
						"Name": "Agusta (Bell) (AB-139/AW-139)"
					}
				},
				"Aircraft": {
					"AircraftId": {
						"Registration": "LWA010F"
					}
				},
				"Route": {
					"CustomType": "International",
					"ViaPoints": [
						{
							"SequenceNumber": "0",
							"AirportCode": {
								"IATA": "SYD",
								"ICAO": "YSSY"
							}
						}
					]
				},
				"Values": {
					"FlightUniqueID": "ARR_872665",
					"SYS_ETA": "2023-09-08T09:25:00"
				},
				"StandSlots": [
					{
						"StartTime": "2023-09-05T09:25:00",
						"EndTime": "2023-09-05T10:56:00"
					}
				],
				"CarouselSlots": [],
				"GateSlots": [],
				"CheckInSlots": [],
				"ChuteSlots": []
			}
		},
		{
			"Action": "STATUS",
			"FlightId": {
				"FlightKind": "Arrival",
				"FlightNumber": "7080",
				"ScheduledDate": "2023-09-05",
				"AirportCode": {
					"IATA": "AUH",
					"ICAO": "OMAA"
				},
				"AirlineDesignator": {
					"IATA": "BA",
					"ICAO": "BAW"
				}
			},
			"FlightState": {
				"ScheduledTime": "2023-09-05T09:25:00",
				"LinkedFlight": {
					"FlightId": {
						"FlightKind": "Departure",
						"FlightNumber": "7081",
						"ScheduledDate": "2023-09-05",
						"AirportCode": {
							"IATA": "AUH",
							"ICAO": "OMAA"
						},
						"AirlineDesignator": {
							"IATA": "BA",
							"ICAO": "BAW"
						}
					},
					"Values": {
						"FlightUniqueID": "DEP_882464",
						"ScheduledTime": "2023-09-05T10:23:00"
					}
				},
				"AircraftType": {
					"AircraftTypeId": {
						"AircraftTypeCode": {
							"IATA": "32A",
							"ICAO": "A320"
						}
					},
					"Values": {
						"Name": "Airbus A320 ( Sharklet )"
					}
				},
				"Aircraft": {
					"AircraftId": {
						"Registration": "HZAS33"
					}
				},
				"Route": {
					"CustomType": "International",
					"ViaPoints": [
						{
							"SequenceNumber": "0",
							"AirportCode": {
								"IATA": "BAH",
								"ICAO": "OBBI"
							}
						}
					]
				},
				"Values": {
					"FlightUniqueID": "ARR_872664",
					"SYS_ETA": "2023-09-08T09:25:00"
				},
				"StandSlots": [
					{
						"StartTime": "2023-09-05T09:25:00",
						"EndTime": "2023-09-05T10:23:00"
					}
				],
				"CarouselSlots": [],
				"GateSlots": [],
				"CheckInSlots": [],
				"ChuteSlots": []
			}
		}
	]
}
~~~