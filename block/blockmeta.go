package block

type blockmeta struct {
	// which log file?
	Log int
	// the offset in the log file
	LogOffset int64
}
