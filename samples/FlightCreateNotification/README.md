

~~~json
{
	"Action": "CREATE",
	"FlightId": {
		"FlightKind": "Arrival",
		"FlightNumber": "123",
		"ScheduledDate": "2023-09-06",
		"AirportCode": {
			"IATA": "AUH",
			"ICAO": "OMAA"
		},
		"AirlineDesignator": {
			"IATA": "QF",
			"ICAO": "QFA"
		}
	},
	"FlightState": {
		"ScheduledTime": "2023-09-06T13:45:00",
		"LinkedFlight": {
			"FlightId": {
				"FlightKind": "Departure",
				"FlightNumber": "456",
				"ScheduledDate": "2023-09-06",
				"AirportCode": {
					"IATA": "AUH",
					"ICAO": "OMAA"
				},
				"AirlineDesignator": {
					"IATA": "QF",
					"ICAO": "QFA"
				}
			},
			"Values": {
				"ScheduledTime": "2023-09-06T15:45:00",
				"FlightUniqueID": "DEP_883620"
			}
		},
		"AircraftType": {
			"AircraftTypeId": {
				"AircraftTypeCode": {
					"IATA": "744",
					"ICAO": "B744"
				}
			},
			"Values": {
				"Name": "Boeing 747-400"
			}
		},
		"Aircraft": {
			"AircraftId": {
				"Registration": ""
			}
		},
		"Route": {
			"CustomType": "International",
			"ViaPoints": [
				{
					"SequenceNumber": "0",
					"AirportCode": {
						"IATA": "MEL",
						"ICAO": "YMML"
					}
				}
			]
		},
		"Values": {
			"FlightUniqueID": "ARR_873807"
		},
		"StandSlots": [
			{
				"StartTime": "2023-09-06T13:45:00",
				"EndTime": "2023-09-06T15:45:00"
			}
		],
		"CarouselSlots": [
			{
				"StartTime": "2023-09-06T13:45:00",
				"EndTime": "2023-09-06T14:25:00",
				"Category": ""
			}
		],
		"GateSlots": [
			{
				"StartTime": "2023-09-06T13:40:00",
				"EndTime": "2023-09-06T13:50:00",
				"Category": "arrival"
			}
		],
		"CheckInSlots": [],
		"ChuteSlots": []
	},
	"Changes": {
		"AircraftTypeChange": null,
		"AircraftChange": null,
		"CarouselSlotsChange": null,
		"GateSlotsChange": null,
		"StandSlotsChange": null,
		"ChuteSlotsChange": null,
		"CheckInSlotsChange": null,
		"RouteChange": null,
		"LinkedFlightChange": null
	},
	"ValueChanges": []
}
~~~