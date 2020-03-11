package snow

// DefaultSnow is the global default snowflake node object.
// nolint gochecknoglobals
var DefaultSnow, _ = NewNode()

// GetOption return the option.
func GetOption() Option { return DefaultSnow.option }

// GetEpoch returns an int64 epoch is snowflake epoch in milliseconds.
func GetEpoch() int64 { return DefaultSnow.epoch }

// GetTime returns an int64 unix timestamp in milliseconds of the snowflake ID time.
func GetTime() int64 { return DefaultSnow.time }

// GetNodeID returns an int64 of the snowflake ID node number
func GetNodeID() int64 { return DefaultSnow.nodeID }

// GetStep returns an int64 of the snowflake step (or sequence) number
func GetStep() int64 { return DefaultSnow.step }

// Next creates and returns a unique snowflake ID
// To help guarantee uniqueness
// - Make sure your system is keeping accurate system time
// - Make sure you never have multiple nodes running with the same node ID
func Next() ID { return DefaultSnow.Next() }
