package event_test

import (
	"reflect"
	"testing"

	"github.com/inkyblackness/hacked/editor/event"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DispatcherSuite struct {
	suite.Suite

	dispatcher *event.Dispatcher
}

func TestDispatcherSuite(t *testing.T) {
	suite.Run(t, new(DispatcherSuite))
}

func (suite *DispatcherSuite) SetupTest() {
	suite.dispatcher = event.NewDispatcher()
}

func (suite *DispatcherSuite) TestRegisterHandlerPanicsWithWrongHandlerType() {
	type actualEvent struct{}
	eType := reflect.TypeOf((*actualEvent)(nil)).Elem()

	assert.Panics(suite.T(), func() {
		suite.dispatcher.RegisterHandler(reflect.TypeOf(int(0)), func(int) {})
	}, "Panic expected for event type not being a struct")
	assert.Panics(suite.T(), func() {
		suite.dispatcher.RegisterHandler(eType, 123)
	}, "Panic expected for not being a function")
	assert.Panics(suite.T(), func() {
		suite.dispatcher.RegisterHandler(eType, func() {})
	}, "Panic expected for missing parameter")
	assert.Panics(suite.T(), func() {
		suite.dispatcher.RegisterHandler(eType, func(int) {})
	}, "Panic expected for wrong parameter")
	assert.Panics(suite.T(), func() {
		suite.dispatcher.RegisterHandler(eType, func(*actualEvent, int) {})
	}, "Panic expected for too many parameters")
}

func (suite *DispatcherSuite) TestEventDispatchesGivenEventToRegisteredHandler() {
	type testedEvent struct{ key int }
	eventA := testedEvent{123}
	var received *testedEvent
	suite.givenRegisteredHandler(reflect.TypeOf(eventA), func(e testedEvent) { received = &e })
	suite.whenEventIsDispatched(eventA)
	require.NotNil(suite.T(), received, "No event received")
	assert.Equal(suite.T(), eventA, *received, "Wrong event received")
}

func (suite *DispatcherSuite) TestEventHandlerCanBeUnregistered() {
	type testedEvent struct{ key int }
	eventA := testedEvent{123}
	var received *testedEvent
	handler := func(e testedEvent) { received = &e }
	suite.givenRegisteredHandler(reflect.TypeOf(eventA), handler)
	suite.givenHandlerWasUnregistered(reflect.TypeOf(eventA), handler)
	suite.whenEventIsDispatched(eventA)
	require.Nil(suite.T(), received, "No call expected")
}

func (suite *DispatcherSuite) TestDeregistrationKeepsOtherHandlersAlive() {
	type testedEvent struct{ key int }
	type differentEvent struct{ key int }
	event1 := testedEvent{123}
	event2 := differentEvent{456}
	var receivedA *testedEvent
	var receivedB *testedEvent
	var receivedC *testedEvent
	var receivedD *differentEvent
	handlerA := func(e testedEvent) { receivedA = &e }
	handlerB := func(e testedEvent) { receivedB = &e }
	handlerC := func(e testedEvent) { receivedC = &e }
	handlerD := func(e differentEvent) { receivedD = &e }
	suite.givenRegisteredHandler(reflect.TypeOf(event1), handlerA, handlerB, handlerC, handlerB)
	suite.givenRegisteredHandler(reflect.TypeOf(event2), handlerD)
	suite.givenHandlerWasUnregistered(reflect.TypeOf(event1), handlerB)
	suite.whenEventIsDispatched(event1)
	suite.whenEventIsDispatched(event2)
	assert.NotNil(suite.T(), receivedA, "Call expected via A")
	assert.Nil(suite.T(), receivedB, "No call expected via B")
	assert.NotNil(suite.T(), receivedC, "Call expected via C")
	assert.NotNil(suite.T(), receivedD, "Call expected via D")
}

func (suite *DispatcherSuite) TestDeregistrationOfHandlerDuringDispatchIsConsidered() {
	type testedEvent struct{ key int }
	eventA := testedEvent{123}
	calls := 0
	handler := func(e testedEvent) { calls++ }
	deregisteringHandler := func(e testedEvent) { suite.dispatcher.UnregisterHandler(reflect.TypeOf(eventA), handler) }
	otherHandlerCalledTimes := 0
	otherHandler := func(e testedEvent) { otherHandlerCalledTimes++ }
	suite.givenRegisteredHandler(reflect.TypeOf(eventA), handler, handler, deregisteringHandler, handler, handler, otherHandler)
	suite.whenEventIsDispatched(eventA)
	assert.Equal(suite.T(), 2, calls, "More calls than expected")
	assert.Equal(suite.T(), 1, otherHandlerCalledTimes, "Handler list is wrong")
}

func (suite *DispatcherSuite) TestRegistrationOfHandlerDuringDispatchIsIgnored() {
	type testedEvent struct{ key int }
	eventA := testedEvent{123}
	var received *testedEvent
	handler := func(e testedEvent) { received = &e }
	registeringHandler := func(e testedEvent) { suite.dispatcher.RegisterHandler(reflect.TypeOf(eventA), handler) }
	dummyHandler := func(e testedEvent) {}
	suite.givenRegisteredHandler(reflect.TypeOf(eventA), registeringHandler, dummyHandler)
	suite.givenHandlerWasUnregistered(reflect.TypeOf(eventA), dummyHandler) // in case same array was reused through cap
	suite.whenEventIsDispatched(eventA)
	require.Nil(suite.T(), received, "No call expected")
}

func (suite *DispatcherSuite) TestEventCanNotBeDispatchedDuringDispatch() {
	type testedEvent struct{ key int }
	eventA := testedEvent{123}
	eventB := testedEvent{456}
	calls := 0
	handler := func(e testedEvent) {
		calls++
		if calls == 1 {
			assert.Panics(suite.T(), func() {
				suite.dispatcher.Event(eventB)
			}, "panic expected")
		}
	}
	suite.givenRegisteredHandler(reflect.TypeOf(eventA), handler)
	suite.whenEventIsDispatched(eventA)
	assert.Equal(suite.T(), 1, calls, "Only one call expected")
}

func (suite *DispatcherSuite) TestRegistrationAndDeregistrationCanBeStackedAndStillBeAsExpected_A() {
	type testedEvent struct{ key int }
	eventA := testedEvent{123}
	calls := 0
	handler := func(e testedEvent) { calls++ }
	registeringHandler := func(e testedEvent) { suite.dispatcher.RegisterHandler(reflect.TypeOf(eventA), handler) }
	unregisteringHandler := func(e testedEvent) { suite.dispatcher.UnregisterHandler(reflect.TypeOf(eventA), handler) }
	suite.givenRegisteredHandler(reflect.TypeOf(eventA), unregisteringHandler, registeringHandler)
	suite.givenEventWasDispatched(eventA)
	require.Equal(suite.T(), 0, calls, "No call expected during first pass")
	suite.givenHandlerWasUnregistered(reflect.TypeOf(eventA), unregisteringHandler)
	suite.givenHandlerWasUnregistered(reflect.TypeOf(eventA), registeringHandler)
	suite.whenEventIsDispatched(eventA)
	assert.Equal(suite.T(), 1, calls, "One call expected as the last call was a registration")
}

func (suite *DispatcherSuite) TestRegistrationAndDeregistrationCanBeStackedAndStillBeAsExpected_B() {
	type testedEvent struct{ key int }
	eventA := testedEvent{123}
	calls := 0
	handler := func(e testedEvent) { calls++ }
	registeringHandler := func(e testedEvent) { suite.dispatcher.RegisterHandler(reflect.TypeOf(eventA), handler) }
	unregisteringHandler := func(e testedEvent) { suite.dispatcher.UnregisterHandler(reflect.TypeOf(eventA), handler) }
	suite.givenRegisteredHandler(reflect.TypeOf(eventA), registeringHandler, unregisteringHandler)
	suite.givenEventWasDispatched(eventA)
	require.Equal(suite.T(), 0, calls, "No call expected during first pass as it was registered during event chain")
	suite.givenHandlerWasUnregistered(reflect.TypeOf(eventA), unregisteringHandler)
	suite.givenHandlerWasUnregistered(reflect.TypeOf(eventA), registeringHandler)
	suite.whenEventIsDispatched(eventA)
	assert.Equal(suite.T(), 0, calls, "No call expected as the deregistration happened after registration")
}

func (suite *DispatcherSuite) TestRegistrationAndDeregistrationCanBeStackedAndStillBeAsExpected_C() {
	type testedEvent struct{ key int }
	eventA := testedEvent{123}
	calls := 0
	handler := func(e testedEvent) { calls++ }
	registeringHandler := func(e testedEvent) { suite.dispatcher.RegisterHandler(reflect.TypeOf(eventA), handler) }
	unregisteringHandler := func(e testedEvent) { suite.dispatcher.UnregisterHandler(reflect.TypeOf(eventA), handler) }
	suite.givenRegisteredHandler(reflect.TypeOf(eventA), registeringHandler, unregisteringHandler, registeringHandler)
	suite.givenEventWasDispatched(eventA)
	require.Equal(suite.T(), 0, calls, "No call expected during first pass as it was registered during event chain")
	suite.givenHandlerWasUnregistered(reflect.TypeOf(eventA), unregisteringHandler)
	suite.givenHandlerWasUnregistered(reflect.TypeOf(eventA), registeringHandler)
	suite.whenEventIsDispatched(eventA)
	assert.Equal(suite.T(), 1, calls, "One call expected as there was one additional registration")
}

func (suite *DispatcherSuite) TestRegistrationReturnsDeregisterShortcut() {
	type testedEvent struct{ key int }
	eventA := testedEvent{123}
	calls := 0
	handler := func(e testedEvent) { calls++ }
	unreg := suite.dispatcher.RegisterHandler(reflect.TypeOf(eventA), handler)
	unreg()
	suite.whenEventIsDispatched(eventA)
	assert.Equal(suite.T(), 0, calls, "No call expected")
}

func (suite *DispatcherSuite) givenRegisteredHandler(t reflect.Type, handlers ...interface{}) {
	for _, handler := range handlers {
		require.NotPanics(suite.T(), func() {
			suite.dispatcher.RegisterHandler(t, handler)
		}, "Registering handler panicked! Not expected.")
	}
}

func (suite *DispatcherSuite) givenHandlerWasUnregistered(t reflect.Type, handler interface{}) {
	suite.dispatcher.UnregisterHandler(t, handler)
}

func (suite *DispatcherSuite) givenEventWasDispatched(e event.Event) {
	suite.whenEventIsDispatched(e)
}

func (suite *DispatcherSuite) whenEventIsDispatched(e event.Event) {
	suite.dispatcher.Event(e)
}
