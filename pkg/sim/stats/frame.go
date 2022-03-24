package stats

import (
	"fmt"
	"io"
	"time"
)

type Framer interface {
	Frame(f FrameTiming)
}

type FrameTiming struct {
	Total   time.Duration
	Capture time.Duration
	Scrub   time.Duration
	Respond time.Duration
	Idle    time.Duration
}

func FFramer(out io.Writer) Framer {
	return fFramer{Writer: out}
}

type fFramer struct {
	io.Writer
}

func (fm fFramer) Frame(f FrameTiming) {
	fmt.Fprintf(fm.Writer, "%v (c:%v s:%v r:%v)\n", f.Total, f.Capture, f.Scrub, f.Respond)
}

func DFramer(d time.Duration, f Framer) Framer {
	return &dFramer{Framer: f, d: d}
}

type dFramer struct {
	Framer
	d, t time.Duration
	fts  []FrameTiming
}

func (f *dFramer) Frame(ft FrameTiming) {
	f.t += ft.Total
	f.fts = append(f.fts, ft)

	if f.t > f.d {
		av := average(f.fts)
		f.t = 0
		f.fts = f.fts[:]
		f.Framer.Frame(av)
	}
}
func CFramer(c int, f Framer) Framer {
	return &cFramer{Framer: f, c: c}
}

type cFramer struct {
	Framer
	c, t int
	fts  []FrameTiming
}

func (f *cFramer) Frame(ft FrameTiming) {
	f.t++
	f.fts = append(f.fts, ft)

	if f.t >= f.c {
		av := average(f.fts)
		f.t = 0
		f.fts = f.fts[:]
		f.Framer.Frame(av)
	}
}

func average(fts []FrameTiming) FrameTiming {
	av := FrameTiming{}
	l := time.Duration(len(fts))
	for _, t := range fts {
		av.Total += t.Total
		av.Capture += t.Capture
		av.Scrub += t.Scrub
		av.Respond += t.Respond
	}
	av.Total = av.Total / l
	av.Capture = av.Capture / l
	av.Scrub = av.Scrub / l
	av.Respond = av.Respond / l

	return av
}
