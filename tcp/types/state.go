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
