package snow

// An ID is a custom type used for a snowflake ID.  This is used so we can attach methods onto the ID.
type ID int64

// TimeOf returns an int64 unix timestamp in milliseconds of the snowflake ID time
func (n *Node) TimeOf(f ID) int64 { return (int64(f) >> n.timeShift) + n.epoch }

// NodeIDOf returns an int64 of the snowflake ID node number
func (n *Node) NodeIDOf(f ID) int64 { return int64(f) & n.nodeMask >> n.nodeShift }

// StepOf returns an int64 of the snowflake step (or sequence) number
func (n *Node) StepOf(f ID) int64 { return int64(f) & n.stepMask }