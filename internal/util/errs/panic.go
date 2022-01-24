package errs

import "io"

func CloseOrPanic(closer io.Closer) {
	if err := closer.Close(); err != nil {
		panic(err)
	}
}

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}
