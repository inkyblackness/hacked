package movie

// Container wraps the information and data of a MOVI container.
type Container struct {
	// EndTimestamp is the time of the end of the movie.
	EndTimestamp Timestamp

	// TODO: merge by Write(), ordered by bucket priority
	// TODO: remove other members, they should all no longer be necessary in the end.
	Audio     Audio
	Video     Video
	Subtitles Subtitles
}
