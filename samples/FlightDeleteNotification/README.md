

~~~json
{
	"Action": "DELETE",
	"FlightId": {
		"FlightKind": "Departure",
		"FlightNumber": "7697",
		"ScheduledDate": "2023-09-06",
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
		"ScheduledTime": "2023-09-06T01:01:00",
		"LinkedFlight": {
			"FlightId": {
				"FlightKind": "Arrival",
				"FlightNumber": "7696",
				"ScheduledDate": "2023-09-06",
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
				"ScheduledTime": "2023-09-06T00:05:00",
				"FlightUniqueID": "ARR_872372"
			}
		},
		"AircraftType": {
			"AircraftTypeId": {
				"AircraftTypeCode": {
					"IATA": "73H",
					"ICAO": "B738"
				}
			},
			"Values": {
				"Name": "Boeing 737-800 ( Winglet ) "
			}
		},
		"Aircraft": {
			"AircraftId": {
				"Registration": "SUGED"
			}
		},
		"Route": {
			"CustomType": "International",
			"ViaPoints": [
				{
					"SequenceNumber": "0",
					"AirportCode": {
						"IATA": "LHR",
						"ICAO": "EGLL"
					}
				}
			]
		},
		"Values": {
			"FlightUniqueID": "DEP_882169"
		},
		"StandSlots": [
			{
				"StartTime": "2023-09-06T00:05:00",
				"EndTime": "2023-09-06T01:01:00"
			}
		],
		"CarouselSlots": [],
		"GateSlots": [
			{
				"StartTime": "2023-09-06T00:11:00",
				"EndTime": "2023-09-06T01:11:00",
				"Category": "departure"
			}
		],
		"CheckInSlots": [
			{
				"StartTime": "2023-09-05T22:01:00",
				"EndTime": "2023-09-06T00:16:00",
				"Category": "Business"
			},
			{
				"StartTime": "2023-09-05T22:01:00",
				"EndTime": "2023-09-06T00:01:00",
				"Category": "Economy"
			},
			{
				"StartTime": "2023-09-05T22:01:00",
				"EndTime": "2023-09-06T00:01:00",
				"Category": "Economy"
			},
			{
				"StartTime": "2023-09-05T22:01:00",
				"EndTime": "2023-09-06T00:01:00",
				"Category": "Economy"
			}
		],
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