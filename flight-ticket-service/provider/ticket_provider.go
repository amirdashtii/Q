package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/amirdashtii/Q/flight-ticket-service/models"
)

type TicketProviderClient struct {
	client *http.Client
}

func NewTicketProviderClient() *TicketProviderClient {
	flightProviderHost = "http://" + os.Getenv("FLIGHT_PROVIDER_HOST") + ":" + os.Getenv("FLIGHT_PROVIDER_PORT")
	tr := &http.Transport{}
	cl := &http.Client{
		Transport: tr,
		Timeout:   httpTimeout,
	}

	return &TicketProviderClient{
		client: cl,
	}
}

func (pc *TicketProviderClient) ReserveTicketWithProvider(reservation *models.TicketReservation) error {

	data := map[string]int{
		"seats": len(reservation.Passengers),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, flightProviderHost+flightProviderEndpoint+"/"+reservation.FlightID.String()+"/reserve", bytes.NewBuffer(jsonData))
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
func (pc *TicketProviderClient) CancelTicketWithProvider(reservation *models.TicketReservation) error {

	data := map[string]int{
		"seats": len(reservation.Passengers),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, flightProviderHost+flightProviderEndpoint+"/"+reservation.FlightID.String()+"/cancel", bytes.NewBuffer(jsonData))
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
