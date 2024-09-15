package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/amirdashtii/Q/flight-ticket-service/models"
)

var (
	flightProviderHost     string
	flightProviderEndpoint = "/flights"
	httpTimeout            = 5 * time.Second
)

type FlightsResponse struct {
	Flights []models.ProviderFlight `json:"flights"`
}

type FlightResponse struct {
	Flight models.ProviderFlight `json:"flight"`
}

type ProviderClient struct {
	client *http.Client
}

func NewProviderClient() *ProviderClient {
	flightProviderHost = os.Getenv("FLIGHT_PROVIDER_HOST")

	tr := &http.Transport{}
	cl := &http.Client{
		Transport: tr,
		Timeout:   httpTimeout,
	}

	return &ProviderClient{
		client: cl,
	}
}

func (pc *ProviderClient) RequestFlights(flightReq *models.FlightSearchRequest, flights *[]models.ProviderFlight) error {

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

func (pc *ProviderClient) RequestFlight(id *string, flight *models.ProviderFlight) error {
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
	err = json.NewDecoder(resp.Body).Decode(flight)
	if err != nil {
		return err
	}

	*flight = response.Flight

	return nil
}

func (pc *ProviderClient) ReserveTicketWithProvider(seats int, flightID string) error {

	data := map[string]int{
		"seats": seats,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, flightProviderHost+flightProviderEndpoint+"/"+flightID+"/reserve", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := pc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil
}

func (pc *ProviderClient) CancelTicketWithProvider(seets int, flightID string) error {

	data := map[string]int{
		"seats": seets,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, flightProviderHost+flightProviderEndpoint+"/"+flightID+"/cancel", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := pc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil
}
