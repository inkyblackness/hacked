package input_test

import (
	"testing"

	"github.com/inkyblackness/hacked/ui/input"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type keyEvent struct {
	isKey bool
	key   input.Key
	mod   input.Modifier
}

type testingStickyKeysListener struct {
	eventMap map[input.Key][]keyEvent
	events   []keyEvent
}

func (listener *testingStickyKeysListener) Key(key input.Key, modifier input.Modifier) {
	listener.addEvent(keyEvent{true, key, modifier})
}

func (listener *testingStickyKeysListener) Modifier(modifier input.Modifier) {
	listener.addEvent(keyEvent{false, 0, modifier})
}

func (listener *testingStickyKeysListener) addEvent(event keyEvent) {
	listener.eventMap[event.key] = append(listener.eventMap[event.key], event)
	listener.events = append(listener.events, event)
}

type StickyKeyBufferSuite struct {
	suite.Suite
	buffer   *input.StickyKeyBuffer
	listener *testingStickyKeysListener
}

func TestProviderBackedStoreSuite(t *testing.T) {
	suite.Run(t, new(StickyKeyBufferSuite))
}

func (suite *StickyKeyBufferSuite) SetupTest() {
	suite.listener = &testingStickyKeysListener{
		eventMap: make(map[input.Key][]keyEvent)}
	suite.buffer = input.NewStickyKeyBuffer(suite.listener)
}

func (suite *StickyKeyBufferSuite) TestRegularEventsAreForwarded_A() {
	suite.buffer.KeyDown(input.KeyF1, input.ModNone)
	suite.buffer.KeyUp(input.KeyF1, input.ModNone)
	suite.buffer.KeyDown(input.KeyF2, input.ModNone)
	suite.buffer.KeyUp(input.KeyF2, input.ModNone)

	assert.Equal(suite.T(), []keyEvent{{true, input.KeyF1, input.ModNone}}, suite.listener.eventMap[input.KeyF1])
	assert.Equal(suite.T(), []keyEvent{{true, input.KeyF2, input.ModNone}}, suite.listener.eventMap[input.KeyF2])
}

func (suite *StickyKeyBufferSuite) TestRegularEventsAreForwarded_B() {
	suite.buffer.KeyDown(input.KeyF1, input.ModNone)
	suite.buffer.KeyDown(input.KeyF2, input.ModNone)
	suite.buffer.KeyUp(input.KeyF2, input.ModNone)
	suite.buffer.KeyUp(input.KeyF1, input.ModNone)

	assert.Equal(suite.T(), []keyEvent{{true, input.KeyF1, input.ModNone}}, suite.listener.eventMap[input.KeyF1])
	assert.Equal(suite.T(), []keyEvent{{true, input.KeyF2, input.ModNone}}, suite.listener.eventMap[input.KeyF2])
	assert.Equal(suite.T(), []keyEvent{{true, input.KeyF1, input.ModNone}, {true, input.KeyF2, input.ModNone}}, suite.listener.events)
}

func (suite *StickyKeyBufferSuite) TestIdenticalKeysAreReportedOnlyOnce() {
	suite.buffer.KeyDown(input.KeyF1, input.ModNone)
	suite.buffer.KeyDown(input.KeyF1, input.ModNone)
	suite.buffer.KeyDown(input.KeyF1, input.ModNone)
	suite.buffer.KeyUp(input.KeyF1, input.ModNone)
	suite.buffer.KeyUp(input.KeyF1, input.ModNone)
	suite.buffer.KeyUp(input.KeyF1, input.ModNone)

	assert.Equal(suite.T(), []keyEvent{{true, input.KeyF1, input.ModNone}}, suite.listener.eventMap[input.KeyF1])
	assert.Equal(suite.T(), []keyEvent{{true, input.KeyF1, input.ModNone}}, suite.listener.events)
}

func (suite *StickyKeyBufferSuite) TestSuperfluousReleasesAreIgnored() {
	suite.buffer.KeyDown(input.KeyF1, input.ModNone)
	suite.buffer.KeyUp(input.KeyF1, input.ModNone)
	suite.buffer.KeyUp(input.KeyF1, input.ModNone)
	suite.buffer.KeyUp(input.KeyF1, input.ModNone)
	suite.buffer.KeyDown(input.KeyF2, input.ModNone)
	suite.buffer.KeyDown(input.KeyF1, input.ModNone)
	suite.buffer.KeyUp(input.KeyF1, input.ModNone)

	assert.Equal(suite.T(), []keyEvent{
		{true, input.KeyF1, input.ModNone}, {true, input.KeyF1, input.ModNone}}, suite.listener.eventMap[input.KeyF1])
	assert.Equal(suite.T(), []keyEvent{
		{true, input.KeyF1, input.ModNone},
		{true, input.KeyF2, input.ModNone},
		{true, input.KeyF1, input.ModNone}}, suite.listener.events)
}

func (suite *StickyKeyBufferSuite) TestReleaseAllResetsRegularKeys() {
	suite.buffer.KeyDown(input.KeyEnter, input.ModNone)
	suite.buffer.KeyDown(input.KeyEnter, input.ModNone)
	suite.buffer.KeyDown(input.KeyTab, input.ModNone)
	suite.buffer.ReleaseAll()
	suite.buffer.KeyDown(input.KeyTab, input.ModNone)

	assert.Equal(suite.T(), []keyEvent{{true, input.KeyEnter, input.ModNone}}, suite.listener.eventMap[input.KeyEnter])
	assert.Equal(suite.T(), []keyEvent{{true, input.KeyTab, input.ModNone}, {true, input.KeyTab, input.ModNone}},
		suite.listener.eventMap[input.KeyTab])
}

func (suite *StickyKeyBufferSuite) TestModifierEventsAreConvertedToModifierStates_DownAugmented() {
	suite.buffer.KeyDown(input.KeyShift, input.ModShift)
	suite.buffer.KeyDown(input.KeyAlt, input.ModShift.With(input.ModAlt))
	suite.buffer.KeyUp(input.KeyAlt, input.ModShift)
	suite.buffer.KeyUp(input.KeyShift, input.ModNone)

	assert.Equal(suite.T(), 0, len(suite.listener.eventMap[input.KeyShift]))
	assert.Equal(suite.T(), 0, len(suite.listener.eventMap[input.KeyAlt]))
	assert.Equal(suite.T(),
		[]keyEvent{
			{false, 0, input.ModShift},
			{false, 0, input.ModShift.With(input.ModAlt)},
			{false, 0, input.ModShift},
			{false, 0, input.ModNone},
		},
		suite.listener.events)
}

func (suite *StickyKeyBufferSuite) TestModifierEventsAreConvertedToModifierStates_UpAugmented() {
	suite.buffer.KeyDown(input.KeyShift, input.ModNone)
	suite.buffer.KeyDown(input.KeyAlt, input.ModShift)
	suite.buffer.KeyUp(input.KeyAlt, input.ModShift.With(input.ModAlt))
	suite.buffer.KeyUp(input.KeyShift, input.ModShift)

	assert.Equal(suite.T(), 0, len(suite.listener.eventMap[input.KeyShift]))
	assert.Equal(suite.T(), 0, len(suite.listener.eventMap[input.KeyAlt]))
	assert.Equal(suite.T(),
		[]keyEvent{
			{false, 0, input.ModShift},
			{false, 0, input.ModShift.With(input.ModAlt)},
			{false, 0, input.ModShift},
			{false, 0, input.ModNone},
		},
		suite.listener.events)
}

func (suite *StickyKeyBufferSuite) TestReleaseAllReleasesModifiers() {
	suite.buffer.KeyDown(input.KeyShift, input.ModNone)
	suite.buffer.KeyDown(input.KeyShift, input.ModNone)
	suite.buffer.KeyDown(input.KeyAlt, input.ModShift)
	suite.buffer.ReleaseAll()

	assert.Equal(suite.T(),
		[]keyEvent{
			{false, 0, input.ModShift},
			{false, 0, input.ModShift.With(input.ModAlt)},
			{false, 0, input.ModNone},
		},
		suite.listener.events)
}

func (suite *StickyKeyBufferSuite) TestActiveModifierReturnsCurrentModifier() {
	suite.buffer.KeyDown(input.KeyShift, input.ModNone)
	suite.buffer.KeyDown(input.KeyControl, input.ModShift)
	suite.buffer.KeyDown(input.KeyAlt, input.ModShift.With(input.ModControl))
	suite.buffer.KeyUp(input.KeyControl, input.ModShift.With(input.ModControl).With(input.ModAlt))

	assert.Equal(suite.T(), input.ModShift.With(input.ModAlt), suite.buffer.ActiveModifier())
}

func (suite *StickyKeyBufferSuite) TestModifierAreUpdatedOnOtherKeys_A() {
	suite.buffer.KeyDown(input.KeyShift, input.ModShift)
	suite.buffer.KeyDown(input.KeyTab, input.ModNone)

	assert.Equal(suite.T(),
		[]keyEvent{
			{false, 0, input.ModShift},
			{false, 0, input.ModNone},
			{true, input.KeyTab, input.ModNone},
		},
		suite.listener.events)
}

func (suite *StickyKeyBufferSuite) TestModifierAreUpdatedOnOtherKeys_B() {
	suite.buffer.KeyDown(input.KeyTab, input.ModNone)
	suite.buffer.KeyDown(input.KeyShift, input.ModShift)
	suite.buffer.KeyUp(input.KeyTab, input.ModNone)

	assert.Equal(suite.T(),
		[]keyEvent{
			{true, input.KeyTab, input.ModNone},
			{false, 0, input.ModShift},
			{false, 0, input.ModNone},
		},
		suite.listener.events)
}

func (suite *StickyKeyBufferSuite) TestModifierAreUpdatedOnOtherKeysWithoutThemselves() {
	suite.buffer.KeyDown(input.KeyShift, input.ModShift)
	suite.buffer.KeyDown(input.KeyShift, input.ModShift)
	suite.buffer.KeyUp(input.KeyShift, input.ModShift)

	assert.Equal(suite.T(), []keyEvent{{false, 0, input.ModShift}}, suite.listener.events)
}
