<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a name="readme-top"></a>



<!-- PROJECT LOGO -->
<br />
<div align="center">
 <!-- <a href="https://github.com/othneildrew/Best-README-Template">
    <img src="images/logo.png" alt="Logo" width="80" height="80">
  </a>
  -->

  <h3 align="center">Get Flights and Resources REST API for SITA AMS 6.6 and 6.7</h3>
  <h4 align="center">Including Webhooks and RabbitMQ distribution of updates and scheduled refreshes</h4>


  <p align="center">
    A Rest API service and notification service for SITA AMS
    <br />
    <a href="https://github.com/daveontour/opapi/issues">Report Bug</a>
    <a href="https://github.com/daveontour/opapi/issues">Request Feature</a>
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#installation-and-execution">Installation and Execution</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage-and-api-reference">Usage and API Reference</a></li>
    <ul>
        <li><a href="#installation-and-execution">Installation and Execution</a></li>
        <li><a href="#/getflights/[airport]?{options}&{options..}">/getFlights</a></li>
    </ul>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

  <img src="images/operation.png" alt="Logo">


SITA AMS has a set of APIs that is quite comprehensive, however:
* They can be difficult to use
* If used without discretion, they can impact system performance.
* The data is provided in a format not easily used with modern frameworks
* They have limited capability to refine the search for particular data
* Response time of the API can be affected by other operations taking place on the AMS Server
* All users receive the same data from the API, regardless of any sensitivities
* The granularity of notifaction messages is very coarse

This project addreses these issues by introducing a Restful API Server to service user requests. <br/>
The service loads and caches flights from AMS and then continually listens for update notification from AMS to maintain the cache. Periodically, the service "advances" the window operating times of flights in the cache by requesting data from AMS. At the same time the service prunes old flights from the cache by removing flights that have fallen out of the window of configured flight operating times<br/>

 The service exposes an API to retrieve flights and to retrieve the allocation of resources out of the cache without needing to go to AMS. Users can provide identity information in their API request to recieve a customised set of data or if no identity is provided, then they will receive a default set of publicly available data<br/>

The service also provides Push notification of changes and regular scheduled updates to subscribers via a WebHooks mechanism or via a RabbitMQ Exchange


<p align="right">(<a href="#readme-top">back to top</a>)</p>


## Design Goals
 
 - Isolate AMS from each individual API call
 - Return results in less than 1 second
 - Cater for multiple airports or instance of AMS using a single instance of the API
 - Discriminate users accessing the API and provide only the data they are entitled to see
 - Provide enough richness in the API to cover a variety of operational Use Cases
 - Minimise memory usage and growth of heap space
 - Do not introduce any new infrastructure components not already present in a typical AMS setup
 - Limit administration tasks to configuration only.
 - Include the tools to test performance and demonstrate capability
 - Provide a change notification service where users can define the types of changes they are interested in
 - Provide a no configuration change notification service for annonymous users

In addition to the main service (opapi.exe) project includes four utility programs:

- **opapiseeder.exe** This utility is used to seed the serive with demonstration flights and allocation when the servoce is run in Demonstration or Performance Test Mode
- **webhookclient.exe** The service has the ability to send change notifications and status updates via a WebHooks implementation. This program allows you to run a demonstration webhook client to deomonstate this capability
- **rabbitclient.exe** Same functionality as the webhookclient except for RabbitMQ clients
- **perftestclient.exe** This utility allows the response time performance of the API to be tested over a period of time by making repeated call to the API and recording response times

### Built With

The service is built using the Go ( version 1.20 ) programming language. No other runtime software components are require other than the service itself

* [Go Programming Language][go-url]


<p align="right">(<a href="#readme-top">back to top</a>)</p>

## DISCLAIMER
This is not official software from SITA. It was built as a learning exercise to learn the Go programming language, using a real world problem that I am familiar with. <br/>
No support or warranty should be inferred. There is no gaurantee for fitness of purpose<br/>
No endorsement of this software be SITA should be inferred <br/>
The information on using the SITA AMS API was obtained from publicly available information<br/>

<!-- GETTING STARTED -->
## Getting Started

To run the service, either download the latest release, or clone this repository and build it yourself in and IDE that supports Go 1.20


### Installation and Execution

The API service can be installed to run as a windows service or from the Windows command line prompt <br/>
The service does *not* require any additional software componets. The service _self hosts_ itself according the the configuration in _*service.json*_ file

1. To install and start as a Windows Service. Must be logged on with Administrator privelege
   ```cmd
   C:\ProgramFiles\opapi\opapi.exe install
   C:\ProgramFiles\opapi\opapi.exe start
   ```
2. To stop the windows service and uninstall. Must be logged on with Administrator privelege
   ```cmd
   C:\ProgramFiles\opapi\opapi.exe stop
   C:\ProgramFiles\opapi\opapi.exe uninstall
   ```
3. To run from command line as any user
   ```cmd
   C:\ProgramFiles\opapi\opapi.exe run
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>


## Sample Output

Responses to the Rest API request, change notifications and subscription notifications are provided in JSON format.
API requests will include meta data in the message to describe the response. 
Samples of the messages can be found at the links below 

- [Get Flights Response Sample][getflights-sample-url]
- [Flight Create Notification Sample][getflights-sample-url]
- [Flight Create Notification Sample][createflight-sample-url]
- [Flight Deleted Notification Sample][deleteflight-sample-url]
- [Get Allocations Sample][getallocations-sample-url]
- [Get Configured Resources Sample][getconfiguredresources-sample-url]

## Usage and API Reference

 This service exposes three API endpoints to retreive data on flights, resource allocations and configured resources:

 * /getFlights
 * /getAllocations
 * /getConfiguredResources

The APIs are accessed via HTTP GET Requests and return data in JSON format

#### Request Header
The HTTP Get request header should include a parameter called "Token". <br />
The value of "Token" is assigned by the system administrator to identify your user profile which defines your rights and capabilities to acces the APIs

If the token header is not present, you will be assigned the rights of the "default user", if one is configured by the administrator


## /getFlights/[Airport]?{options}&{options..}
Retreive flight details

|Parameter|Description|Example|
|----|----------|-----|
|**Airport**|Three letter IATA airport code to the desired airport|/getFlights/APT|
|**al or airline**|Two letter IATA code for the airline, eg. BA, DL, LH, MH (default: all airlines)|/getFlights/APT?al=QF|
|**flt or flight**|Flight Number, eg. QF001, EK23, BA007 (default: all flights|/getFlights/APT?flt=QF001|
|**d or direction**|The direction of flight relative to the home airport. either 'Arr' or 'Dep'|/getFlights/APT?d=Arr|
|**r or route**|The route of the flight|/getFlights/APT?r=MEL|
|**from**|Number of hours relative to 'now' for the earliest scheduled time of operation for the flight, eg. -3 (default: -12|/getFlights/APT?from=-12|
|**to**|Number of hours relative to 'now' for the latest scheduled time of operation for the flight, eg. 12 (default: 24)|/getFlights/APT?to=48|
|**updatedSince**|Return records that have been updated from the date, e.g. 2023-07-16T13:00:00|/getFlights/APT?upatedSince=2023-07-16T13:00:0|
|**{custom field name}**|Return records have the specified custom field name equal to the specified value. Users will only be able to search on the custom fields they have been granted access to. (see below)|/getFlights/APT?Sh--_GroundHandler=EAS|
|**max**|Limit the total number of flights returned to the number specified|/getFlights/APT?max=5|


### Examples 

  Find the flights from now until 12 hours from now<br />
  **/getFlights/APT?from=0&amp;to=12**<br />

  Find the Qantas flights from now until 12 hours from now<br />
  **/getFlights/APT?al=QF&amp;from=0&amp;to=12**<br />

  Find the flights arriving from Melbourne<br />
  **/getFlights/APT?route=MEL&amp;d=Arr**<br />

  Find all the flight where the custom field name **Sh--_GroundHandler** of the assigned flight is  EAS<br />
  **/getFlights/APT?Sh--_GroundHandler=EAS**<br />

[Get Flights Response Sample][getflights-sample-url]

## /getAllocations/[Airport]?{options}
Retreive flights allocated to resources

|Parameter|Description|Example|
|----|----------|-----|
|**Airport**|Three letter IATA airport code to the desired airport|/getAllocations/APT|
|**flt or flight**|Flight Number, eg. QF001, EK23, BA007 (default: all flights)|/getAllocations/APT?flt=QF001|
|**al or airline**|Two letter IATA code for the airline, eg. BA, DL, LH, MH (default: all airlines|/getAllocations/APT?flt=QF|
|**rt or resourceType**|One of CheckIn, Gate, Stand, Carousel, Chute. (default: all types are returned)|getAllocations/APT?rt=Gate|
|**id or resource**|The name of the individual reource to query. Query must include the resourceType parameter (default: all resources)|/getAllocations/APT?rt=Gate&amp;id=100|
|**from**|Number of hours relative to 'now' to start looking for allocations, eg. -3 (default:-12)|getAllocations/APT?from=-12|
|**to**|Number of hours relative to 'now' to stop looking for allocations, eg. 12 (default: 24)|getResources/APT?to=7|
|**sort**|Either "resource" or "time" to specify the sort order of the allocations returned (default: resource)|/getAllocations/APT?sort=time|
|**updatedSince**|Return records that have been updated from the date, e.g. 2023-07-16T13:00:00|/getAllocations/APT?upatedSince=2023-07-16T13:00:00|



### Examples

  Find the flights allocated to checkin desk 100 from now until 12 hours from now<br />
    **/getResources/APT?from=0&amp;to=12&amp;rt=CheckIn&amp;id=100**<br />
    <br />
    Find all the resources allocated to flight QF100<br />
    **/getResources/APT?flt=QF100**<br />
    <br />
    Find all the resources allocated to Emirates (EK)<br />
    **/getResources/APT?al=EK**<br />
    <br />
    Find all the resources allocated to British Airways (BA) for the next 3 days<br />
    **/getResources/APT?al=BA&amp;from=0&amp;to=72**<br />
    <br />
    Find all the resources where the custom field name **Sh--_GroundHandler** of the assigned flight is
    EAS<br />
    **/getResources/APT?Sh--_GroundHandler=EAS**
  
## /getConfiguredResources/[Airport]/{ResourceType}
Retreive the configured resources for the airport

|Parameter|Description|Example|
|----|----------|-----|
|**Airport**|Three letter IATA airport code to the desired airport|/getConfiguredResources/APT|
|**{Resource Type}**|One of CheckIn, Gate, Stand, Carousel, Chute. (default: all types are returned)|/getConfiguredResources/APT/Gate|


## Sample Outputs




# Configuring the Service

The execution of the service is controlled by the configuration in the file **service.json** in the directory the application is installed in.
An example of the contents of the service.json file is shown below

~~~json
{
    "ServiceName": "GetFlightAndResourceService",
    "ServiceDisplayName": "GetFlight and Resource Rest API",
    "ServiceDescription": "A  HTTP/JSON  Rest Service for retrieving flights and resource allocations from AMS",
    "ServiceIPport": "127.0.0.1:8081",
    "ScheduleUpdateJob": "02:00:00",
    "ScheduleUpdateJobIntervalInHours": 1,
    "ScheduleUpdateJobIntervalInMinutes": -1,
    "DebugService": true,
    "TraceService": false,
    "UseHTTPS": false,
    "UseHTTPSUntrusted": false,
    "KeyFile": "keyFile",
    "CertFile": "certFile",
    "TestHTTPServer": false,
    "LogFile": "c:/Users/dave_/Desktop/Logs/process.log",
    "RequestLogFile": "c:/Users/dave_/Desktop/Logs/request.log",
    "MaxLogFileSizeInMB": 50,
    "MaxNumberLogFiles": 3,
    "EnableMetrics": true,
    "MetricsLogFile": "c:/Users/dave_/Desktop/Logs/performance.log",
    "AdminToken": "davewashere",
    "NumberOfChangePushWorkers":7,
    "NumberOfSchedulePushWorkers":5
}
~~~

|Parameter|Description|Notes|
|----|----------|----|
|ServiceName|GetFlightAndResourceService||
|ServiceDisplayName|The service display name e.g."GetFlight and Resource Rest API"||
|ServiceDescription|The service description e.g. "A HTTP/JSON  Rest Service for retrieving flights and resource allocations from AMS"||
|ServiceIPport|The IP address and port number for the service to bind to, e.g. 127.0.0.1:8081||
|ScheduleUpdateJob|02:00:00||
|ScheduleUpdateJobIntervalInHours|The interval between running update jobs to update the repository from AMS in hours||
|ScheduleUpdateJobIntervalInMinutes|The interval between running update jobs to update the repository from AMS in minutes|This is a seperate schedule to the one specified by the "hours" schedule |
|DebugService|||
|TraceService|||
|UseHTTPS|Set to "true" to run the service using HTTPS. If true, then the KeyFile and CertFile paramters must be set |"false" by default|
|UseHTTPSUntrusted|Set to true to run the service using HTTPS with a locally generated certificate||
|KeyFile|keyFile||
|CertFile|certFile||
|TestHTTPServer|||
|LogFile|c:/Users/dave_/Desktop/Logs/process.log||
|RequestLogFile|"c:/Users/dave_/Desktop/Logs/request.log||
|MaxLogFileSizeInMB|The maximum size of the application log file||
|MaxNumberLogFiles|When the log file have reached their maximum size they will be archived. This parameter specifies how many archive logs to keep||
|EnableMetrics|true||
|MetricsLogFile|c:/Users/dave_/Desktop/Logs/performance.log||
|AdminToken|The header token to identify the user with administrator capability ||
|NumberOfChangePushWorkers|The number of worker Go processes for managing change notifications||
|NumberOfSchedulePushWorkers|The number of worker Go processes for managing distribution of subscription updates|Maximum value should not exceed the total number of subscriptions|


# Configuring Airports


Each airport to be serviced by the API service needs to be configured in the file **airports.json** in the directory the application is installed in. More than one airport can be configured. 

An example of the contents of the service.json file is shown below
~~~json
{
    "Repositories": [
        {
            "AMSAirport": "ABC",
            "AMSToken": "0ab7d73d-e93a-480b-ba8c-a35943161cb0",
            "AMSSOAPServiceURL": "http://localhost/SITAAMSIntegrationService/v2/SITAAMSIntegrationService",
            "AMSRestServiceURL": "http://localhost/api/v1/",
            "FlightSDOWindowMinimumInDaysFromNow": -3,
            "FlightSDOWindowMaximumInDaysFromNow": 20,
            "ListenerType":"MSMQ",
            "NotificationListenerQueue": ".\\private$\\tow_tracker",
            "LoadFlightChunkSizeInDays": 4,
            "RabbitMQConnectionString": "amqp://amsauh:amsauh@localhost:5672/amsauh",
            "RabbitMQExchange": "Test",
            "RabbitMQTopic": "AMSX.Notify",
            "PublishChangesRabbitMQConnectionString": "amqp://amsauh:amsauh@localhost:5672/amsauh",
            "PublishChangesRabbitMQExchange": "Test",
            "PublishChangesRabbitMQTopic": "AMSJSON.Notify"
        },
        {
            "AMSAirport": "DEF",
            "AMSToken": "0ab7d73d-e93a-480b-ba8c-a35943161cb0",
            "AMSSOAPServiceURL": "http://localhost/SITAAMSIntegrationService/v2/SITAAMSIntegrationService",
            "AMSRestServiceURL": "http://localhost/api/v1/",
            "FlightSDOWindowMinimumInDaysFromNow": -3,
            "FlightSDOWindowMaximumInDaysFromNow": 20,
            "ListenerType":"MSMQ",
            "NotificationListenerQueue": ".\\private$\\tow_tracker",
            "LoadFlightChunkSizeInDays": 4,
            "RabbitMQConnectionString": "amqp://amsauh:amsauh@localhost:5672/amsauh",
            "RabbitMQExchange": "Test",
            "RabbitMQTopic": "AMSX.Notify",
            "PublishChangesRabbitMQConnectionString": "amqp://amsauh:amsauh@localhost:5672/amsauh",
            "PublishChangesRabbitMQExchange": "Test",
            "PublishChangesRabbitMQTopic": "AMSJSON.Notify"
        }
    ]
}
~~~

|Parameter|Description|Notes|
|----|----------|----|
|AMSAirport| "DEF"||
|AMSToken| "0ab7d73d-e93a-480b-ba8c-a35943161cb0"||
|AMSSOAPServiceURL| "http://localhost/SITAAMSIntegrationService/v2/SITAAMSIntegrationService"||
|AMSRestServiceURL| "http://localhost/api/v1/"||
|FlightSDOWindowMinimumInDaysFromNow| -3||
|FlightSDOWindowMaximumInDaysFromNow| 20||
|ListenerType|"MSMQ"||
|NotificationListenerQueue| ".\\private$\\tow_tracker"||
|LoadFlightChunkSizeInDays| 4||
|RabbitMQConnectionString| "amqp://amsauh:amsauh@localhost:5672/amsauh"||
|RabbitMQExchange| "Test"||
|RabbitMQTopic| "AMSX.Notify"||
|PublishChangesRabbitMQConnectionString| "amqp://amsauh:amsauh@localhost:5672/amsauh"||
|PublishChangesRabbitMQExchange| "Test"||
|PublishChangesRabbitMQTopic| "AMSJSON.Notify"||

# Configuring Users

Access to the service can be controlled by tokens that are passed in the header of the request. The configuration of the user controls the airports, airline and custom fields that the user has access to. The user configuration can also include subscription to update notification and change notifications. COnfiguration of users is in the **users.json** file

If the user does not provide a token for a valid configured user, then they are granted the access rights of the "default" user if configured. 

Below is an example of the **users.json** file with two users configured. 

- The "Default User" has restricted access just the ABC and APT airports. They are also restricted to flight information on WY and QF flights. The AllowedCustomFields array defines the custom field data that will be returned to them if available; not other custom fields will be returned to the user. 
- The second defined user, "Super User" has access to all airports, airlines and custome fields. The "*" in the configuration gives access to ALL.

Between these extremes, other users can be configured with appropriate access rights configured for each user


~~~json
{
    "Users": [
        {
            "Enabled": true,
            "UserName": "Default User",
            "Key": "default",
            "AllowedAirports": [
                "ABC", "APT"
            ],
            "AllowedAirlines": [
                "WY",
                "QF",
            ],
            "AllowedCustomFields": [
                "FlightUniqueID",
                "SYS_ETA",
                "FlightUniqueID",
                "Il--_DisembarkingFirstPad",
                "Sh--_GroundHandler",
                "Il--_TotalMalePax",
                "S---_Terminal",
            ],
            "UserChangeSubscriptions": [],
            "UserPushSubscriptions": []
        },
        {
            "Enabled": true,
            "UserName": "Super User",
            "Key": "reallyhardtoguessstring!@#*&^",
            "AllowedAirports": [
                "*"
            ],
            "AllowedAirlines": [
                "*"
            ],
            "AllowedCustomFields": [
                "*",
            ],
            "UserChangeSubscriptions": [],
            "UserPushSubscriptions": []
        },
    ]
}
~~~
## Configuring User Push Subscriptions

A user can have zero, one or more push subscriptions configured. A push subscription will send the current state for a set of flights or resources configured in the subscription. The push will occur at regular intervals as defined in the subscription. The data is sent to either a WebHook endpoint and/or published to a RabbitMQ exchange as configured in the subscription. 

## Configuring User Change Subscriptions

A user can have zero, one or more change notification subscriptions confiured. The change subscription defines the types of changes that the user in iterested in, e.g an aircraft type change, a flight delete, etc. When one of the interested changes occurs, the service will push the notification to the defined end point for the user. The defined end point can be a WebHooks client or a RabbitMQ Exchange 

# Running in Demonstrration Mode

The service can be run in Demonstration Mode to allow usage and demonstration of the API capability **without** a connection to AMS. When run in demonstration mode, the service uses the configuration in the file testfiles\test.json to define the airport and the resources available at the airport. Flights and Allocation can be loaded into the system using the **opapiseeder.exe** program (see below) 

   ```cmd
   C:\ProgramFiles\opapi\opapi.exe demo
   ```

## Seeding Demonstration Mode with Data

When the service is running in Demonstration Mode, it can be seeded with flight and resource allocation data by using the **opapiseeder.exe** program. 
The program must be run on the same server as the main opapi service. An example is shown below 

   ```cmd
   C:\ProgramFiles\opapi\opapiseeder.exe demo 5000 500 true
   ```

In this example, the service is seeded with 5000 flights, with each flight having 500 custom field entries. The final "true" parameter indicates that the seeder should continue to run an produce random updates every 20 seconds. Setting this final parameter to "false" prevents these random updates. 

**opapiseeder.exe** uses the configuration file testfiles\test.json to hold the configured resources and connection details for connecting to the main opapi service

# Running in Performance Test Mode
Performance Test Mode is very similar to Demonstration Mode, the exception being is that the flights are loaded and updated by the **opapiseeder.exe** program via a RabbitMQ interface, rather than a direct interface in Demonstration Mode. This allows closer approximation to a real life implementation than Demonstration Mode

   ```cmd
   C:\ProgramFiles\opapi\opapi.exe perfTest
   ```

## Seeding Performance Test Mode with Data
   ```cmd
   C:\ProgramFiles\opapi\opapiseeder.exe perfTest 5000 500 true
   ```

## Running with the Example Webhook Client

The example WebHooks client program can act help demonstrate the functionn of the service by acting as a demonstration client.
The client will receive Webhook updates from the service and log data on the received content of the message


  - To log all received messages to files.
   ```cmd
   C:\ProgramFiles\opapi\webhookclient.exe run localhost:8082 log
   ``` 

  - To print out the first 500 characters of the message on the console
   ```cmd
   C:\ProgramFiles\opapi\webhookclient.exe run localhost:8082 debug
   ``` 

  - To print out the entire contents of the message on the console
   ``` cmd
   C:\ProgramFiles\opapi\webhookclient.exe run localhost:8082 trace
   ``` 

  - To print out a message received message on the console
   ```cmd
   C:\ProgramFiles\opapi\webhookclient.exe run localhost:8082
   ``` 

## Performance Characteristics

### API Response Times

### Memory Useage

### GoLint 




<!-- ROADMAP -->
## Roadmap

- [ ] Add Changelog
- [ ] Add UI to manage configuration
- [ ] Add a realtime management and monitoring user interface
- [ ] Verify compatibility with UTF-16
- [ ] Add change notifications for resource allocation changes


See the [open issues](https://github.com/daveontour/opapi/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Dave Burton: daveontour57@gmail.com

Project Link: [https://github.com/daveontour/opapi](https://github.com/daveontour/opapi)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[getflights-sample-url]: https://github.com/daveontour/opapi/tree/main/samples/GetFlights "Get Flight Response Sample"
[createflight-sample-url]: https://github.com/daveontour/opapi/tree/main/samples/FlightCreateNotification "Flight Create Notification Sample"
[updateflight-sample-url]: https://github.com/daveontour/opapi/tree/main/samples/FlightUpdateNotification "Flight Update Notification Sample"
[deleteflight-sample-url]: https://github.com/daveontour/opapi/tree/main/samples/FlightDeleteNotification "Flight Deleted Notification Sample"
[getallocations-sample-url]: https://github.com/daveontour/opapi/tree/main/samples/GetAllocations "Get Allocations Sample"
[getconfiguredresources-sample-url]: https://github.com/daveontour/opapi/tree/main/samples/GetConfiguredResources "Get Configured Resources Sample"
[go-url]: https://go.dev/ "Go Programming Language"
