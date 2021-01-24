package gnet

import "time"

type (
	// EventHandler represents the server events' callbacks for the Serve call.
	// Each event has an Action return value that is used manage the state
	// of the connection and server.
	EventHandler interface {
		// OnInitComplete fires when the server is ready for accepting connections.
		// The parameter:server has information and various utilities.
		OnInitComplete(server Server) (action Action)

		// OnShutdown fires when the server is being shut down, it is called right after
		// all event-loops and connections are closed.
		OnShutdown(server Server)

		// OnOpened fires when a new connection has been opened.
		// The parameter:c has information about the connection such as it's local and remote address.
		// Parameter:out is the return value which is going to be sent back to the client.
		OnOpened(c Conn) (out []byte, action Action)

		// OnClosed fires when a connection has been closed.
		// The parameter:err is the last known connection error.
		OnClosed(c Conn, err error) (action Action)

		// PreWrite fires just before any data is written to any client socket, this event function is usually used to
		// put some code of logging/counting/reporting or any prepositive operations before writing data to client.
		PreWrite()

		// React fires when a connection sends the server data.
		// Call c.Read() or c.ReadN(n) within the parameter:c to read incoming data from client.
		// Parameter:out is the return value which is going to be sent back to the client.
		React(frame []byte, c Conn) (out []byte, action Action)

		// Tick fires immediately after the server starts and will fire again
		// following the duration specified by the delay return value.
		Tick() (delay time.Duration, action Action)
	}
)

// EventServer is a built-in implementation of EventHandler which sets up each method with a default implementation,
// you can compose it with your own implementation of EventHandler when you don't want to implement all methods
// in EventHandler.
type EventServer struct {
}

// OnInitComplete fires when the server is ready for accepting connections.
// The parameter:server has information and various utilities.
func (es *EventServer) OnInitComplete(svr Server) (action Action) {
	return
}

// OnShutdown fires when the server is being shut down, it is called right after
// all event-loops and connections are closed.
func (es *EventServer) OnShutdown(svr Server) {
}

// OnOpened fires when a new connection has been opened.
// The parameter:c has information about the connection such as it's local and remote address.
// Parameter:out is the return value which is going to be sent back to the client.
func (es *EventServer) OnOpened(c Conn) (out []byte, action Action) {
	return
}

// OnClosed fires when a connection has been closed.
// The parameter:err is the last known connection error.
func (es *EventServer) OnClosed(c Conn, err error) (action Action) {
	return
}

// PreWrite fires just before any data is written to any client socket, this event function is usually used to
// put some code of logging/counting/reporting or any prepositive operations before writing data to client.
func (es *EventServer) PreWrite() {
}

// React fires when a connection sends the server data.
// Call c.Read() or c.ReadN(n) within the parameter:c to read incoming data from client.
// Parameter:out is the return value which is going to be sent back to the client.
func (es *EventServer) React(frame []byte, c Conn) (out []byte, action Action) {
	return
}

// Tick fires immediately after the server starts and will fire again
// following the duration specified by the delay return value.
func (es *EventServer) Tick() (delay time.Duration, action Action) {
	return
}
