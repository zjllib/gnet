// Copyright (c) 2019 Andy Pan
// Copyright (c) 2018 Joshua J Baker
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package gnet

import (
	"github.com/zjllib/gnet/errors"
	"github.com/zjllib/gnet/internal/logging"
	"net"
	"runtime"
	"sync/atomic"
	"time"
)

// Action is an action that occurs after the completion of an event.
type Action int

const (
	None     Action = iota // indicates that no action should occur following an event.
	Close                  // closes the connection.
	Shutdown               // shutdowns the server.
)

// Server represents a server context which provides information about the
// running server and has control functions for managing state.
type Server struct {
	// svr is the internal server struct.
	svr *server
	// Multicore indicates whether the server will be effectively created with multi-cores, if so,
	// then you must take care of synchronizing the shared data between all event callbacks, otherwise,
	// it will run the server with single thread. The number of threads in the server will be automatically
	// assigned to the value of logical CPUs usable by the current process.
	Multicore bool

	// The Addr parameter is the listening address that align
	// with the addr string passed to the Serve function.
	Addr net.Addr

	// NumEventLoop is the number of event-loops that the server is using.
	NumEventLoop int

	// ReusePort indicates whether SO_REUSEPORT is enable.
	ReusePort bool

	// TCPKeepAlive (SO_KEEPALIVE) socket option.
	TCPKeepAlive time.Duration
}

// CountConnections counts the number of currently active connections and returns it.
func (s Server) CountConnections() (count int) {
	s.svr.subEventLoopSet.iterate(func(i int, el *eventloop) bool {
		count += int(atomic.LoadInt32(&el.connCount))
		return true
	})
	return
}

// Serve starts handling events for the specified address.
//
// Address should use a scheme prefix and be formatted
// like `tcp://192.168.0.10:9851` or `unix://socket`.
// Valid network schemes:
//  tcp   - bind to both IPv4 and IPv6
//  tcp4  - IPv4
//  tcp6  - IPv6
//  udp   - bind to both IPv4 and IPv6
//  udp4  - IPv4
//  udp6  - IPv6
//  unix  - Unix Domain Socket
//
// The "tcp" network scheme is assumed when one is not specified.
func Serve(eventHandler EventHandler, protoAddr string, opts ...Option) (err error) {
	options := loadOptions(opts...)

	if options.Logger != nil {
		logging.DefaultLogger = options.Logger
	}
	defer logging.Cleanup()

	// The maximum number of operating system threads that the Go program can use is initially set to 10000,
	// which should be the maximum amount of I/O event-loops locked to OS threads users can start up.
	if options.LockOSThread && options.NumEventLoop > 10000 {
		logging.DefaultLogger.Errorf("too many event-loops under LockOSThread mode, should be less than 10,000 "+
			"while you are trying to set up %d\n", options.NumEventLoop)
		return errors.ErrTooManyEventLoopThreads
	}

	network, addr := parseProtoAddr(protoAddr)

	var ln *listener
	if ln, err = initListener(network, addr, options.ReusePort); err != nil {
		return
	}
	defer ln.close()

	return serve(eventHandler, ln, options)
}

func sniffErrorAndLog(err error) {
	if err != nil {
		logging.DefaultLogger.Errorf(err.Error())
	}
}

// channelBuffer determines whether the channel should be a buffered channel to get the best performance.
func channelBuffer(preset int) int {
	// Use blocking channel if GOMAXPROCS=1.
	// This switches context from sender to receiver immediately,
	// which results in higher performance (under go1.5 at least).
	if runtime.GOMAXPROCS(0) == 1 {
		return 0
	}

	// Use non-blocking workerChan if GOMAXPROCS>1,
	// since otherwise the sender might be dragged down if the receiver is CPU-bound.
	return preset
}
