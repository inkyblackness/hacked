package world

// IsConsideredCyberspaceByDefault returns true for those level identifier
// that are hardcoded to be cyberspace levels.
// For the vanilla engine, this is true for level 10 and any above 13.
func IsConsideredCyberspaceByDefault(levelID int) bool {
	return (levelID == 10) || (levelID > 13)
}
