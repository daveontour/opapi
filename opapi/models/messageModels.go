package models

import (
	"bufio"
	"fmt"
	"time"
)

type Request struct {
	Direction                  string
	Airline                    string
	FltNum                     string
	From                       string
	To                         string
	UpdatedSince               string
	Route                      string
	UserProfile                UserProfile
	PresentQueryableParameters []ParameterValuePair
}
type Response struct {
	User             string               `json:"User,omitempty"`
	AirportCode      string               `json:"AirportCode,omitempty"`
	Route            string               `json:"Route,omitempty"`
	From             string               `json:"FlightsFrom,omitempty"`
	To               string               `json:"FlightsTo,omitempty"`
	Airline          string               `json:"Airline,omitempty"`
	Flight           string               `json:"Flight,omitempty"`
	FromResource     string               `json:"ResourcessFrom,omitempty"`
	ToResource       string               `json:"ResourceTo,omitempty"`
	NumberOfFlights  int                  `json:"NumberOfFlights,omitempty"`
	Direction        string               `json:"Direction,omitempty"`
	CustomFieldQuery []ParameterValuePair `json:"CustomFieldQueries,omitempty"`
	Warnings         []string             `json:"Warnings,omitempty"`
	Errors           []string             `json:"Errors,omitempty"`
	ResponseFlights  []FlightResponseItem `json:"Flights,omitempty"`
}
type ResourceResponse struct {
	User                string                           `json:"User,omitempty"`
	AirportCode         string                           `json:"AirportCode,omitempty"`
	From                string                           `json:"FlightsFrom,omitempty"`
	To                  string                           `json:"FlightsTo,omitempty"`
	NumberOfFlights     int                              `json:"NumberOfFlights,omitempty"`
	FromResource        string                           `json:"ResourcessFrom,omitempty"`
	ToResource          string                           `json:"ResourceTo,omitempty"`
	Direction           string                           `json:"Direction,omitempty"`
	CustomFieldQuery    []ParameterValuePair             `json:"CustomFieldQueries,omitempty"`
	Warnings            []string                         `json:"Warnings,omitempty"`
	Errors              []string                         `json:"Errors,omitempty"`
	ResourceType        string                           `json:"ResourceType,omitempty"`
	ResourceID          string                           `json:"ResourceID,omitempty"`
	FlightID            string                           `json:"FlightID,omitempty"`
	Airline             string                           `json:"Airline,omitempty"`
	Allocations         []AllocationResponseItem         `json:"Allocations,omitempty"`
	ConfiguredResources []ConfiguredResourceResponseItem `json:"ConfiguredResources,omitempty"`
}
type Repositories struct {
	Repositories []Repository `json:"airports"`
}
type Repository struct {
	AMSAirport                             string `json:"AMSAirport"`
	AMSSOAPServiceURL                      string `json:"AMSSOAPServiceURL"`
	AMSRestServiceURL                      string `json:"AMSRestServiceURL"`
	AMSToken                               string `json:"AMSToken"`
	FlightSDOWindowMinimumInDaysFromNow    int    `json:"FlightSDOWindowMinimumInDaysFromNow"`
	FlightSDOWindowMaximumInDaysFromNow    int    `json:"FlightSDOWindowMaximumInDaysFromNow"`
	ListenerType                           string `json:"ListenerType"`
	RabbitMQConnectionString               string `json:"RabbitMQConnectionString"`
	RabbitMQExchange                       string `json:"RabbitMQExchange"`
	RabbitMQTopic                          string `json:"RabbitMQTopic"`
	PublishChangesRabbitMQConnectionString string `json:"PublishChangesRabbitMQConnectionString"`
	PublishChangesRabbitMQExchange         string `json:"PublishChangesRabbitMQExchange"`
	PublishChangesRabbitMQTopic            string `json:"PublishChangesRabbitMQTopic"`
	NotificationListenerQueue              string `json:"NotificationListenerQueue"`
	LoadFlightChunkSizeInDays              int    `json:"LoadFlightChunkSizeInDays"`
	FlightLinkedList                       FlightLinkedList
	CurrentLowerLimit                      time.Time
	CurrentUpperLimit                      time.Time
	CheckInList                            ResourceLinkedList
	StandList                              ResourceLinkedList
	GateList                               ResourceLinkedList
	CarouselList                           ResourceLinkedList
	ChuteList                              ResourceLinkedList
}
type FixedResources struct {
	Values []FixedResource `xml:"FixedResource"`
}
type FixedResource struct {
	ResourceTypeCode string `xml:"ResourceTypeCode"`
	Name             string `xml:"Name"`
	Area             string `xml:"Area"`
}

type ChangePushJob struct {
	Sub         UserChangeSubscription
	Flight      *Flight
	UserProfile *UserProfile
}
type SchedulePushJob struct {
	Sub         UserPushSubscription
	UserToken   string
	UserName    string
	UserProfile *UserProfile
}
type FlightUpdateChannelMessage struct {
	FlightID    string
	AirportCode string
}

type AllocationItem struct {
	PrevNode             *AllocationItem `xml:"-" json:"-"`
	NextNode             *AllocationItem `xml:"-" json:"-"`
	ResourceID           string
	From                 time.Time
	To                   time.Time
	FlightID             string
	Direction            string
	Route                string
	AircraftType         string
	AircraftRegistration string
	AirportCode          string
	LastUpdate           time.Time
}

type AllocationResponseItem struct {
	ResourceType string `xml:"ResourceTypeCode"`
	Name         string `xml:"Name"`
	Area         string `xml:"Area"`
	AllocationItem
}

type ConfiguredResourceResponseItem struct {
	ResourceTypeCode string `xml:"ResourceTypeCode"`
	Name             string `xml:"Name"`
	Area             string `xml:"Area"`
}

type AllocationLinkedList struct {
	Head *AllocationItem
	Tail *AllocationItem
}

type ResourceAllocationStruct struct {
	PrevNode              *ResourceAllocationStruct
	NextNode              *ResourceAllocationStruct
	Resource              FixedResource
	FlightAllocationsList AllocationLinkedList
}

type ParameterValuePair struct {
	Parameter string `json:"Parameter,omitempty"`
	Value     string `json:"Value,omitempty"`
}

type PropertyValuePair struct {
	Text         string `xml:",chardata"`
	PropertyName string `xml:"propertyName,attr"`
}

type NumberOfAllocations struct {
	ResourceName        string
	NumberofAllocations int
}

type MetricsReport struct {
	Airport                          string
	NumberOfFlights                  int
	NumberOfCheckins                 int
	NumberOfGates                    int
	TotalNumberOfGateAllocations     int
	NumberOfStands                   int
	TotalNumberOfStandAllocations    int
	NumberOfCarousels                int
	TotalNumberOfCarouselAllocations int
	NumberOfChutes                   int
	TotalNumberOfChuteAllocations    int
	TotalNumberOfCheckinAllocations  int
	CheckInAllocationMetrics         []NumberOfAllocations
	GateAllocationMetrics            []NumberOfAllocations
	StandAllocationMetrics           []NumberOfAllocations
	CarouselAllocationMetrics        []NumberOfAllocations
	ChuteAllocationMetrics           []NumberOfAllocations
	// MemAllocMB                       int
	// MemHeapAllocMB                   int
	// MemTotaAllocMB                   int
	// MemSysMB                         int
	// MemNumGC                         int
}
type MetricsReportNow struct {
	Airport                             string
	NumberOfFlights                     int
	NumberOfCheckins                    int
	TotalNumberOfCheckinAllocationsNow  int
	NumberOfGates                       int
	TotalNumberOfGateAllocationsNow     int
	NumberOfStands                      int
	TotalNumberOfStandAllocationsNow    int
	NumberOfCarousels                   int
	TotalNumberOfCarouselAllocationsNow int
	NumberOfChutes                      int
	TotalNumberOfChuteAllocationsNow    int
	CheckInAllocationMetricsNow         []NumberOfAllocations
	GateAllocationMetricsNow            []NumberOfAllocations
	StandAllocationMetricsNow           []NumberOfAllocations
	CarouselAllocationMetricsNow        []NumberOfAllocations
	ChuteAllocationMetricsNow           []NumberOfAllocations
	// MemAllocMB                          int
	// MemHeapAllocMB                      int
	// MemTotaAllocMB                      int
	// MemSysMB                            int
	// MemNumGC                            int
}

type Users struct {
	Users []UserProfile `json:"users"`
}
type UserProfile struct {
	Enabled                      bool                     `json:"Enabled"`
	UserName                     string                   `json:"UserName"`
	Key                          string                   `json:"Key"`
	AllowedAirports              []string                 `json:"AllowedAirports"`
	AllowedAirlines              []string                 `json:"AllowedAirlines"`
	AllowedCustomFields          []string                 `json:"AllowedCustomFields"`
	DefaultAirline               string                   `json:"DefaultAirline"`
	DefaultQueryableCustomFields []ParameterValuePair     `json:"DefaultQueryableCustomFields"`
	UserPushSubscriptions        []UserPushSubscription   `json:"UserPushSubscriptions"`
	UserChangeSubscriptions      []UserChangeSubscription `json:"UserChangeSubscriptions"`
}
type UserPushSubscription struct {
	Enabled                               bool
	EnableInDemoMode                      bool
	PushOnStartUp                         bool
	Airport                               string
	DestinationURL                        string
	HeaderParameters                      []ParameterValuePair
	SubscriptionType                      string
	Time                                  string
	ReptitionHours                        int
	ReptitionMinutes                      int
	Airline                               string
	From                                  int
	To                                    int
	QueryableCustomFields                 []ParameterValuePair
	ResourceType                          string
	ResourceID                            string
	Route                                 string
	Direction                             string
	TrustBadCertificates                  bool
	PublishStatusRabbitMQConnectionString string
	PublishStatusRabbitMQExchange         string
	PublishStatusRabbitMQTopic            string
	HTTPEnabled                           bool
	RMQEnabled                            bool
}
type UserChangeSubscription struct {
	Enabled                                bool
	Airport                                string
	DestinationURL                         string
	HeaderParameters                       []ParameterValuePair
	CheckInChange                          bool
	GateChange                             bool
	StandChange                            bool
	CarouselChange                         bool
	ChuteChange                            bool
	AircraftTypeOrRegoChange               bool
	RouteChange                            bool
	LinkedFlightChange                     bool
	EventChange                            bool
	CreateFlight                           bool
	DeleteFlight                           bool
	UpdateFlight                           bool
	All                                    bool
	ParameterChange                        []string
	TrustBadCertificates                   bool
	UserKey                                string
	PublishChangesRabbitMQConnectionString string
	PublishChangesRabbitMQExchange         string
	PublishChangesRabbitMQTopic            string
	HTTPEnabled                            bool
	RMQEnabled                             bool
}

type ResourceLinkedList struct {
	Head *ResourceAllocationStruct
	Tail *ResourceAllocationStruct
}

type FlightResponseItem struct {
	// This structure is used as an element in an array for the flights to be returned in a response
	// The array is used so the STO can be used to sort the returned order

	FlightPtr *Flight
	STO       time.Time
}

type FlightLinkedList struct {
	Head *Flight
	Tail *Flight
}

type GetFlightsError struct {
	StatusCode int
	Err        error
}

func (r *Response) AddWarning(w string) {
	r.Warnings = append(r.Warnings, w)
}
func (r *Response) AddError(w string) {
	r.Errors = append(r.Errors, w)
}
func (r *ResourceResponse) AddWarning(w string) {
	r.Warnings = append(r.Warnings, w)
}
func (r *ResourceResponse) AddError(w string) {
	r.Errors = append(r.Errors, w)
}

// JSON Reeiver function to write the JSON on an structure

func (d ResourceResponse) WriteJSON(fwb *bufio.Writer) (err error) {

	_, err = fwb.WriteString("{" +
		"\"Airport\":\"" + d.AirportCode + "\"," +
		"\"ResourceType\":\"" + d.ResourceType + "\"," +
		"\"ResourceName\":\"" + d.ResourceID + "\"," +
		"\"AllocationStart\":\"" + d.FromResource + "\"," +
		"\"AllocationEnd\":\"" + d.ToResource + "\"," +
		"\"FlightNumber\":\"" + d.FlightID + "\"," +
		"\"Airline\":\"" + d.Airline + "\"," +
		"\"CustomFieldQuery\":[")
	if err != nil {
		return
	}

	for idx, w := range d.CustomFieldQuery {
		if idx > 0 {
			_, err = fwb.WriteString(",")
			if err != nil {
				return
			}
		}
		_, err = fwb.WriteString("{\"" + w.Parameter + "\":\"" + w.Value + "\"}")
		if err != nil {
			return
		}
	}
	_, err = fwb.WriteString("],")
	if err != nil {
		return
	}

	_, err = fwb.WriteString("\"Warnings\":[")
	if err != nil {
		return
	}
	for idx, w := range d.Warnings {
		if idx > 0 {
			_, err = fwb.WriteString(",")
			if err != nil {
				return
			}
		}
		_, err = fwb.WriteString("\"" + w + "\"")
		if err != nil {
			return
		}
	}

	_, err = fwb.WriteString("],\"Errors\":[")
	if err != nil {
		return
	}
	for idx, w := range d.Errors {
		if idx > 0 {
			_, err = fwb.WriteString(",")
			if err != nil {
				return
			}
		}
		_, err = fwb.WriteString("\"" + w + "\"")
		if err != nil {
			return
		}
	}

	_, err = fwb.WriteString("],\"Allocations\": [")
	if err != nil {
		return
	}
	for idx, a := range d.Allocations {
		if idx > 0 {
			_, err = fwb.WriteString(",")
			if err != nil {
				return
			}
		}
		err = a.WriteJSON(fwb)
		if err != nil {
			return
		}
	}
	_, err = fwb.WriteString("]}")
	if err != nil {
		return
	}

	return nil
}

func WriteJSONConfiguredResources(fwb *bufio.Writer, d ResourceResponse) (err error) {

	_, err = fwb.WriteString("{" +
		"\"Airport\":\"" + d.AirportCode + "\"," +
		"\"ResourceType\":\"" + d.ResourceType + "\"")
	if err != nil {
		return
	}

	_, err = fwb.WriteString(",\"ConfiguredResources\": [")
	if err != nil {
		return
	}

	for idx, a := range d.ConfiguredResources {
		if idx > 0 {
			_, err = fwb.WriteString(",")
			if err != nil {
				break
			}
		}
		fwb.WriteString("{\"ResourceType\":\"" + a.ResourceTypeCode + "\",\"ResourceName\":\"" + a.Name + "\",\"Area\":\"" + a.Area + "\"}")
	}
	_, err = fwb.WriteString("]}")
	if err != nil {
		return
	}

	return nil
}
func (d AllocationResponseItem) WriteJSON(fwb *bufio.Writer) error {

	_, err := fwb.WriteString("{" +
		"\"ResourceType\":\"" + d.ResourceType + "\"," +
		"\"Name\":\"" + d.Name + "\"," +
		"\"Area\":\"" + d.Area + "\"," +
		fmt.Sprintf("\"AllocationStart\":\"%s\",", d.AllocationItem.From) +
		fmt.Sprintf("\"AllocationEnd\":\"%s\",", d.AllocationItem.To) +
		"\"Flight\": {" +
		"\"FlightID\":\"" + d.AllocationItem.FlightID + "\"," +
		"\"Direction\":\"" + d.AllocationItem.Direction + "\"," +
		"\"Route\":\"" + d.AllocationItem.Route + "\"")
	if err != nil {
		return err
	}
	if d.AllocationItem.AircraftRegistration != "" {
		_, err = fwb.WriteString(",\"AircraftRegistration\":\"" + d.AllocationItem.AircraftRegistration + "\"")
		if err != nil {
			return err
		}
	}
	if d.AllocationItem.AircraftType != "" {
		_, err = fwb.WriteString(",\"AircraftType\":\"" + d.AllocationItem.AircraftType + "\"")
		if err != nil {
			return err
		}
	}
	_, err = fwb.WriteString(" }}")
	if err != nil {
		return err
	}

	return nil
}
func WriteFlightsInJSON(fwb *bufio.Writer, flights []FlightResponseItem, userProfile *UserProfile, statusOnly bool, maxFlights int) error {
	_, err := fwb.WriteString(`"Flights":[`)
	if err != nil {
		return err
	}
	for idx, currentNode := range flights {
		if maxFlights > 0 && idx > maxFlights {
			break
		}
		if idx > 0 {
			err = fwb.WriteByte(',')
			if err != nil {
				return err
			}
		}
		err = currentNode.FlightPtr.WriteJSON(fwb, userProfile, statusOnly)
		if err != nil {
			return err
		}
	}

	err = fwb.WriteByte(']')
	if err != nil {
		return err
	}
	return nil
}

func (r *GetFlightsError) Error() string {
	return fmt.Sprintf("status %d: err %v", r.StatusCode, r.Err)
}

func (ll *AllocationLinkedList) RemoveFlightAllocations(flightID string) {
	currentNode := ll.Head

	for currentNode != nil {
		if currentNode.FlightID == flightID {

			if currentNode.PrevNode != nil {
				currentNode.PrevNode.NextNode = currentNode.NextNode
			} else {
				ll.Head = currentNode.NextNode
			}

			if currentNode.NextNode != nil {
				currentNode.NextNode.PrevNode = currentNode.PrevNode
			} else {
				ll.Tail = currentNode.PrevNode
			}

			currentNode.PrevNode = nil
			currentNode.NextNode = nil

			//return // Node found and removed, exit the function
		}

		currentNode = currentNode.NextNode
	}
}
func (ll *AllocationLinkedList) Len() int {
	currentNode := ll.Head
	count := 0

	for currentNode != nil {
		count++
		currentNode = currentNode.NextNode
	}

	return count
}
func (ll *AllocationLinkedList) AddNode(newNode AllocationItem) {

	newNode.PrevNode = ll.Tail
	newNode.NextNode = nil

	if ll.Tail != nil {
		ll.Tail.NextNode = &newNode
	}

	ll.Tail = &newNode

	if ll.Head == nil {
		ll.Head = &newNode
	}
}
func (n1 *FixedResource) Equals(n2 *FixedResource) bool {

	if n1.Area != n2.Area {
		return false
	}
	if n1.Name != n2.Name {
		return false
	}
	if n1.ResourceTypeCode != n2.ResourceTypeCode {
		return false
	}
	return true
}
func (ll *ResourceLinkedList) AddAllocation(node AllocationItem) {
	currentNode := ll.Head

	for currentNode != nil {
		if currentNode.Resource.Name == node.ResourceID {
			currentNode.FlightAllocationsList.AddNode(node)
			break
		}
		currentNode = currentNode.NextNode
	}
}
func (ll *ResourceLinkedList) ReplaceOrAddNodes(nodes []FixedResource) {
	for _, node := range nodes {
		ll.ReplaceOrAddNode(node)
	}
}
func (ll *ResourceLinkedList) ReplaceOrAddNode(fr FixedResource) {
	currentNode := ll.Head
	node := ResourceAllocationStruct{Resource: fr}

	for currentNode != nil {
		if currentNode.Resource.Equals(&node.Resource) {
			// Don't need to do anything
			return
		}
		currentNode = currentNode.NextNode
	}

	ll.AddNode(node)
}
func (ll *ResourceLinkedList) AddNode(newNode ResourceAllocationStruct) {

	newNode.PrevNode = ll.Tail
	newNode.NextNode = nil

	if ll.Tail != nil {
		ll.Tail.NextNode = &newNode
	}

	ll.Tail = &newNode

	if ll.Head == nil {
		ll.Head = &newNode
	}
}
func (ll *ResourceLinkedList) RemoveFlightAllocation(flightID string) {
	currentNode := ll.Head

	for currentNode != nil {
		currentNode.FlightAllocationsList.RemoveFlightAllocations(flightID)
		currentNode = currentNode.NextNode
	}
}
func (ll *ResourceLinkedList) Len() int {
	currentNode := ll.Head
	count := 0

	for currentNode != nil {
		count++
		currentNode = currentNode.NextNode
	}

	return count
}
func (ll *ResourceLinkedList) NumberOfFlightAllocations() (n int) {
	currentNode := ll.Head

	for currentNode != nil {
		n = n + currentNode.FlightAllocationsList.Len()
		currentNode = currentNode.NextNode
	}
	return
}
func (ll *ResourceLinkedList) NumberOfFlightAllocationsNow() (n int) {
	currentNode := ll.Head
	now := time.Now()

	for currentNode != nil {
		x := currentNode.FlightAllocationsList.Head

		for x != nil {
			if inTimeSpan(x.From, x.To, now) {
				n++
			}
			x = x.NextNode
		}

		currentNode = currentNode.NextNode
	}
	return
}

func (ll *ResourceLinkedList) AllocationsMetrics() (metrics []NumberOfAllocations) {

	currentNode := ll.Head
	for currentNode != nil {
		n := currentNode.FlightAllocationsList.Len()
		metrics = append(metrics, NumberOfAllocations{ResourceName: currentNode.Resource.Name, NumberofAllocations: n})
		currentNode = currentNode.NextNode
	}

	return
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func (ll *ResourceLinkedList) AllocationsMetricsNow() (metrics []NumberOfAllocations) {

	now := time.Now()
	currentNode := ll.Head
	for currentNode != nil {
		n := 0

		x := currentNode.FlightAllocationsList.Head

		for x != nil {
			if inTimeSpan(x.From, x.To, now) {
				n++
			}
			x = x.NextNode
		}

		metrics = append(metrics, NumberOfAllocations{ResourceName: currentNode.Resource.Name, NumberofAllocations: n})
		currentNode = currentNode.NextNode
	}

	return
}

func (ll *FlightLinkedList) RemoveNode(removeNode Flight) {
	currentNode := ll.Head

	for currentNode != nil {
		if currentNode.GetFlightID() == removeNode.GetFlightID() {
			if currentNode.PrevNode != nil {
				currentNode.PrevNode.NextNode = currentNode.NextNode
			} else {
				ll.Head = currentNode.NextNode
			}

			if currentNode.NextNode != nil {
				currentNode.NextNode.PrevNode = currentNode.PrevNode
			} else {
				ll.Tail = currentNode.PrevNode
			}

			currentNode.PrevNode = nil
			currentNode.NextNode = nil

			return // Node found and removed, exit the function
		}

		currentNode = currentNode.NextNode
	}
}

func (ll *FlightLinkedList) GetFlight(flightID string) *Flight {
	currentNode := ll.Head

	for currentNode != nil {
		if currentNode.GetFlightID() == flightID {
			return currentNode
		}

		currentNode = currentNode.NextNode
	}
	return nil
}

func (rep *Repository) GetFlight(flightID string) *Flight {
	return rep.FlightLinkedList.GetFlight(flightID)
}

func (ll *FlightLinkedList) Len() int {
	currentNode := ll.Head
	count := 0

	for currentNode != nil {
		count++
		currentNode = currentNode.NextNode
	}

	return count
}

func (ll *FlightLinkedList) AddNode(newNode Flight) {
	// AddNode adds a new node to the end of the doubly linked list.
	newNode.PrevNode = ll.Tail
	newNode.NextNode = nil

	if ll.Tail != nil {
		ll.Tail.NextNode = &newNode
	}

	ll.Tail = &newNode

	if ll.Head == nil {
		ll.Head = &newNode
	}
}
func (ll *FlightLinkedList) ReplaceOrAddNode(node Flight) {
	currentNode := ll.Head

	for currentNode != nil {
		if currentNode.GetFlightID() == node.GetFlightID() {

			// Replace the entire node
			node.PrevNode = currentNode.PrevNode
			node.NextNode = currentNode.NextNode

			if currentNode.PrevNode != nil {
				currentNode.PrevNode.NextNode = &node
			} else {
				ll.Head = &node
			}

			if currentNode.NextNode != nil {
				currentNode.NextNode.PrevNode = &node
			} else {
				ll.Tail = &node
			}

			currentNode.PrevNode = nil
			currentNode.NextNode = nil

			// Node found and replaced, exit the function
			return
		}
		currentNode = currentNode.NextNode
	}

	ll.AddNode(node)
}
func (ll *FlightLinkedList) RemoveExpiredNodes(from time.Time) {
	currentNode := ll.Head

	for currentNode != nil {
		if currentNode.GetSDO().Before(from) {
			if currentNode.PrevNode != nil {
				currentNode.PrevNode.NextNode = currentNode.NextNode
			} else {
				ll.Head = currentNode.NextNode
			}

			if currentNode.NextNode != nil {
				currentNode.NextNode.PrevNode = currentNode.PrevNode
			} else {
				ll.Tail = currentNode.PrevNode
			}

			currentNode = currentNode.NextNode

		} else {
			currentNode = currentNode.NextNode
		}
	}
}

func (r *Repository) RemoveFlightAllocation(flightID string) {
	r.CheckInList.RemoveFlightAllocation(flightID)
	r.GateList.RemoveFlightAllocation(flightID)
	r.StandList.RemoveFlightAllocation(flightID)
	r.CarouselList.RemoveFlightAllocation(flightID)
	r.ChuteList.RemoveFlightAllocation(flightID)
}
func (r *Repository) UpdateLowerLimit(t time.Time) {
	r.CurrentLowerLimit = t
}
func (r *Repository) UpdateUpperLimit(t time.Time) {
	r.CurrentUpperLimit = t
}
