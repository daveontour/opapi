<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a name="readme-top"></a>
<!--
*** Thanks for checking out the Best-README-Template. If you have a suggestion
*** that would make this better, please fork the repo and create a pull request
*** or simply open an issue with the tag "enhancement".
*** Don't forget to give the project a star!
*** Thanks again! Now go create something AMAZING! :D
-->



<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]



<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/othneildrew/Best-README-Template">
    <img src="images/logo.png" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">Get Flights and Resources REST API for SITA AMS 6.6 and 6.7</h3>
  <h4 align="center">Including Webhooks and RabbitMQ distribution of updates and scheduled refreshes</h4>


  <p align="center">
    A Rest API service and notification service for SITA AMS
    <br />
    <a href="https://github.com/daveontour/opapi">**Explore the docs »**</a>
    <br />
    <br />
    <a href="https://github.com/daveontour/opapi">View Demo</a>
    ·
    <a href="https://github.com/daveontour/opapi/issues">Report Bug</a>
    ·
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
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
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
* All users recieve the same data from the API, regardless of any sensitivities

This project addreses these issues by introducing a Restful API Server to service user requests. <br/>
The service loads and caches flights from AMS and then continually listed for updates to update the cache. The service exposes an API to retrieve flights and to retrieve the allocation of resources out of the cache without needing to go to AMS. Users can provide identity information in their API request to recieve a customised set of data or if no identity is provided, then they will receive a default set of publicly available data<br/>
<br/>
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


### Built With

The service is built using the Go ( version 1.20 ) programming language. No other runtime software components are require other than the service itself

* [![Go][Go]][Go-url]


<p align="right">(<a href="#readme-top">back to top</a>)</p>

## DISCLAIMER
This is not official software from SITA. It was built as a learning exercise to learn the Go programming language, using a real world problem that I am familiar with. <br/>
No support or warranty should be inferred. There is no gaurantee for fitness of purpose<br/>
No endorsement of this software bt SITA should be inferred <br/>
The information on using the SITA AMS API was obtained from publicly available information<br/>

<!-- GETTING STARTED -->
## Getting Started

To run the service, either download the latest release, or clone this repository and build it yourself in and IDE that supports Go 1.20


### Installation and Execution

The API service can be installed to run as a windows service or from the Windows command line prompt <br/>
The service does *not* require any additional software componets. The service _self hosts_ itself according the the configuration in _*service.json*_ 

1. To install and start as a Windows Service. Must be logged on with Administrator privelege
   ```cmd
   C:\ProgramFiles\opapi\opapi.exe uninstall
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



<!-- USAGE EXAMPLES -->
## Usage

Use this space to show useful examples of how a project can be used. Additional screenshots, code examples and demos work well in this space. You may also link to more resources.

_For more examples, please refer to the [Documentation](https://example.com)_

<p align="right">(<a href="#readme-top">back to top</a>)</p>



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

Use this space to list resources you find helpful and would like to give credit to. I've included a few of my favorites to kick things off!

* [Choose an Open Source License](https://choosealicense.com)
* [GitHub Emoji Cheat Sheet](https://www.webpagefx.com/tools/emoji-cheat-sheet)
* [Malven's Flexbox Cheatsheet](https://flexbox.malven.co/)
* [Malven's Grid Cheatsheet](https://grid.malven.co/)
* [Img Shields](https://shields.io)
* [GitHub Pages](https://pages.github.com)
* [Font Awesome](https://fontawesome.com)
* [React Icons](https://react-icons.github.io/react-icons/search)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Useage and API Reference

 This service exposes three API endpoints to retreive data on flights, resource allocations and configured resources:

 * /getFlights
 * /getAllocations
 * /getConfiguredResources

The APIs are accessed via HTTP GET Requests and return data in JSON format

 
#### Request Header
The HTTP Get request header should include a parameter called "Token". <br />
The value of "Token" is assigned by the system administrator to identify your user profile which defines your rights and capabilities to acces the APIs

If the token header is not present, you will be assigned the rights of the "default user", if one is configured by the administrator


## /getFlights/[Airport]?{options}
Retreive flight details<

|Parameter|Description|Example|
|----|----------|-----|
|Airport|Three letter IATA airport code to the desired airport|/getFlights/APT|
|al or airline|Two letter IATA code for the airline, eg. BA, DL, LH, MH (default: all airlines)|/getFlights/APT?al=QF|
|flt or flight|Flight Number, eg. QF001, EK23, BA007 (default: all flights|/getFlights/APT?flt=QF001|
|d or direction|The direction of flight relative to the home airport. either 'Arr' or 'Dep'|/getFlights/APT?d=Arr|
|r or route|The route of the flight|/getFlights/APT?r=MEL|
|from|Number of hours relative to 'now' for the earliest scheduled time of operation for the flight, eg. -3 (default: -12|/getFlights/APT?from=-12|
|to|Number of hours relative to 'now' for the latest scheduled time of operation for the flight, eg. 12 (default: 24)|/getFlights/APT?to=48|
|updatedSince|Return records that have been updated from the date, e.g. 2023-07-16T13:00:00|/getFlights/APT?upatedSince=2023-07-16T13:00:0|
|{custom field name}|Return records have the specified custom field name equal to the specified value|/getFlights/APT?Sh--_GroundHandler=EAS|


### Examples 

  Find the flights from now until 12 hours from now<br />
  **/getFlights/APT?from=0&amp;to=12**<br />

  Find the Qantas flights from now until 12 hours from now<br />
  **/getFlights/APT?al=QF&amp;from=0&amp;to=12**<br />

  Find the flights arriving from Melbourne<br />
  **/getFlights/APT?route=MEL&amp;d=Arr**<br />

  Find all the flight where the custom field name **Sh--_GroundHandler** of the assigned flight is  EAS<br />
  **/getFlights/APT?Sh--_GroundHandler=EAS**<br />


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


# Configuring the Service

The execution of the service is controlled by the configuration in the file **service.json** in the directory the application is installed in.
An example of the conents of the service.json file is shown below

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

# Configuring Airports


The execution of the service is controlled by the configuration in the file **airports.json** in the directory the application is installed in.
An example of the conents of the service.json file is shown below
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

# Configuring Users

## Configuring User Push Subsctiptions

## Configuring User Change Subscriptions

# Running in Demonstrration Mode

## Seeding Demonstration Mode with Data

# Running in Performance Test Mode

## Seeding Performance Test Mode with Data

## Running with the Example Webhook Client


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/othneildrew/Best-README-Template.svg?style=for-the-badge
[contributors-url]: https://github.com/othneildrew/Best-README-Template/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/othneildrew/Best-README-Template.svg?style=for-the-badge
[forks-url]: https://github.com/othneildrew/Best-README-Template/network/members
[stars-shield]: https://img.shields.io/github/stars/othneildrew/Best-README-Template.svg?style=for-the-badge
[stars-url]: https://github.com/othneildrew/Best-README-Template/stargazers
[issues-shield]: https://img.shields.io/github/issues/othneildrew/Best-README-Template.svg?style=for-the-badge
[issues-url]: https://github.com/othneildrew/Best-README-Template/issues
[license-shield]: https://img.shields.io/github/license/othneildrew/Best-README-Template.svg?style=for-the-badge
[license-url]: https://github.com/othneildrew/Best-README-Template/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/othneildrew
[product-screenshot]: images/screenshot.png
[Go-url]: https://https://go.dev/
[Next.js]: https://img.shields.io/badge/next.js-000000?style=for-the-badge&logo=nextdotjs&logoColor=white
[Next-url]: https://nextjs.org/
[React.js]: https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB
[React-url]: https://reactjs.org/
[Vue.js]: https://img.shields.io/badge/Vue.js-35495E?style=for-the-badge&logo=vuedotjs&logoColor=4FC08D
[Vue-url]: https://vuejs.org/
[Angular.io]: https://img.shields.io/badge/Angular-DD0031?style=for-the-badge&logo=angular&logoColor=white
[Angular-url]: https://angular.io/
[Svelte.dev]: https://img.shields.io/badge/Svelte-4A4A55?style=for-the-badge&logo=svelte&logoColor=FF3E00
[Svelte-url]: https://svelte.dev/
[Laravel.com]: https://img.shields.io/badge/Laravel-FF2D20?style=for-the-badge&logo=laravel&logoColor=white
[Laravel-url]: https://laravel.com
[Bootstrap.com]: https://img.shields.io/badge/Bootstrap-563D7C?style=for-the-badge&logo=bootstrap&logoColor=white
[Bootstrap-url]: https://getbootstrap.com
[JQuery.com]: https://img.shields.io/badge/jQuery-0769AD?style=for-the-badge&logo=jquery&logoColor=white
[JQuery-url]: https://jquery.com 