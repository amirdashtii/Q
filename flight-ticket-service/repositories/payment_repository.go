package repositories

import "github.com/amirdashtii/Q/flight-ticket-service/models"

func (p *Postgres) CreatePaymentReceipt(receivedPaymentRequest *models.PaymentReceipt) error {
	result := p.db.Create(receivedPaymentRequest)
	return result.Error
}
