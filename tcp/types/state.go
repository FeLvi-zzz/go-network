package types

type State uint8

const (
	State_CLOSED     State = iota
	State_LISTEN     State = iota
	State_SYN_RCVD   State = iota
	State_SYN_SENT   State = iota
	State_ESTAB      State = iota
	State_FIN_WAIT_1 State = iota
	State_FIN_WAIT_2 State = iota
	State_CLOSING    State = iota
	State_CLOSE_WAIT State = iota
	State_LAST_ACK   State = iota
	State_TIME_WAIT  State = iota
)

func (s State) String() string {
	switch s {
	case State_CLOSED:
		return "CLOSED"
	case State_LISTEN:
		return "LISTEN"
	case State_SYN_RCVD:
		return "SYN_RCVD"
	case State_SYN_SENT:
		return "SYN_SENT"
	case State_ESTAB:
		return "ESTAB"
	case State_FIN_WAIT_1:
		return "FIN_WAIT_1"
	case State_FIN_WAIT_2:
		return "FIN_WAIT_2"
	case State_CLOSING:
		return "CLOSING"
	case State_CLOSE_WAIT:
		return "CLOSE_WAIT"
	case State_LAST_ACK:
		return "LAST_ACK"
	case State_TIME_WAIT:
		return "TIME_WAIT"
	}
	return "UNKNOWN"
}
