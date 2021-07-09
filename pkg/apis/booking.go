package apis

import (
	"context"
	"time"

	"github.com/smart-core-os/sc-api/go/traits"
	scTime "github.com/smart-core-os/sc-api/go/types/time"
	"github.com/smart-core-os/sc-golang/pkg/memory"
	"github.com/smart-core-os/sc-golang/pkg/router"
	"github.com/smart-core-os/sc-golang/pkg/server"
	"github.com/smart-core-os/sc-golang/pkg/wrap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func BookingApi() server.GrpcApi {
	r := router.NewBookingApiRouter()
	r.Factory = func(name string) (traits.BookingApiClient, error) {
		return wrap.BookingApiServer(newBookingApiServer(name)), nil
	}
	return r
}

func newBookingApiServer(name string) *memory.BookingApi {
	api := memory.NewBookingApi()

	// we want the start of the day in local time, .Truncate(24*Hour) uses UTC
	now := time.Now().In(time.Local)
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// create a few bookings to make working with the data easier
	_, _ = api.CreateBooking(context.Background(), &traits.CreateBookingRequest{Name: name, Booking: &traits.Booking{
		Title:     "Test booking at 12:00",
		OwnerName: "Memo Ry",
		Booked: &scTime.Period{
			StartTime: timestamppb.New(startOfToday.Add(12 * time.Hour)),
			EndTime:   timestamppb.New(startOfToday.Add(13 * time.Hour)),
		},
	}})
	_, _ = api.CreateBooking(context.Background(), &traits.CreateBookingRequest{Name: name, Booking: &traits.Booking{
		Title:     "Test booking at 16:10",
		OwnerName: "Memo Ry",
		Booked: &scTime.Period{
			StartTime: timestamppb.New(startOfToday.Add(16*time.Hour + 10*time.Minute)),
			EndTime:   timestamppb.New(startOfToday.Add(17 * time.Hour)),
		},
	}})
	return api
}
