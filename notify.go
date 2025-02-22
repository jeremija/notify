package notify

import "github.com/pkg/errors"

// Compile-time check to ensure Notify implements Notifier.
var _ Notifier = (*Notify)(nil)

// ErrSendNotification signals that the notifier failed to send a notification.
var ErrSendNotification = errors.New("send notification")

// Notify is the central struct for managing notification services and sending messages to them.
type Notify struct {
	Disabled  bool
	notifiers []Notifier
}

// Option is a function that can be used to configure a Notify instance. It is used by the WithOptions and
// NewWithOptions functions. It's a function because it's a bit more flexible than using a struct. The only required
// parameter is the Notify instance itself.
type Option func(*Notify)

// Enable is an Option function that enables the Notify instance. This is the default behavior.
func Enable(n *Notify) {
	if n != nil {
		n.Disabled = false
	}
}

// Disable is an Option function that disables the Notify instance. It is enabled by default.
func Disable(n *Notify) {
	if n != nil {
		n.Disabled = true
	}
}

// WithOptions applies the given options to the Notify instance. If no options are provided, it returns the Notify
// instance unchanged.
func (n *Notify) WithOptions(options ...Option) *Notify {
	if options == nil {
		return n
	}

	for _, option := range options {
		if option != nil {
			option(n)
		}
	}

	return n
}

// NewWithOptions returns a new instance of Notify with the given options. If no options are provided, it returns a new
// Notify instance with default options. By default, the Notify instance is enabled.
func NewWithOptions(options ...Option) *Notify {
	n := &Notify{
		Disabled:  false,               // Enabled by default.
		notifiers: make([]Notifier, 0), // Avoid nil list.
	}

	return n.WithOptions(options...)
}

// New returns a new instance of Notify. It returns a new Notify instance with default options. By default, the Notify
// instance is enabled.
func New() *Notify {
	return NewWithOptions()
}

// Create the package level Notify instance.
var std = New()

// Default returns the standard Notify instance used by the package-level send function.
func Default() *Notify {
	return std
}
