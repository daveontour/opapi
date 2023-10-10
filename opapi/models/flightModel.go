package models

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/daveontour/opapi/opapi/timeservice"
)

type AirlineDesignator struct {
	CodeContext string `xml:"codeContext,attr"`
	Text        string `xml:",chardata"`
}

type AirportCode struct {
	CodeContext string `xml:"codeContext,attr"`
	Text        string `xml:",chardata"`
}

type FlightId struct {
	FlightKind        string              `xml:"FlightKind"`
	AirlineDesignator []AirlineDesignator `xml:"AirlineDesignator"`
	FlightNumber      string              `xml:"FlightNumber"`
	ScheduledDate     string              `xml:"ScheduledDate"`
	AirportCode       []AirportCode       `xml:"AirportCode"`
}

func (d *FlightId) WriteJSON(fwb *bufio.Writer) error {

	_, err := fwb.WriteString("\"FlightId\":{" +
		"\"FlightKind\":\"" + d.FlightKind + "\"," +
		"\"FlightNumber\":\"" + d.FlightNumber + "\"," +
		"\"ScheduledDate\":\"" + string(d.ScheduledDate) + "\"")

	if err != nil {
		return err
	}

	if d.AirportCode != nil {
		_, err = fwb.WriteString(",\"AirportCode\":{")
		if err != nil {
			return err
		}
		for idx, apt := range d.AirportCode {
			if idx > 0 {
				_, err = fwb.WriteString(",")
				if err != nil {
					return err
				}
			}
			_, err = fwb.WriteString("\"" + apt.CodeContext + "\":\"" + apt.Text + "\"")
			if err != nil {
				return err
			}
		}
		_, err = fwb.WriteString("}")
		if err != nil {
			return err
		}
	}
	if d.AirlineDesignator != nil {
		_, err = fwb.WriteString(",\"AirlineDesignator\":{")
		if err != nil {
			return err
		}
		for idx, al := range d.AirlineDesignator {
			if idx > 0 {
				_, err = fwb.WriteString(",")
				if err != nil {
					return err
				}
			}
			_, err = fwb.WriteString("\"" + al.CodeContext + "\":\"" + al.Text + "\"")
			if err != nil {
				return err
			}
		}
		_, err = fwb.WriteString("}")
		if err != nil {
			return err
		}
	}

	_, err = fwb.WriteString("}")
	if err != nil {
		return err
	}
	return nil
}

func (d *FlightId) MarshalJSON() ([]byte, error) {

	fwb := strings.Builder{}
	fwb.WriteString("{" +
		"\"FlightKind\":\"" + d.FlightKind + "\"," +
		"\"FlightNumber\":\"" + d.FlightNumber + "\"," +
		"\"ScheduledDate\":\"" + string(d.ScheduledDate) + "\"")

	if d.AirportCode != nil {
		fwb.WriteString(",\"AirportCode\":{")
		for idx, apt := range d.AirportCode {
			if idx > 0 {
				fwb.WriteString(",")
			}
			fwb.WriteString("\"" + apt.CodeContext + "\":\"" + apt.Text + "\"")
		}
		fwb.WriteString("}")
	}
	if d.AirlineDesignator != nil {
		fwb.WriteString(",\"AirlineDesignator\":{")
		for idx, al := range d.AirlineDesignator {
			if idx > 0 {
				fwb.WriteString(",")
			}
			fwb.WriteString("\"" + al.CodeContext + "\":\"" + al.Text + "\"")
		}
		fwb.WriteString("}")
	}

	fwb.WriteString("}")

	var sendText = fwb.String()

	return []byte(sendText), nil
}

type Value struct {
	PropertyName string `xml:"propertyName,attr"`
	Text         string `xml:",chardata"`
}

func (d *Value) WriteJSON(fwb *bufio.Writer) error {
	_, err := fwb.WriteString("{\"" + d.PropertyName + "\":\"" + d.Text + "\"}")
	if err != nil {
		return err
	}
	return nil
}

type LinkedFlight struct {
	FlightId FlightId `xml:"FlightId"`
	Value    []Value  `xml:"Value"`
}

func (d *LinkedFlight) WriteJSON(fwb *bufio.Writer) error {

	_, err := fwb.WriteString("{\"FlightId\":{" +
		"\"FlightKind\":\"" + d.FlightId.FlightKind + "\"," +
		"\"FlightNumber\":\"" + d.FlightId.FlightNumber + "\"," +
		"\"ScheduledDate\":\"" + string(d.FlightId.ScheduledDate) + "\"," +
		"\"AirportCode\":{")
	if err != nil {
		return err
	}
	for idx, apt := range d.FlightId.AirportCode {
		if idx > 0 {
			_, err = fwb.WriteString(",")
			if err != nil {
				return err
			}
		}
		_, err = fwb.WriteString("\"" + apt.CodeContext + "\":\"" + apt.Text + "\"")
		if err != nil {
			return err
		}
	}
	_, err = fwb.WriteString("},")
	if err != nil {
		return err
	}

	_, err = fwb.WriteString("\"AirlineDesignator\":{")
	if err != nil {
		return err
	}
	for idx, al := range d.FlightId.AirlineDesignator {
		if idx > 0 {
			_, err = fwb.WriteString(",")
			if err != nil {
				return err
			}
		}
		_, err = fwb.WriteString("\"" + al.CodeContext + "\":\"" + al.Text + "\"")
		if err != nil {
			return err
		}
	}
	_, err = fwb.WriteString("}" +
		"}," +
		"\"Values\":")
	if err != nil {
		return err
	}
	err = MarshalValuesArrayJSON(d.Value, fwb)
	if err != nil {
		return err
	}

	_, err = fwb.WriteString("}")
	if err != nil {
		return err
	}

	return nil
}

type AircraftTypeCode struct {
	CodeContext string `xml:"codeContext,attr"`
	Text        string `xml:",chardata"`
}
type AircraftTypeId struct {
	//	Text             string             `xml:",chardata" json:"-"`
	AircraftTypeCode []AircraftTypeCode `xml:"AircraftTypeCode"`
}

func (tid *AircraftTypeId) WriteJSON(fwb *bufio.Writer) error {

	_, err := fwb.WriteString("{")
	if err != nil {
		return err
	}

	if tid.AircraftTypeCode != nil {
		_, err = fwb.WriteString("\"AircraftTypeCode\":{")
		if err != nil {
			return err
		}

		for idx, tc := range tid.AircraftTypeCode {
			if idx > 0 {
				_, err = fwb.WriteString(",")
				if err != nil {
					return err
				}
			}
			_, err = fwb.WriteString("\"" + tc.CodeContext + "\":\"" + tc.Text + "\"")
			if err != nil {
				return err
			}
		}
		_, err = fwb.WriteString("}")
		if err != nil {
			return err
		}
	}
	_, err = fwb.WriteString("}")
	if err != nil {
		return err
	}
	return nil
}

type AircraftType struct {
	AircraftTypeId AircraftTypeId `xml:"AircraftTypeId"`
	Value          []Value        `xml:"Value"`
}

func (t *AircraftType) WriteJSON(fwb *bufio.Writer) error {

	_, err := fwb.WriteString("{")
	if err != nil {
		return err
	}
	_, err = fwb.WriteString("\"AircraftTypeId\":")
	if err != nil {
		return err
	}
	err = t.AircraftTypeId.WriteJSON(fwb)
	if err != nil {
		return err
	}
	if len(t.Value) > 0 {
		_, err = fwb.WriteString(",")
		if err != nil {
			return err
		}
	}

	_, err = fwb.WriteString("\"Values\":")
	if err != nil {
		return err
	}
	err = MarshalValuesArrayJSON(t.Value, fwb)
	if err != nil {
		return err
	}
	_, err = fwb.WriteString("}")
	if err != nil {
		return err
	}
	return nil
}

type RouteViaPoint struct {
	SequenceNumber string        `xml:"sequenceNumber,attr"`
	AirportCode    []AirportCode `xml:"AirportCode"`
}

type ViaPoints struct {
	RouteViaPoint []RouteViaPoint `xml:"RouteViaPoint"`
}

func (r *ViaPoints) WriteJSON(fwb *bufio.Writer) error {

	_, err := fwb.WriteString("[")
	if err != nil {
		return err
	}

	for idx, rvp := range r.RouteViaPoint {
		if idx > 0 {
			_, err = fwb.WriteString(",")
			if err != nil {
				return err
			}
		}
		_, err = fwb.WriteString("{" +
			"\"SequenceNumber\":\"" + rvp.SequenceNumber + "\"," +
			"\"AirportCode\":{")
		if err != nil {
			return err
		}

		for idx2, apt := range rvp.AirportCode {
			if idx2 > 0 {
				_, err = fwb.WriteString(",")
				if err != nil {
					return err
				}
			}
			_, err = fwb.WriteString("\"" + apt.CodeContext + "\":\"" + apt.Text + "\"")
			if err != nil {
				return err
			}
		}

		_, err = fwb.WriteString("}}")
		if err != nil {
			return err
		}
	}

	_, err = fwb.WriteString("]")
	if err != nil {
		return err
	}

	return nil
}

func (r *ViaPoints) MarshalJSON() ([]byte, error) {

	fwb := strings.Builder{}

	_, err := fwb.WriteString("[")
	if err != nil {
		return nil, err
	}

	for idx, rvp := range r.RouteViaPoint {
		if idx > 0 {
			_, err = fwb.WriteString(",")
			if err != nil {
				return nil, err
			}
		}
		_, err = fwb.WriteString("{" +
			"\"SequenceNumber\":\"" + rvp.SequenceNumber + "\"," +
			"\"AirportCode\":{")
		if err != nil {
			return nil, err
		}

		for idx2, apt := range rvp.AirportCode {
			if idx2 > 0 {
				_, err = fwb.WriteString(",")
				if err != nil {
					return nil, err
				}
			}
			_, err = fwb.WriteString("\"" + apt.CodeContext + "\":\"" + apt.Text + "\"")
			if err != nil {
				return nil, err
			}
		}

		_, err = fwb.WriteString("}}")
		if err != nil {
			return nil, err
		}
	}

	_, err = fwb.WriteString("]")
	if err != nil {
		return nil, err
	}

	var sendText = fwb.String()

	return []byte(sendText), nil
}

type Route struct {
	CustomsType string    `xml:"customsType,attr"`
	ViaPoints   ViaPoints `xml:"ViaPoints"`
}

func (r *Route) WriteJSON(fwb *bufio.Writer) error {

	_, err := fwb.WriteString("{")
	if err != nil {
		return err
	}
	if r.CustomsType != "" {
		_, err = fwb.WriteString("\"CustomType\":\"" + r.CustomsType + "\",")
		if err != nil {
			return err
		}
	}
	_, err = fwb.WriteString("\"ViaPoints\":")
	if err != nil {
		return err
	}
	err = r.ViaPoints.WriteJSON(fwb)
	if err != nil {
		return err
	}
	_, err = fwb.WriteString("}")
	if err != nil {
		return err
	}
	return nil
}

type TableValue struct {
	PropertyName string  `xml:"propertyName,attr"`
	Value        []Value `xml:"Value"`
}

func (ss *StandSlots) WriteJSON(fwb *bufio.Writer) error {

	_, err := fwb.WriteString("[")
	if err != nil {
		return err
	}

	for idx2, s := range ss.StandSlot {

		if idx2 > 0 {
			_, err = fwb.WriteString(",")
			if err != nil {
				return err
			}
		}
		_, err = fwb.WriteString("{")
		if err != nil {
			return err
		}
		for idx3, v := range s.Value {
			if idx3 > 0 {
				_, err = fwb.WriteString(",")
				if err != nil {
					return err
				}
			}
			_, err = fwb.WriteString("\"" + v.PropertyName + "\":\"" + v.Text + "\"")
			if err != nil {
				return err
			}
		}
		for _, v := range s.Stand.Value {
			_, err = fwb.WriteString(",\"" + v.PropertyName + "\":\"" + v.Text + "\"")
			if err != nil {
				return err
			}
		}

		for _, v := range s.Stand.Area.Value {
			_, err = fwb.WriteString(",\"Area" + v.PropertyName + "\":\"" + v.Text + "\"")
			if err != nil {
				return err
			}
		}
		_, err = fwb.WriteString("}")
		if err != nil {
			return err
		}

	}

	_, err = fwb.WriteString("]")
	if err != nil {
		return err
	}

	return nil
}

func (ss *CarouselSlots) WriteJSON(fwb *bufio.Writer) error {

	_, err := fwb.WriteString("[")
	if err != nil {
		return err
	}

	for idx2, s := range ss.CarouselSlot {

		if idx2 > 0 {
			_, err = fwb.WriteString(",")
			if err != nil {
				return err
			}
		}
		_, err = fwb.WriteString("{")
		if err != nil {
			return err
		}
		for idx3, v := range s.Value {
			if idx3 > 0 {
				_, err = fwb.WriteString(",")
				if err != nil {
					return err
				}
			}
			_, err = fwb.WriteString("\"" + v.PropertyName + "\":\"" + v.Text + "\"")
			if err != nil {
				return err
			}
		}
		for _, v := range s.Carousel.Value {
			_, err = fwb.WriteString(",\"" + v.PropertyName + "\":\"" + v.Text + "\"")
			if err != nil {
				return err
			}
		}

		for _, v := range s.Carousel.Area.Value {
			_, err = fwb.WriteString(",\"Area" + v.PropertyName + "\":\"" + v.Text + "\"")
			if err != nil {
				return err
			}
		}
		_, err = fwb.WriteString("}")
		if err != nil {
			return err
		}

	}

	_, err = fwb.WriteString("]")
	if err != nil {
		return err
	}

	return nil
}

func (ss *GateSlots) WriteJSON(fwb *bufio.Writer) error {

	_, err := fwb.WriteString("[")
	if err != nil {
		return err
	}
	for idx2, s := range ss.GateSlot {

		if idx2 > 0 {
			_, err = fwb.WriteString(",")
			if err != nil {
				return err
			}
		}
		_, err = fwb.WriteString("{")
		if err != nil {
			return err
		}
		for idx3, v := range s.Value {
			if idx3 > 0 {
				_, err = fwb.WriteString(",")
				if err != nil {
					return err
				}
			}
			_, err = fwb.WriteString("\"" + v.PropertyName + "\":\"" + v.Text + "\"")
			if err != nil {
				return err
			}
		}
		for _, v := range s.Gate.Value {
			_, err = fwb.WriteString(",\"" + v.PropertyName + "\":\"" + v.Text + "\"")
			if err != nil {
				return err
			}
		}

		for _, v := range s.Gate.Area.Value {
			_, err = fwb.WriteString(",\"Area" + v.PropertyName + "\":\"" + v.Text + "\"")
			if err != nil {
				return err
			}
		}

		_, err = fwb.WriteString("}")
		if err != nil {
			return err
		}

	}

	_, err = fwb.WriteString("]")
	if err != nil {
		return err
	}

	return nil
}

func (ss *GateSlots) MarshalJSON() ([]byte, error) {

	fwb := strings.Builder{}

	_, err := fwb.WriteString("{")
	if err != nil {
		return nil, err
	}
	for idx2, s := range ss.GateSlot {

		if idx2 > 0 {
			_, err = fwb.WriteString(",")
			if err != nil {
				return nil, err
			}
		}
		fwb.WriteString("{")
		for idx3, v := range s.Value {
			if idx3 > 0 {
				_, err = fwb.WriteString(",")
				if err != nil {
					return nil, err
				}
			}
			_, err = fwb.WriteString("\"" + v.PropertyName + "\":\"" + v.Text + "\"")
			if err != nil {
				return nil, err
			}
		}
		for _, v := range s.Gate.Value {
			_, err = fwb.WriteString(",\"" + v.PropertyName + "\":\"" + v.Text + "\"")
			if err != nil {
				return nil, err
			}
		}

		for _, v := range s.Gate.Area.Value {
			_, err = fwb.WriteString(",\"Area" + v.PropertyName + "\":\"" + v.Text + "\"")
			if err != nil {
				return nil, err
			}
		}

		_, err = fwb.WriteString("}")
		if err != nil {
			return nil, err
		}

	}

	_, err = fwb.WriteString("]")
	if err != nil {
		return nil, err
	}

	var sendText = fwb.String()

	return []byte(sendText), nil
}

func (ss *CheckInSlots) WriteJSON(fwb *bufio.Writer) error {

	_, err := fwb.WriteString("[")
	if err != nil {
		return err
	}

	for idx2, s := range ss.CheckInSlot {

		if idx2 > 0 {
			_, err = fwb.WriteString(",")
			if err != nil {
				return err
			}
		}
		_, err = fwb.WriteString("{")
		if err != nil {
			return err
		}
		for idx3, v := range s.Value {
			if idx3 > 0 {
				_, err = fwb.WriteString(",")
				if err != nil {
					return err
				}
			}
			_, err = fwb.WriteString("\"" + v.PropertyName + "\":\"" + v.Text + "\"")
			if err != nil {
				return err
			}
		}
		for _, v := range s.CheckIn.Value {
			_, err = fwb.WriteString(",\"" + v.PropertyName + "\":\"" + v.Text + "\"")
			if err != nil {
				return err
			}
		}
		for _, v := range s.CheckIn.Area.Value {
			_, err = fwb.WriteString(",\"Area" + v.PropertyName + "\":\"" + v.Text + "\"")
			if err != nil {
				return err
			}
		}
		_, err = fwb.WriteString("}")
		if err != nil {
			return err
		}
	}

	_, err = fwb.WriteString("]")
	if err != nil {
		return err
	}

	return nil
}

func (s *CheckInSlot) MarshalJSON() ([]byte, error) {

	fwb := strings.Builder{}
	_, err := fwb.WriteString("{")
	if err != nil {
		return nil, err
	}

	for idx3, v := range s.Value {
		if idx3 > 0 {
			_, err = fwb.WriteString(",")
			if err != nil {
				return nil, err
			}
		}
		_, err = fwb.WriteString("\"" + v.PropertyName + "\":\"" + v.Text + "\"")
		if err != nil {
			return nil, err
		}
	}
	for _, v := range s.CheckIn.Value {
		_, err = fwb.WriteString(",\"" + v.PropertyName + "\":\"" + v.Text + "\"")
		if err != nil {
			return nil, err
		}
	}
	for _, v := range s.CheckIn.Area.Value {
		_, err = fwb.WriteString(",\"Area" + v.PropertyName + "\":\"" + v.Text + "\"")
		if err != nil {
			return nil, err
		}
	}
	_, err = fwb.WriteString("}")
	if err != nil {
		return nil, err
	}

	var sendText = fwb.String()

	return []byte(sendText), nil
}
func (ss *ChuteSlots) WriteJSON(fwb *bufio.Writer) error {

	_, err := fwb.WriteString("[")
	if err != nil {
		return err
	}

	for idx2, s := range ss.ChuteSlot {

		if idx2 > 0 {
			_, err = fwb.WriteString(",")
			if err != nil {
				return err
			}

		}
		_, err = fwb.WriteString("{")
		if err != nil {
			return err
		}
		for idx3, v := range s.Value {
			if idx3 > 0 {
				_, err = fwb.WriteString(",")
				if err != nil {
					return err
				}

			}
			_, err = fwb.WriteString("\"" + v.PropertyName + "\":\"" + v.Text + "\"")
			if err != nil {
				return err
			}
		}
		for _, v := range s.Chute.Value {
			_, err = fwb.WriteString(",\"" + v.PropertyName + "\":\"" + v.Text + "\"")
			if err != nil {
				return err
			}
		}
		for _, v := range s.Chute.Area.Value {
			_, err = fwb.WriteString(",\"Area" + v.PropertyName + "\":\"" + v.Text + "\"")
			if err != nil {
				return err
			}
		}
		_, err = fwb.WriteString("}")
		if err != nil {
			return err
		}

	}

	_, err = fwb.WriteString("]")
	if err != nil {
		return err
	}

	return nil
}

type AircraftId struct {
	Registration string `xml:"Registration" json:"Registration" `
}
type Aircraft struct {
	AircraftId AircraftId `xml:"AircraftId" json:"AircraftId"`
}

func (d *Aircraft) WriteJSON(fwb *bufio.Writer) error {

	_, err := fwb.WriteString("{\"AircraftId\": {\"Registration\": \"" + d.AircraftId.Registration + "\"}}")
	if err != nil {
		return err
	}
	return nil
}

type FlightState struct {
	ScheduledTime string        `xml:"ScheduledTime" `
	LinkedFlight  LinkedFlight  `xml:"LinkedFlight"`
	AircraftType  AircraftType  `xml:"AircraftType"`
	Aircraft      Aircraft      `xml:"Aircraft" json:"Aircraft"`
	Route         Route         `xml:"Route" json:"-"`
	Values        []Value       `xml:"Value" json:"Values,omitempty"`
	TableValue    []TableValue  `xml:"TableValue" json:"TableValues,omitempty"`
	StandSlots    StandSlots    `xml:"StandSlots" json:"StandSlots,omitempty"`
	CarouselSlots CarouselSlots `xml:"CarouselSlots" json:"CarouselSlots,omitempty"`
	GateSlots     GateSlots     `xml:"GateSlots" json:"GateSlots,omitempty"`
	CheckInSlots  CheckInSlots  `xml:"CheckInSlots" json:"CheckInSlots,omitempty"`
	ChuteSlots    ChuteSlots    `xml:"ChuteSlots" json:"ChuteSlots,omitempty"`
}

func MarshalValuesArrayJSON(vs []Value, fwb *bufio.Writer) (err error) {

	_, err = fwb.WriteString("{")
	if err != nil {
		return
	}

	set := false
	for _, f := range vs {
		if set {
			_, err = fwb.WriteString(",")
			if err != nil {
				return
			}
		}
		set = true
		_, err = fwb.WriteString("\"" + f.PropertyName + "\":\"" + strings.Replace(f.Text, "\n", "", -1) + "\"")
		if err != nil {
			return
		}

	}

	_, err = fwb.WriteString("}")
	if err != nil {
		return err
	}

	return nil
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func MarshalCustomFieldArrayJSON(vs []Value, fwb *bufio.Writer, userProfile *UserProfile) (err error) {

	_, err = fwb.WriteString("{")
	if err != nil {
		return
	}

	if userProfile != nil {
		set := false
		for _, f := range vs {
			// Do the pruning of custom fields to only the allowed set for the user
			if !contains(userProfile.AllowedCustomFields, f.PropertyName) && !contains(userProfile.AllowedCustomFields, "*") {
				continue
			}
			if set {
				_, err = fwb.WriteString(",")
				if err != nil {
					return
				}
			}
			set = true
			_, err = fwb.WriteString("\"" + f.PropertyName + "\":\"" + strings.Replace(f.Text, "\n", "", -1) + "\"")
			if err != nil {
				return
			}
		}
	}

	_, err = fwb.WriteString("}")
	if err != nil {
		return
	}

	return nil
}
func (d *FlightState) WriteJSON(fwb *bufio.Writer, userProfile *UserProfile) error {

	_, err := fwb.WriteString("{" +
		"\"ScheduledTime\":\"" + d.ScheduledTime + "\"," +
		"\"LinkedFlight\":")
	if err != nil {
		return err
	}
	err = d.LinkedFlight.WriteJSON(fwb)
	if err != nil {
		return err
	}
	_, err = fwb.WriteString(",\"AircraftType\":")
	if err != nil {
		return err
	}
	err = d.AircraftType.WriteJSON(fwb)
	if err != nil {
		return err
	}
	_, err = fwb.WriteString(",\"Aircraft\":")
	if err != nil {
		return err
	}
	err = d.Aircraft.WriteJSON(fwb)
	if err != nil {
		return err
	}
	_, err = fwb.WriteString(",\"Route\":")
	if err != nil {
		return err
	}
	err = d.Route.WriteJSON(fwb)
	if err != nil {
		return err
	}
	_, err = fwb.WriteString(",\"Values\":")
	if err != nil {
		return err
	}
	err = MarshalCustomFieldArrayJSON(d.Values, fwb, userProfile)
	if err != nil {
		return err
	}
	_, err = fwb.WriteString(",\"StandSlots\":")
	if err != nil {
		return err
	}
	err = d.StandSlots.WriteJSON(fwb)
	if err != nil {
		return err
	}
	_, err = fwb.WriteString(",\"CarouselSlots\":")
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	err = d.CarouselSlots.WriteJSON(fwb)
	if err != nil {
		return err
	}
	_, err = fwb.WriteString(",\"GateSlots\":")
	if err != nil {
		return err
	}
	err = d.GateSlots.WriteJSON(fwb)
	if err != nil {
		return err
	}
	_, err = fwb.WriteString(",\"CheckInSlots\":")
	if err != nil {
		return err
	}
	err = d.CheckInSlots.WriteJSON(fwb)
	if err != nil {
		return err
	}
	_, err = fwb.WriteString(",\"ChuteSlots\":")
	if err != nil {
		return err
	}
	err = d.ChuteSlots.WriteJSON(fwb)
	if err != nil {
		return err
	}
	err = fwb.Flush()
	if err != nil {
		return err
	}

	return nil
}

type Change struct {
	PropertyName string `xml:"propertyName,attr"`
	OldValue     string `xml:"OldValue"`
	NewValue     string `xml:"NewValue"`
}

type GateSlotsChange struct {
	OldValue struct {
		GateSlot []struct {
			Value []struct {
				Text         string `xml:",chardata"`
				PropertyName string `xml:"propertyName,attr"`
			} `xml:"Value"`
			Gate struct {
				Value []struct {
					Text         string `xml:",chardata"`
					PropertyName string `xml:"propertyName,attr"`
				} `xml:"Value"`
				Area struct {
					Value []struct {
						Text         string `xml:",chardata"`
						PropertyName string `xml:"propertyName,attr"`
					} `xml:"Value"`
				} `xml:"Area"`
			} `xml:"Gate"`
		} `xml:"GateSlot"`
	} `xml:"OldValue"`
	NewValue struct {
		GateSlot []struct {
			Value []struct {
				Text         string `xml:",chardata"`
				PropertyName string `xml:"propertyName,attr"`
			} `xml:"Value"`
			Gate struct {
				Value []struct {
					Text         string `xml:",chardata"`
					PropertyName string `xml:"propertyName,attr"`
				} `xml:"Value"`
				Area struct {
					Value []struct {
						Text         string `xml:",chardata"`
						PropertyName string `xml:"propertyName,attr"`
					} `xml:"Value"`
				} `xml:"Area"`
			} `xml:"Gate"`
		} `xml:"GateSlot"`
	} `xml:"NewValue"`
}

func (sc *GateSlotsChange) MarshalJSON() ([]byte, error) {

	fwb := strings.Builder{}

	fwb.WriteString("{")
	fwb.WriteString("\"OldValue\":[")

	for idx, s := range sc.OldValue.GateSlot {
		if idx > 0 {
			fwb.WriteString(",")
		}
		fwb.WriteString("{")
		for idx3, v := range s.Value {
			if idx3 > 0 {
				fwb.WriteString(",")
			}
			fwb.WriteString("\"" + v.PropertyName + "\":\"" + v.Text + "\"")
		}
		for _, v := range s.Gate.Value {
			fwb.WriteString(",")
			fwb.WriteString("\"" + v.PropertyName + "\":\"" + v.Text + "\"")
		}

		for _, v := range s.Gate.Area.Value {
			fwb.WriteString(",")
			fwb.WriteString("\"Area" + v.PropertyName + "\":\"" + v.Text + "\"")
		}
		fwb.WriteString("}")
	}
	fwb.WriteString("],")

	fwb.WriteString("\"NewValue\":[")

	for idx, s := range sc.NewValue.GateSlot {
		if idx > 0 {
			fwb.WriteString(",")
		}
		fwb.WriteString("{")
		for idx3, v := range s.Value {
			if idx3 > 0 {
				fwb.WriteString(",")
			}
			fwb.WriteString("\"" + v.PropertyName + "\":\"" + v.Text + "\"")
		}
		for _, v := range s.Gate.Value {
			fwb.WriteString(",")
			fwb.WriteString("\"" + v.PropertyName + "\":\"" + v.Text + "\"")
		}

		for _, v := range s.Gate.Area.Value {
			fwb.WriteString(",")
			fwb.WriteString("\"Area" + v.PropertyName + "\":\"" + v.Text + "\"")
		}
		fwb.WriteString("}")
	}
	fwb.WriteString("]")

	fwb.WriteString("}")
	var sendText = fwb.String()

	return []byte(sendText), nil

}

type StandSlotsChange struct {
	OldValue struct {
		StandSlot []struct {
			Value []struct {
				Text         string `xml:",chardata"`
				PropertyName string `xml:"propertyName,attr"`
			} `xml:"Value"`
			Stand struct {
				Value []struct {
					Text         string `xml:",chardata"`
					PropertyName string `xml:"propertyName,attr"`
				} `xml:"Value"`
				Area struct {
					Value []struct {
						Text         string `xml:",chardata"`
						PropertyName string `xml:"propertyName,attr"`
					} `xml:"Value"`
				} `xml:"Area"`
			} `xml:"Stand"`
		} `xml:"StandSlot"`
	} `xml:"OldValue"`
	NewValue struct {
		StandSlot []struct {
			Value []struct {
				Text         string `xml:",chardata"`
				PropertyName string `xml:"propertyName,attr"`
			} `xml:"Value"`
			Stand struct {
				Value []struct {
					Text         string `xml:",chardata"`
					PropertyName string `xml:"propertyName,attr"`
				} `xml:"Value"`
				Area struct {
					Value []struct {
						Text         string `xml:",chardata"`
						PropertyName string `xml:"propertyName,attr"`
					} `xml:"Value"`
				} `xml:"Area"`
			} `xml:"Stand"`
		} `xml:"StandSlot"`
	} `xml:"NewValue"`
}

func (sc *StandSlotsChange) MarshalJSON() ([]byte, error) {

	fwb := strings.Builder{}

	fwb.WriteString("{")
	fwb.WriteString("\"OldValue\":[")

	for idx, s := range sc.OldValue.StandSlot {
		if idx > 0 {
			fwb.WriteString(",")
		}
		fwb.WriteString("{")
		for idx3, v := range s.Value {
			if idx3 > 0 {
				fwb.WriteString(",")
			}
			fwb.WriteString("\"" + v.PropertyName + "\":\"" + v.Text + "\"")
		}
		for _, v := range s.Stand.Value {
			fwb.WriteString(",")
			fwb.WriteString("\"" + v.PropertyName + "\":\"" + v.Text + "\"")
		}

		for _, v := range s.Stand.Area.Value {
			fwb.WriteString(",")
			fwb.WriteString("\"Area" + v.PropertyName + "\":\"" + v.Text + "\"")
		}
		fwb.WriteString("}")
	}
	fwb.WriteString("],")

	fwb.WriteString("\"NewValue\":[")

	for idx, s := range sc.NewValue.StandSlot {
		if idx > 0 {
			fwb.WriteString(",")
		}
		fwb.WriteString("{")
		for idx3, v := range s.Value {
			if idx3 > 0 {
				fwb.WriteString(",")
			}
			fwb.WriteString("\"" + v.PropertyName + "\":\"" + v.Text + "\"")
		}
		for _, v := range s.Stand.Value {
			fwb.WriteString(",")
			fwb.WriteString("\"" + v.PropertyName + "\":\"" + v.Text + "\"")
		}

		for _, v := range s.Stand.Area.Value {
			fwb.WriteString(",")
			fwb.WriteString("\"Area" + v.PropertyName + "\":\"" + v.Text + "\"")
		}
		fwb.WriteString("}")
	}
	fwb.WriteString("]")

	fwb.WriteString("}")
	var sendText = fwb.String()

	return []byte(sendText), nil

}

type CarouselSlotsChange struct {
	OldValue struct {
		CarouselSlot []struct {
			Value    []PropertyValuePair `xml:"Value"`
			Carousel struct {
				Value []PropertyValuePair `xml:"Value"`
				Area  struct {
					Value []struct {
						Text         string `xml:",chardata"`
						PropertyName string `xml:"propertyName,attr"`
					} `xml:"Value"`
				} `xml:"Area"`
			} `xml:"Carousel"`
		} `xml:"CarouselSlot"`
	} `xml:"OldValue"`
	NewValue struct {
		CarouselSlot []struct {
			Value    []PropertyValuePair `xml:"Value"`
			Carousel struct {
				Value []PropertyValuePair `xml:"Value"`
				Area  struct {
					Value []struct {
						Text         string `xml:",chardata"`
						PropertyName string `xml:"propertyName,attr"`
					} `xml:"Value"`
				} `xml:"Area"`
			} `xml:"Carousel"`
		} `xml:"CarouselSlot"`
	} `xml:"NewValue"`
}

func (sc *CarouselSlotsChange) MarshalJSON() ([]byte, error) {

	fwb := strings.Builder{}

	fwb.WriteString("{")
	fwb.WriteString("\"OldValue\":[")

	for idx, s := range sc.OldValue.CarouselSlot {
		if idx > 0 {
			fwb.WriteString(",")
		}
		fwb.WriteString("{")
		for idx3, v := range s.Value {
			if idx3 > 0 {
				fwb.WriteString(",")
			}
			fwb.WriteString("\"" + v.PropertyName + "\":\"" + v.Text + "\"")
		}
		for _, v := range s.Carousel.Value {
			fwb.WriteString(",")
			fwb.WriteString("\"" + v.PropertyName + "\":\"" + v.Text + "\"")
		}

		for _, v := range s.Carousel.Area.Value {
			fwb.WriteString(",")
			fwb.WriteString("\"Area" + v.PropertyName + "\":\"" + v.Text + "\"")
		}
		fwb.WriteString("}")
	}
	fwb.WriteString("],")

	fwb.WriteString("\"NewValue\":[")

	for idx, s := range sc.NewValue.CarouselSlot {
		if idx > 0 {
			fwb.WriteString(",")
		}
		fwb.WriteString("{")
		for idx3, v := range s.Value {
			if idx3 > 0 {
				fwb.WriteString(",")
			}
			fwb.WriteString("\"" + v.PropertyName + "\":\"" + v.Text + "\"")
		}
		for _, v := range s.Carousel.Value {
			fwb.WriteString(",")
			fwb.WriteString("\"" + v.PropertyName + "\":\"" + v.Text + "\"")
		}

		for _, v := range s.Carousel.Area.Value {
			fwb.WriteString(",")
			fwb.WriteString("\"Area" + v.PropertyName + "\":\"" + v.Text + "\"")
		}
		fwb.WriteString("}")
	}
	fwb.WriteString("]")

	fwb.WriteString("}")
	var sendText = fwb.String()

	return []byte(sendText), nil

}

type ChuteSlotsChange struct {
	OldValue struct {
		ChuteSlot struct {
			Value PropertyValuePair `xml:"Value"`
			Chute struct {
				Value []PropertyValuePair `xml:"Value"`
				Area  struct {
					Value PropertyValuePair `xml:"Value"`
				} `xml:"Area"`
			} `xml:"Chute"`
		} `xml:"ChuteSlot"`
	} `xml:"OldValue"`
	NewValue struct {
		ChuteSlot struct {
			Value PropertyValuePair `xml:"Value"`
			Chute struct {
				Value []PropertyValuePair `xml:"Value"`
				Area  struct {
					Value PropertyValuePair `xml:"Value"`
				} `xml:"Area"`
			} `xml:"Chute"`
		} `xml:"ChuteSlot"`
	} `xml:"NewValue"`
}
type AircraftTypeChange struct {
	OldValue struct {
		AircraftType struct {
			AircraftTypeId struct {
				AircraftTypeCode []struct {
					Text        string `xml:",chardata"`
					CodeContext string `xml:"codeContext,attr"`
				} `xml:"AircraftTypeCode"`
			} `xml:"AircraftTypeId"`
			Value PropertyValuePair `xml:"Value"`
		} `xml:"AircraftType"`
	} `xml:"OldValue"`
	NewValue struct {
		AircraftType struct {
			AircraftTypeId struct {
				AircraftTypeCode []struct {
					Text        string `xml:",chardata"`
					CodeContext string `xml:"codeContext,attr"`
				} `xml:"AircraftTypeCode"`
			} `xml:"AircraftTypeId"`
			Value PropertyValuePair `xml:"Value"`
		} `xml:"AircraftType"`
	} `xml:"NewValue"`
}
type AircraftChange struct {
	OLdValue struct {
		Aircraft struct {
			AircraftId struct {
				Registration string `xml:"Registration"`
			} `xml:"AircraftId"`
			Value PropertyValuePair `xml:"Value"`
		} `xml:"Aircraft"`
	} `xml:"OldValue"`
	NewValue struct {
		Aircraft struct {
			AircraftId struct {
				Registration string `xml:"Registration"`
			} `xml:"AircraftId"`
			Value PropertyValuePair `xml:"Value"`
		} `xml:"Aircraft"`
	} `xml:"NewValue"`
}
type RouteChange struct {
	OldValue struct {
		Route struct {
			CustomsType string    `xml:"customsType,attr"`
			ViaPoints   ViaPoints `xml:"ViaPoints"`
		} `xml:"Route"`
	} `xml:"OldValue"`
	NewValue struct {
		Route struct {
			CustomsType string    `xml:"customsType,attr"`
			ViaPoints   ViaPoints `xml:"ViaPoints"`
		} `xml:"Route"`
	} `xml:"NewValue"`
}

type CheckInSlotsChange struct {
	OldValue struct {
		CheckInSlot []CheckInSlot `xml:"CheckInSlot" json:"CheckInSlots"`
	} `xml:"OldValue"`
	NewValue struct {
		CheckInSlot []CheckInSlot `xml:"CheckInSlot" json:"CheckInSlots"`
	} `xml:"NewValue"`
}

type LinkedFlightChange struct {
	OldValue struct {
		LinkedFlight struct {
			FlightId FlightId `xml:"FlightId"`
			Value    []struct {
				Text         string `xml:",chardata"`
				PropertyName string `xml:"propertyName,attr"`
			} `xml:"Value"`
		} `xml:"LinkedFlight"`
	} `xml:"OldValue"`
	NewValue struct {
		LinkedFlight struct {
			FlightId FlightId `xml:"FlightId"`
			Value    []struct {
				Text         string `xml:",chardata"`
				PropertyName string `xml:"propertyName,attr"`
			} `xml:"Value"`
		} `xml:"LinkedFlight"`
	} `xml:"NewValue"`
}
type FlightChanges struct {
	AircraftTypeChange  *AircraftTypeChange  `xml:"AircraftTypeChange" json:"AircraftTypeChange"`
	AircraftChange      *AircraftChange      `xml:"AircraftChange" json:"AircraftChange"`
	CarouselSlotsChange *CarouselSlotsChange `xml:"CarouselSlotsChange" json:"CarouselSlotsChange"`
	GateSlotsChange     *GateSlotsChange     `xml:"GateSlotsChange" json:"GateSlotsChange"`
	StandSlotsChange    *StandSlotsChange    `xml:"StandSlotsChange" json:"StandSlotsChange"`
	ChuteSlotsChange    *ChuteSlotsChange    `xml:"ChuteSlotsChange" json:"ChuteSlotsChange"`
	CheckinSlotsChange  *CheckInSlotsChange  `xml:"CheckInSlotsChange" json:"CheckInSlotsChange"`
	RouteChange         *RouteChange         `xml:"RouteChange" json:"RouteChange"`
	LinkedFlightChange  *LinkedFlightChange  `xml:"LinkedFlightChange" json:"LinkedFlightChange"`
	Changes             []Change             `xml:"Change"  json:"-" `
}
type Flight struct {
	PrevNode      *Flight       `xml:"-" json:"-"`
	NextNode      *Flight       `xml:"-" json:"-"`
	Action        string        `xml:"Action" json:"Action"`
	FlightId      FlightId      `xml:"FlightId" json:"FlightId"`
	FlightState   FlightState   `xml:"FlightState" json:"FlightState"`
	FlightChanges FlightChanges `xml:"FlightChanges" json:"FlightChanges"`
	LastUpdate    time.Time     `xml:"LastUpdate" json:"LastUpdate"`
}

func (d *Flight) GetGobParameters() (id, airport, airline, fltnum, kind, route, sto string, stoUnix int) {

	id = d.GetFlightID()
	kind = d.GetFlightDirection()
	airline = d.GetIATAAirline()
	route = d.GetFlightRoute()
	sto = d.GetSTO().String()
	stoUnix = int(d.GetSTO().Unix())
	fltnum = d.FlightId.FlightNumber
	airport = d.GetIATAAirport()
	return
}

func (d *Flight) WriteJSON(fwb *bufio.Writer, userProfile *UserProfile, statusOnly bool) error {

	_, err := fwb.WriteString("{")
	if err != nil {
		return err
	}

	if d.Action != "" {
		_, err = fwb.WriteString(`"Action":"` + d.Action + `",`)
		if err != nil {
			return err
		}
	}
	err = d.FlightId.WriteJSON(fwb)
	if err != nil {
		return err
	}
	_, err = fwb.WriteString(`,"FlightState":`)
	if err != nil {
		return err
	}
	err = d.FlightState.WriteJSON(fwb, userProfile)
	if err != nil {
		return err
	}

	if statusOnly {
		_, err = fwb.WriteString("}}")
		if err != nil {
			return err
		}
	} else {
		// Using the built in JSON serialiser for the Changes because I'm too lazy to write a custom serilaizer
		flightChanges, _ := json.Marshal(d.FlightChanges)
		_, err = fwb.WriteString("},\n\"Changes\":")
		if err != nil {
			return err
		}
		_, err = fwb.Write(flightChanges)
		if err != nil {
			return err
		}
		_, err = fwb.WriteString(",\n\"ValueChanges\":[")
		if err != nil {
			return err
		}
		f := false
		for idx, c := range d.FlightChanges.Changes {
			if contains(userProfile.AllowedCustomFields, c.PropertyName) || contains(userProfile.AllowedCustomFields, "*") {
				if idx > 0 && f {
					_, err = fwb.WriteString(",")
					if err != nil {
						return err
					}
				}
				f = true
				_, err = fwb.WriteString("{\"PropertyName\":\"" + c.PropertyName + "\", \"OldValue\":\"" + c.OldValue + "\",\"NewValue\":\"" + c.NewValue + "\"}")
				if err != nil {
					return err
				}
			}
		}

		_, err = fwb.WriteString("]}")
		if err != nil {
			return err
		}
	}

	return nil

}

type Flights struct {
	Flight []Flight `xml:"Flight" json:"Flights"`
}

type StandAllocation struct {
	Stand  Stand
	From   time.Time
	To     time.Time
	Flight FlightId
}

type StandAllocations struct {
	Allocations []StandAllocation
}

// Resource definitions

type Area struct {
	Value []Value `xml:"Value"`
}

type Stand struct {
	Value []Value `xml:"Value" json:"Slot,omitempty"`
	Area  Area    `xml:"Area" json:"Area,omitempty"`
}

type StandSlot struct {
	Value []Value `xml:"Value" json:"Slot,omitempty"`
	Stand Stand   `xml:"Stand" json:"Area,omitempty"`
}
type StandSlots struct {
	StandSlot []StandSlot `xml:"StandSlot" json:"StandSlot,omitempty"`
}
type Carousel struct {
	Value []Value `xml:"Value" json:"Slot,omitempty"`
	Area  Area    `xml:"Area" json:"Area,omitempty"`
}
type CarouselSlot struct {
	Value    []Value  `xml:"Value" json:"Slot,omitempty"`
	Carousel Carousel `xml:"Carousel" json:"Carousel,omitempty"`
}
type CarouselSlots struct {
	CarouselSlot []CarouselSlot `xml:"CarouselSlot" json:"CarouselSlot,omitempty"`
}

type Gate struct {
	Value []Value `xml:"Value"`
	Area  Area    `xml:"Area"`
}

type GateSlot struct {
	Value []Value `xml:"Value"`
	Gate  Gate    `xml:"Gate"`
}
type GateSlots struct {
	GateSlot []GateSlot `xml:"GateSlot" json:"GateSlot,omitempty"`
}
type CheckIn struct {
	Value []Value `xml:"Value"`
	Area  Area    `xml:"Area"`
}
type CheckInSlot struct {
	Value   []Value `xml:"Value"`
	CheckIn CheckIn `xml:"CheckIn"`
}
type CheckInSlots struct {
	CheckInSlot []CheckInSlot `xml:"CheckInSlot" json:"CheckInSlot,omitempty"`
}
type Chute struct {
	Value []Value `xml:"Value"`
	Area  Area    `xml:"Area"`
}
type ChuteSlot struct {
	Value []Value `xml:"Values"`
	Chute Chute   `xml:"Chute"`
}
type ChuteSlots struct {
	ChuteSlot []ChuteSlot `xml:"ChuteSlot" json:"ChuteSlot,omitempty"`
}

type Envelope struct {
	Body struct {
		GetFlightsResponse struct {
			GetFlightsResult struct {
				WebServiceResult struct {
					ApiResponse struct {
						Data struct {
							Flights Flights `xml:"Flights"`
						} `xml:"Data"`
					} `xml:"ApiResponse"`
				} `xml:"WebServiceResult"`
			} `xml:"GetFlightsResult"`
		} `xml:"GetFlightsResponse"`
	} `xml:"Body"`
}

type FlightCreatedNotificationEnvelope struct {
	Content struct {
		FlightCreatedNotification struct {
			Flight Flight `xml:"Flight"`
		} `xml:"FlightCreatedNotification"`
	} `xml:"Content"`
}
type FlightUpdatedNotificationEnvelope struct {
	Content struct {
		FlightUpdatedNotification struct {
			Flight Flight `xml:"Flight"`
		} `xml:"FlightUpdatedNotification"`
	} `xml:"Content"`
}
type FlightDeletedNotificationEnvelope struct {
	Content struct {
		FlightDeletedNotification struct {
			Flight Flight `xml:"Flight"`
		} `xml:"FlightDeletedNotification"`
	} `xml:"Content"`
}

func (f Flight) GetSDO() time.Time {

	sdo := f.FlightId.ScheduledDate
	sdod, _ := time.Parse("2006-01-02", sdo)
	return sdod
}
func (f Flight) GetProperty(property string) string {
	for _, v := range f.FlightState.Values {
		if v.PropertyName == property {
			return v.Text
		}
	}
	return ""
}
func (f Flight) IsArrival() bool {
	if f.FlightId.FlightKind == "Arrival" {
		return true
	} else {
		return false
	}
}
func (f Flight) GetIATAAirline() string {
	for _, v := range f.FlightId.AirlineDesignator {
		if v.CodeContext == "IATA" {
			return v.Text
		}
	}
	return ""
}
func (f Flight) GetIATAAirport() string {
	for _, v := range f.FlightId.AirportCode {
		if v.CodeContext == "IATA" {
			return v.Text
		}
	}
	return ""
}
func (f Flight) GetICAOAirline() string {
	for _, v := range f.FlightId.AirlineDesignator {
		if v.CodeContext == "ICAO" {
			return v.Text
		}
	}
	return ""
}
func (f Flight) GetFlightID() string {

	airline := f.GetIATAAirline()
	fltNum := f.FlightId.FlightNumber
	sto := f.FlightState.ScheduledTime
	// kind := "D"
	// if f.IsArrival() {
	// 	kind = "A"
	// }
	return airline + fltNum + "@" + sto
}
func (f Flight) GetFlightDirection() string {

	if f.IsArrival() {
		return "Arrival"
	} else {
		return "Departure"
	}
}
func (f Flight) GetFlightRoute() string {

	var sb strings.Builder
	idx := 0

	for _, rp := range f.FlightState.Route.ViaPoints.RouteViaPoint {
		for _, ap := range rp.AirportCode {
			if idx > 0 && ap.CodeContext == "IATA" {
				sb.WriteString(",")
			}

			if ap.CodeContext == "IATA" {
				sb.WriteString(ap.Text)
				idx++
			}

		}
	}

	return sb.String()
}
func (f Flight) GetAircraftType() string {

	sb := "-"

	for _, rp := range f.FlightState.AircraftType.AircraftTypeId.AircraftTypeCode {

		if rp.CodeContext == "IATA" {
			sb = rp.Text
		}
	}

	return sb
}
func (f Flight) GetAircraftRegistration() string {

	if f.FlightState.Aircraft.AircraftId.Registration != "" {
		return f.FlightState.Aircraft.AircraftId.Registration
	} else {
		return "-"
	}
}
func (f Flight) GetSTO() time.Time {

	sto := f.FlightState.ScheduledTime

	if sto != "" {
		stot, err := time.ParseInLocation("2006-01-02T15:04:05", sto, timeservice.Loc)
		if err == nil {
			return stot
		}
		return time.Now()
	}

	return time.Now()
}

func (p CheckInSlot) GetResourceID() (name string, from time.Time, to time.Time) {

	for _, v := range p.Value {

		if v.PropertyName == "StartTime" {
			from, _ = time.ParseInLocation(timeservice.Layout, v.Text, timeservice.Loc)
			continue
		}
		if v.PropertyName == "EndTime" {
			to, _ = time.ParseInLocation(timeservice.Layout, v.Text, timeservice.Loc)
			continue
		}
	}

	for _, v := range p.CheckIn.Value {
		if v.PropertyName == "Name" {
			name = v.Text
			continue
		}
	}
	return
}

func (p StandSlot) GetResourceID() (name string, from time.Time, to time.Time) {

	for _, v := range p.Value {

		if v.PropertyName == "StartTime" {
			from, _ = time.ParseInLocation(timeservice.Layout, v.Text, timeservice.Loc)
			continue
		}
		if v.PropertyName == "EndTime" {
			to, _ = time.ParseInLocation(timeservice.Layout, v.Text, timeservice.Loc)
			continue
		}
	}

	for _, v := range p.Stand.Value {
		if v.PropertyName == "Name" {
			name = v.Text
			continue
		}
	}
	return
}
func (p CarouselSlot) GetResourceID() (name string, from time.Time, to time.Time) {

	for _, v := range p.Value {

		if v.PropertyName == "StartTime" {
			from, _ = time.ParseInLocation(timeservice.Layout, v.Text, timeservice.Loc)
			continue
		}
		if v.PropertyName == "EndTime" {
			to, _ = time.ParseInLocation(timeservice.Layout, v.Text, timeservice.Loc)
			continue
		}
	}

	for _, v := range p.Carousel.Value {
		if v.PropertyName == "Name" {
			name = v.Text
			continue
		}
	}
	return
}

func (p ChuteSlot) GetResourceID() (name string, from time.Time, to time.Time) {

	for _, v := range p.Value {

		if v.PropertyName == "StartTime" {
			from, _ = time.ParseInLocation(timeservice.Layout, v.Text, timeservice.Loc)
			continue
		}
		if v.PropertyName == "EndTime" {
			to, _ = time.ParseInLocation(timeservice.Layout, v.Text, timeservice.Loc)
			continue
		}
	}

	for _, v := range p.Chute.Value {
		if v.PropertyName == "Name" {
			name = v.Text
			continue
		}
	}
	return
}

func (p GateSlot) GetResourceID() (name string, from time.Time, to time.Time) {

	for _, v := range p.Value {

		if v.PropertyName == "StartTime" {
			from, _ = time.ParseInLocation(timeservice.Layout, v.Text, timeservice.Loc)
			continue
		}
		if v.PropertyName == "EndTime" {
			to, _ = time.ParseInLocation(timeservice.Layout, v.Text, timeservice.Loc)
			continue
		}
	}

	for _, v := range p.Gate.Value {
		if v.PropertyName == "Name" {
			name = v.Text
			continue
		}
	}
	return name, from, to
}

func (r Repository) MinimumProperties(min int) {

	fmt.Printf("Setting the number of Custom Fields in sample flights to %v", min)
	if len(r.FlightLinkedList.Head.FlightState.Values) < min {

		currentNode := r.FlightLinkedList.Head

		for currentNode != nil {
			i := len(currentNode.FlightState.Values)
			for len(currentNode.FlightState.Values) <= min {
				prop := Value{
					PropertyName: "Custom_Field_Name_%" + strconv.Itoa(i),
					Text:         "Custom_Field_Value_%" + strconv.Itoa(i),
				}
				currentNode.FlightState.Values = append(currentNode.FlightState.Values, prop)
				i++
			}
			currentNode = currentNode.NextNode
		}
	} else if min < len(r.FlightLinkedList.Head.FlightState.Values) {
		currentNode := r.FlightLinkedList.Head

		for currentNode != nil {
			currentNode.FlightState.Values = currentNode.FlightState.Values[:min]
			currentNode = currentNode.NextNode
		}
	}

	fmt.Printf(" - Completed\n")

}
