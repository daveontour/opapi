Example of a flight update message

~~~json
{
	"Action": "UPDATE",
	"FlightId": {
		"FlightKind": "Departure",
		"FlightNumber": "6111",
		"ScheduledDate": "2023-09-06",
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
		"ScheduledTime": "2023-09-06T11:40:00",
		"LinkedFlight": {
			"FlightId": {
				"FlightKind": "Arrival",
				"FlightNumber": "6110",
				"ScheduledDate": "2023-09-06",
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
				"ScheduledTime": "2023-09-06T10:31:00",
				"FlightUniqueID": "ARR_873758"
			}
		},
		"AircraftType": {
			"AircraftTypeId": {
				"AircraftTypeCode": {
					"IATA": "LOH",
					"ICAO": "C130"
				}
			},
			"Values": {
				"Name": "Lockheed Martin"
			}
		},
		"Aircraft": {
			"AircraftId": {
				"Registration": "NJM306"
			}
		},
		"Route": {
			"CustomType": "International",
			"ViaPoints": [
				{
					"SequenceNumber": "0",
					"AirportCode": {
						"IATA": "DEL",
						"ICAO": "VIDP"
					}
				}
			]
		},
		"Values": {
			"FlightUniqueID": "DEP_883547"
		},
		"StandSlots": [
			{
				"StartTime": "2023-09-06T10:31:00",
				"EndTime": "2023-09-06T12:26:00",
				"Name": "122",
				"ExternalName": "122",
				"AreaName": "Apron 1"
			}
		],
		"CarouselSlots": [],
		"GateSlots": [
			{
				"StartTime": "2023-09-06T11:26:00",
				"EndTime": "2023-09-06T12:36:00",
				"Category": "departure",
				"Name": "04",
				"ExternalName": "4",
				"AreaName": "T1"
			}
		],
		"CheckInSlots": [
			{
				"StartTime": "2023-09-06T08:40:00",
				"EndTime": "2023-09-06T10:55:00",
				"Category": "Business"
			},
			{
				"StartTime": "2023-09-06T08:40:00",
				"EndTime": "2023-09-06T10:40:00",
				"Category": "Economy"
			},
			{
				"StartTime": "2023-09-06T08:40:00",
				"EndTime": "2023-09-06T10:40:00",
				"Category": "Economy"
			},
			{
				"StartTime": "2023-09-06T08:40:00",
				"EndTime": "2023-09-06T10:40:00",
				"Category": "Economy"
			}
		],
		"ChuteSlots": []
	},
	"Changes": {
		"AircraftTypeChange": null,
		"AircraftChange": null,
		"CarouselSlotsChange": null,
		"GateSlotsChange": {
			"OldValue": [
				{
					"StartTime": "2023-09-06T11:21:00",
					"EndTime": "2023-09-06T12:31:00",
					"Category": "departure",
					"Name": "04",
					"ExternalName": "4",
					"AreaName": "T1"
				}
			],
			"NewValue": [
				{
					"StartTime": "2023-09-06T11:26:00",
					"EndTime": "2023-09-06T12:36:00",
					"Category": "departure",
					"Name": "04",
					"ExternalName": "4",
					"AreaName": "T1"
				}
			]
		},
		"StandSlotsChange": {
			"OldValue": [
				{
					"StartTime": "2023-09-06T10:31:00",
					"EndTime": "2023-09-06T12:21:00",
					"Name": "122",
					"ExternalName": "122",
					"AreaName": "Apron 1"
				}
			],
			"NewValue": [
				{
					"StartTime": "2023-09-06T10:31:00",
					"EndTime": "2023-09-06T12:26:00",
					"Name": "122",
					"ExternalName": "122",
					"AreaName": "Apron 1"
				}
			]
		},
		"ChuteSlotsChange": null,
		"CheckInSlotsChange": null,
		"RouteChange": null,
		"LinkedFlightChange": null
	},
	"ValueChanges": []
}
~~~