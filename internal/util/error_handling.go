package util

import "io"

func CloseOrPanic(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		panic(err)
	}
}

func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}
