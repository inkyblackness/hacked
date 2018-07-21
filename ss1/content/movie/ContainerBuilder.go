package movie

import (
	"github.com/inkyblackness/hacked/ss1/content/bitmap"
)

// ContainerBuilder is the builder implementation for a Container.
type ContainerBuilder struct {
	container *memoryContainer
}

// NewContainerBuilder returns a new builder for creating a new container.
func NewContainerBuilder() *ContainerBuilder {
	builder := &ContainerBuilder{container: &memoryContainer{}}

	return builder
}

// Build returns the immutable instance of a new container.
func (builder *ContainerBuilder) Build() Container {
	return builder.container
}

// MediaDuration sets the duration for the new container in seconds.
func (builder *ContainerBuilder) MediaDuration(value float32) *ContainerBuilder {
	builder.container.mediaDuration = value
	return builder
}

// VideoHeight sets the video height for the new container.
func (builder *ContainerBuilder) VideoHeight(value uint16) *ContainerBuilder {
	builder.container.videoHeight = value
	return builder
}

// VideoWidth sets the video width for the new container.
func (builder *ContainerBuilder) VideoWidth(value uint16) *ContainerBuilder {
	builder.container.videoWidth = value
	return builder
}

// StartPalette sets the initial palette of the new container
func (builder *ContainerBuilder) StartPalette(palette bitmap.Palette) *ContainerBuilder {
	builder.container.startPalette = palette
	return builder
}

// AudioSampleRate sets the video width for the new container.
func (builder *ContainerBuilder) AudioSampleRate(value uint16) *ContainerBuilder {
	builder.container.audioSampleRate = value
	return builder
}

// AddEntry adds the given entry to the list.
func (builder *ContainerBuilder) AddEntry(entry Entry) *ContainerBuilder {
	builder.container.entries = append(builder.container.entries, entry)
	return builder
}
