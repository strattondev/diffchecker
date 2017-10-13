package diffchecker

type DiffCheckerExpiry int

const (
	HOUR DiffCheckerExpiry = iota
	DAY
	FOREVER
)

func (expiry DiffCheckerExpiry) String() string {
	switch expiry {
	case HOUR:
		return "hour"
	case DAY:
		return "day"
	case FOREVER:
		fallthrough
	default:
		return "forever"
	}
}
