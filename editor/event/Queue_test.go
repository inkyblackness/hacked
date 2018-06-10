package event_test

import (
	"testing"

	"github.com/inkyblackness/hacked/editor/event"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type QueueSuite struct {
	suite.Suite

	queue event.Queue

	events       []event.Event
	pendingEvent event.Event
}

func TestQueueSuite(t *testing.T) {
	suite.Run(t, new(QueueSuite))
}

func (suite *QueueSuite) SetupTest() {
	suite.queue = event.Queue{}
	suite.events = nil
}

func (suite *QueueSuite) TestNewQueueIsEmpty() {
	suite.thenQueueShouldBeEmpty()
}

func (suite *QueueSuite) TestQueueIsNotEmptyWhenAddingEvent() {
	suite.whenAddingEvent(suite.someEvent())
	suite.thenQueueShouldNotBeEmpty()
}

func (suite *QueueSuite) TestNilEventsAreIgnored() {
	suite.whenAddingEvent(nil)
	suite.thenQueueShouldBeEmpty()
}

func (suite *QueueSuite) TestTypedNilEventsAreNotIgnored() {
	var typedEvent *struct{}
	suite.whenAddingEvent(typedEvent)
	suite.thenQueueShouldNotBeEmpty()
}

func (suite *QueueSuite) TestForwardToCallsListenerForEachEvent() {
	eventA := suite.keyedEvent(1)
	eventB := suite.keyedEvent(2)
	suite.givenQueuedEvents(eventA, eventB)
	suite.whenQueueIsForwarded()
	suite.thenForwardedEventsShouldBe(eventA, eventB)
}

func (suite *QueueSuite) TestQueueIsEmptyAfterForward() {
	suite.givenQueuedEvents(suite.someEvent(), suite.someEvent())
	suite.whenQueueIsForwarded()
	suite.thenQueueShouldBeEmpty()
}

func (suite *QueueSuite) TestEventsCanBeAddedWhileForwarding() {
	eventA := suite.keyedEvent(1)
	eventB := suite.keyedEvent(2)
	eventC := suite.keyedEvent(3)
	suite.givenQueuedEvents(eventA, eventB)
	suite.givenEventWillBeQueuedAtNextForward(eventC)
	suite.whenQueueIsForwarded()
	suite.thenForwardedEventsShouldBe(eventA, eventB)

	suite.whenQueueIsForwarded()
	suite.thenForwardedEventsShouldBe(eventA, eventB, eventC)
}

func (suite *QueueSuite) givenQueuedEvents(events ...event.Event) {
	for _, e := range events {
		suite.queue.Event(e)
	}
}

func (suite *QueueSuite) givenEventWillBeQueuedAtNextForward(e event.Event) {
	suite.pendingEvent = e
}

func (suite *QueueSuite) whenAddingEvent(e event.Event) {
	suite.queue.Event(e)
}

func (suite *QueueSuite) whenQueueIsForwarded() {
	suite.queue.ForwardTo(suite)
}

func (suite *QueueSuite) thenQueueShouldBeEmpty() {
	assert.True(suite.T(), suite.queue.IsEmpty(), "Queue is not empty")
}

func (suite *QueueSuite) thenQueueShouldNotBeEmpty() {
	assert.False(suite.T(), suite.queue.IsEmpty(), "Queue is empty")
}

func (suite *QueueSuite) thenForwardedEventsShouldBe(expected ...event.Event) {
	assert.Equal(suite.T(), expected, suite.events, "Events don't match")
}

func (suite *QueueSuite) Event(e event.Event) {
	suite.events = append(suite.events, e)
	if suite.pendingEvent != nil {
		suite.queue.Event(suite.pendingEvent)
		suite.pendingEvent = nil
	}
}

func (suite *QueueSuite) someEvent() event.Event {
	return struct {
		data string
	}{"test"}
}

func (suite *QueueSuite) keyedEvent(key int) event.Event {
	return struct {
		key int
	}{key}
}
