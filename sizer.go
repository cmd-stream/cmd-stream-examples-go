package exmpls

// Sizer defines a Size method.
//
// This interface is used by the client Codec to determine the command size.
type Sizer interface {
	Size() (size int)
}
