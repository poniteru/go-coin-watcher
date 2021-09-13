package market

type Direction int

const (
	UNDEFINED Direction = iota
	UP
	DOWN
)

func (d Direction) String() string {
	switch d {
	case UNDEFINED:
		return "UNDEFINED"
	case UP:
		return "UP"
	case DOWN:
		return "DOWN"
	}
	return "UNDEFINED"
}

func GetDirection(s string) Direction {
	switch s {
	case "up":
		fallthrough
	case "UP":
		return UP
	case "down":
		fallthrough
	case "DOWN":
		return DOWN
	default:
		return UNDEFINED
	}
}
