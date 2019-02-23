package godatabus

import (
	"testing"

	"github.com/streadway/amqp"
)

func TestHandleRejectWithNoRejectsRequeues(t *testing.T) {
	ackMock := NewAcknowledgerMock()
	delivery := amqp.Delivery{Headers: make(map[string]interface{})}
	delivery.Acknowledger = ackMock

	handleRejectWithRetries(&delivery)

	if *ackMock.CallsRejectWithRequeue != 1 || *ackMock.CallsRejectWithoutRequeue != 0 {
		t.Errorf("Delivery has been rejected without requeue %d, %d", ackMock.CallsRejectWithRequeue, ackMock.CallsRejectWithoutRequeue)
	}
	if delivery.Headers["retry-attempts"].(int) != 1 {
		t.Error("Retry attempts header was not set")
	}
}

func TestHandleRejectWithZeroRejectsRequeues(t *testing.T) {
	ackMock := NewAcknowledgerMock()
	delivery := amqp.Delivery{Headers: make(map[string]interface{})}
	delivery.Headers["retry-attempts"] = 0
	delivery.Acknowledger = ackMock

	handleRejectWithRetries(&delivery)

	if *ackMock.CallsRejectWithRequeue != 1 || *ackMock.CallsRejectWithoutRequeue != 0 {
		t.Error("Delivery has been rejected without requeue")
	}
	if delivery.Headers["retry-attempts"].(int) != 1 {
		t.Error("Retry attempts header was not incremented")
	}
}

func TestHandleRejectWithTwoRejectsRequeues(t *testing.T) {
	ackMock := NewAcknowledgerMock()
	delivery := amqp.Delivery{Headers: make(map[string]interface{})}
	delivery.Headers["retry-attempts"] = 2
	delivery.Acknowledger = ackMock

	handleRejectWithRetries(&delivery)

	if *ackMock.CallsRejectWithRequeue != 1 || *ackMock.CallsRejectWithoutRequeue != 0 {
		t.Error("Delivery has been rejected without requeue")
	}
	if delivery.Headers["retry-attempts"].(int) != 3 {
		t.Error("Retry attempts header was not incremented")
	}
}

func TestHandleRejectWithThreeRejectsDoesntRequeue(t *testing.T) {
	ackMock := NewAcknowledgerMock()
	delivery := amqp.Delivery{Headers: make(map[string]interface{})}
	delivery.Headers["retry-attempts"] = 3
	delivery.Acknowledger = ackMock

	handleRejectWithRetries(&delivery)

	if *ackMock.CallsRejectWithRequeue != 0 || *ackMock.CallsRejectWithoutRequeue != 1 {
		t.Error("Delivery has been rejected with requeue")
	}
	if delivery.Headers["retry-attempts"].(int) != 4 {
		t.Error("Retry attempts header was not incremented")
	}
}

type AcknowledgerMock struct {
	CallsRejectWithRequeue    *int
	CallsRejectWithoutRequeue *int
}

func NewAcknowledgerMock() AcknowledgerMock {
	return AcknowledgerMock{CallsRejectWithoutRequeue: new(int), CallsRejectWithRequeue: new(int)}
}

func (ack AcknowledgerMock) Ack(tag uint64, multiple bool) error {
	return nil
}

func (ack AcknowledgerMock) Nack(tag uint64, multiple bool, requeue bool) error {
	return nil
}

func (ack AcknowledgerMock) Reject(tag uint64, requeue bool) error {
	if requeue {
		*ack.CallsRejectWithRequeue++
	} else {
		*ack.CallsRejectWithoutRequeue++
	}
	return nil
}
