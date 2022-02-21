package run

import (
	"bytes"
	"log"
	"math"
	"reflect"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	pref "google.golang.org/protobuf/reflect/protoreflect"
)

// LogRepeatChanges returns a StreamServerInterceptor that logs out when streams return the same message more than
// 10 times in a row. The interceptor is aware of the patterns used in Smart Core and ignores differences in the
// PullResponse.Change.change_time property.
func LogRepeatChanges() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, &logIdenticalResponses{ServerStream: ss})
	}
}

type logIdenticalResponses struct {
	grpc.ServerStream
	lastMsg interface{}
	count   int
}

func (l *logIdenticalResponses) SendMsg(m interface{}) error {
	if l.changeEqual(l.lastMsg, m) {
		l.count++
		if l.count == 10 {
			json, err := protojson.Marshal(m.(proto.Message))
			if err != nil {
				log.Printf("WARN: 10+ repeated changes %v", m)
			} else {
				log.Printf("WARN: 10+ repeated changes %s", json)
			}
		}
	} else {
		l.count = 0
	}
	l.lastMsg = m
	return l.ServerStream.SendMsg(m)
}

func (l *logIdenticalResponses) changeEqual(a, b interface{}) bool {
	// There are some common patterns in smart core Pull methods, they all look like this:
	// Response { Changes: []Change{ Name, ChangeTime, Value } } for resources and
	// Response { Changes: []Change{ Name, ChangeTime, OldValue, NewValue, ChangeType } } for lists.
	//
	// This pattern breaks our equality checks because the ChangeTime value is generally different each time, even if
	// the value properties haven't changed. This is the error we're trying to avoid... churn

	if a == nil || b == nil {
		return a == nil && b == nil
	}
	am, ok := a.(proto.Message)
	if !ok {
		return false
	}
	bm, ok := b.(proto.Message)
	if !ok {
		return false
	}

	// this code is heavily borrowed from proto.Equal, but adjusted so we can ignore PullResponse.Change.change_time
	mx, my := am.ProtoReflect(), bm.ProtoReflect()
	if mx.IsValid() != my.IsValid() {
		return false
	}
	return equalMessage(mx, my)
}

// equalMessage compares two messages.
func equalMessage(mx, my pref.Message) bool {
	if mx.Descriptor() != my.Descriptor() {
		return false
	}

	nx := 0
	equal := true
	mx.Range(func(fd pref.FieldDescriptor, vx pref.Value) bool {
		nx++
		vy := my.Get(fd)
		equal = my.Has(fd) && equalField(fd, vx, vy)
		return equal
	})
	if !equal {
		return false
	}
	ny := 0
	my.Range(func(fd pref.FieldDescriptor, vx pref.Value) bool {
		ny++
		return true
	})
	if nx != ny {
		return false
	}

	return equalUnknown(mx.GetUnknown(), my.GetUnknown())
}

// equalField compares two fields.
func equalField(fd pref.FieldDescriptor, x, y pref.Value) bool {
	switch {
	// This is the case we've added, ignore PullResponse.Change.change_time
	case fd.Name() == "change_time" && fd.ContainingMessage().Name() == "Change":
		return true
	case fd.IsList():
		return equalList(fd, x.List(), y.List())
	case fd.IsMap():
		return equalMap(fd, x.Map(), y.Map())
	default:
		return equalValue(fd, x, y)
	}
}

// equalMap compares two maps.
func equalMap(fd pref.FieldDescriptor, x, y pref.Map) bool {
	if x.Len() != y.Len() {
		return false
	}
	equal := true
	x.Range(func(k pref.MapKey, vx pref.Value) bool {
		vy := y.Get(k)
		equal = y.Has(k) && equalValue(fd.MapValue(), vx, vy)
		return equal
	})
	return equal
}

// equalList compares two lists.
func equalList(fd pref.FieldDescriptor, x, y pref.List) bool {
	if x.Len() != y.Len() {
		return false
	}
	for i := x.Len() - 1; i >= 0; i-- {
		if !equalValue(fd, x.Get(i), y.Get(i)) {
			return false
		}
	}
	return true
}

// equalValue compares two singular values.
func equalValue(fd pref.FieldDescriptor, x, y pref.Value) bool {
	switch fd.Kind() {
	case pref.BoolKind:
		return x.Bool() == y.Bool()
	case pref.EnumKind:
		return x.Enum() == y.Enum()
	case pref.Int32Kind, pref.Sint32Kind,
		pref.Int64Kind, pref.Sint64Kind,
		pref.Sfixed32Kind, pref.Sfixed64Kind:
		return x.Int() == y.Int()
	case pref.Uint32Kind, pref.Uint64Kind,
		pref.Fixed32Kind, pref.Fixed64Kind:
		return x.Uint() == y.Uint()
	case pref.FloatKind, pref.DoubleKind:
		fx := x.Float()
		fy := y.Float()
		if math.IsNaN(fx) || math.IsNaN(fy) {
			return math.IsNaN(fx) && math.IsNaN(fy)
		}
		return fx == fy
	case pref.StringKind:
		return x.String() == y.String()
	case pref.BytesKind:
		return bytes.Equal(x.Bytes(), y.Bytes())
	case pref.MessageKind, pref.GroupKind:
		return equalMessage(x.Message(), y.Message())
	default:
		return x.Interface() == y.Interface()
	}
}

// equalUnknown compares unknown fields by direct comparison on the raw bytes
// of each individual field number.
func equalUnknown(x, y pref.RawFields) bool {
	if len(x) != len(y) {
		return false
	}
	if bytes.Equal([]byte(x), []byte(y)) {
		return true
	}

	mx := make(map[pref.FieldNumber]pref.RawFields)
	my := make(map[pref.FieldNumber]pref.RawFields)
	for len(x) > 0 {
		fnum, _, n := protowire.ConsumeField(x)
		mx[fnum] = append(mx[fnum], x[:n]...)
		x = x[n:]
	}
	for len(y) > 0 {
		fnum, _, n := protowire.ConsumeField(y)
		my[fnum] = append(my[fnum], y[:n]...)
		y = y[n:]
	}
	return reflect.DeepEqual(mx, my)
}
