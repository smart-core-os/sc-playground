package th

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
)

// Diff is like cmp.Diff but with protocmp.Transform and cmpopts.EquateApprox included as options by default.
func Diff(x, y interface{}, opts ...cmp.Option) string {
	allOpts := []cmp.Option{
		protocmp.Transform(),
		cmpopts.EquateApprox(0, 0.00001),
	}
	allOpts = append(allOpts, opts...)
	return cmp.Diff(x, y, allOpts...)
}
