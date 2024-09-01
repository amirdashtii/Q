package provider

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/amirdashtii/Q/flight-ticket-service/models"
)

var flightProviderHost string

const (
	flightProviderEndpoint = "/flights"
	httpTimeout            = 5 * time.Second
)

type FlightsResponse struct {
	Flights []models.Flight `json:"flights"`
}

type FlightResponse struct {
	Flight models.Flight `json:"flight"`
}

type FlightProviderClient struct {
	client *http.Client
}

func NewFlightProviderClient() *FlightProviderClient {
	flightProviderHost = "http://" + os.Getenv("FLIGHT_PROVIDER_HOST") + ":" + os.Getenv("FLIGHT_PROVIDER_PORT")
	tr := &http.Transport{}
	cl := &http.Client{
		Transport: tr,
		Timeout:   httpTimeout,
	}

	return &FlightProviderClient{
		client: cl,
	}
}

func (pc *FlightProviderClient) RequestFlights(flightReq *models.FlightSearchRequest, flights *[]models.Flight) error {

	u, err := url.Parse(flightProviderHost + flightProviderEndpoint)
	if err != nil {
		return err
	}

	query := u.Query()
	query.Set("source", flightReq.Source)
	query.Set("destination", flightReq.Destination)
	query.Set("departure_date", flightReq.DepartureDate)
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := pc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response FlightsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}

	*flights = response.Flights

	return nil
}

func (pc *FlightProviderClient) RequestFlight(id *string, flight *models.Flight) error {
	req, err := http.NewRequest(http.MethodGet, flightProviderHost+flightProviderEndpoint+"/"+*id, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := pc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var response FlightResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}

	*flight = response.Flight

	return nil
}
