package netx

type Check interface {
	Valid() bool
	Report() string
}

type LocalInvariant interface {
	LocalCheck() bool
}

type GlobalInvariant interface {
	GlobalCheck() bool
}

type AnomalyDetector interface {
	Detect() bool
}
