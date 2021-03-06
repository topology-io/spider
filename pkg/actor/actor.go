package actor

import (
	"time"

	"github.com/topology-io/spider/pkg/promise"
)

// Actor is a concurrency abstraction.
//
// Spider guarantees the following two "happens before" rules:
//
// 1. The actor send rule: the send of the message to an actor happens before
//    the receive of that message by the same actor.
// 2. The actor subsequent processing rule: processing one message happens
//    before processing the next message by the same sender.
type Actor interface {
	// Prestart is invoked by the actor system after this actor has been created
	// but before it begins processing its mailbox.
	// Actor implementations should save the context locally for later use, e.g.
	// `otherActor.Send(cx.Self, "ping")`
	Prestart(cx Context)

	// Receive is invoked by the actor system after a message has been sent
	// to this actor. The actor system treats the entire Receive method as
	// a critical section; the actor will not receive subsequent messages
	// until returning. For this reason, do your best to avoid blocking in the
	// receive method.
	//
	// Message values should be immutable, or at least mutation should be
	// avoided. The type system doesn't let us enforce that, so just be careful.
	Receive(rx ReceiveContext, msg interface{})
}

// Ref is an opaque interface to an actor in the system.
type Ref interface {
	// Returns this actor ref's address.
	Address() Address

	// Sends a fire-and-forget message with no delivery guarantees.
	//
	// Message values should be immutable, or at least mutation should be
	// avoided. The type system doesn't let us enforce that, so just be careful.
	Send(replyTo Ref, msg interface{})

	// Sends a message and awaits a reply up to a timeout duration. Note that
	// Ask is more expensive and less scalable than one-way message passing
	// with Send.
	//
	// Message values should be immutable, or at least mutation should be
	// avoided. The type system doesn't let us enforce that, so just be careful.
	Ask(msg interface{}, timeout time.Duration) (interface{}, error)

	// AddWatcher provides death notifications in the actor system.
	// The watcher will receive an actor.Stopped message after this actor ref
	// dies.
	AddWatcher(Ref)

	// Life returns a promise that will complete after the underlying actor is
	// stopped.
	Life() promise.ReadOnlyPromise
}

// Address encodes the receive address of an actor ref.
type Address string
