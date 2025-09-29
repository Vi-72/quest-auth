package contracts

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/Vi-72/quest-auth/internal/core/ports"
)

type EventPublisherContractSuite struct {
	suite.Suite
	publisher ports.EventPublisher
	ctx       context.Context
}

// simple event impl
type testEvent struct{ id uuid.UUID }

func (e testEvent) GetID() uuid.UUID          { return e.id }
func (e testEvent) GetName() string           { return "test" }
func (e testEvent) GetAggregateID() uuid.UUID { return e.id }
func (e testEvent) At() time.Time             { return time.Now() }

func TestNullEventPublisherContract(t *testing.T) {
	s := &EventPublisherContractSuite{publisher: &ports.NullEventPublisher{}, ctx: context.Background()}
	suite.Run(t, s)
}

func (s *EventPublisherContractSuite) TestPublishSingleEvent() {
	ev := testEvent{id: uuid.New()}
	err := s.publisher.Publish(s.ctx, ev)
	s.Assert().NoError(err)
}

func (s *EventPublisherContractSuite) TestPublishMultipleEvents() {
	e1 := testEvent{id: uuid.New()}
	e2 := testEvent{id: uuid.New()}
	err := s.publisher.Publish(s.ctx, e1, e2)
	s.Assert().NoError(err)
}

func (s *EventPublisherContractSuite) TestPublishEmpty() {
	err := s.publisher.Publish(s.ctx)
	s.Assert().NoError(err)
}

func (s *EventPublisherContractSuite) TestPublishAsync() {
	e := testEvent{id: uuid.New()}
	s.publisher.PublishAsync(s.ctx, e)
}
