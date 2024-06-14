package bdsaas

import (
	"context"
	"os"
	"testing"
	"time"
)

func Test(t *testing.T) {
	client := NewClient(os.Getenv("APP_KEY"))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	seats, err := client.GetSeats(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(seats) == 0 {
		t.Fatal("no seats")
	}
	t.Log("Seats:", seats)
	records, err := client.Query(ctx, "badu1801461545239609344")
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 {
		t.Fatal("unexpected records")
	}
	record := records[0]
	t.Logf("Record: %+v", record)
}
