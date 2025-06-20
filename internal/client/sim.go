package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/adrianpk/rida/internal/telemetry"
	"github.com/google/uuid"
)

const (
	APIHost      = "localhost"
	APIPort      = 8080
	ScootersPath = "/api/v1/scooters"

	FindScootersRetryDelay = 1 * time.Second
	NoScootersRestDelay    = 2 * time.Second
	NoValidIDsRestDelay    = 2 * time.Second
	PreRideDelay           = 200 * time.Millisecond
	TripMinDuration        = 10 // seconds
	TripDurationJitter     = 6  // seconds
	UpdateLocationInterval = 3 * time.Second
	RestMinDuration        = 2 // seconds
	RestDurationJitter     = 4 // seconds
	LatJitter              = 0.01
	LngJitter              = 0.01
)

type Sim struct {
	ID     uuid.UUID
	APIKey string
	Client *http.Client
	Lat    float64
	Lng    float64
	Tag    string
}

func NewSim(apiKey string, lat, lng float64, tag string) *Sim {
	return &Sim{
		ID:     uuid.New(),
		APIKey: apiKey,
		Client: &http.Client{Timeout: 2 * time.Second},
		Lat:    lat,
		Lng:    lng,
		Tag:    tag,
	}
}

func (c *Sim) FindScooters(ctx context.Context) ([]telemetry.Scooter, error) {
	area := c.BoundingBox(400)
	status := telemetry.StatusFree

	baseURL := fmt.Sprintf("http://%s:%d%s", APIHost, APIPort, ScootersPath)

	params := url.Values{}
	params.Set("minLat", fmt.Sprintf("%f", area.MinLat))
	params.Set("maxLat", fmt.Sprintf("%f", area.MaxLat))
	params.Set("minLng", fmt.Sprintf("%f", area.MinLng))
	params.Set("maxLng", fmt.Sprintf("%f", area.MaxLng))
	params.Set("status", string(status))

	fullURL := baseURL + "?" + params.Encode()

	req, err := c.newReq(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var scooters []telemetry.Scooter
	if err := json.NewDecoder(resp.Body).Decode(&scooters); err != nil {
		return nil, err
	}

	return scooters, nil
}

func (c *Sim) StartRide(ctx context.Context, scooterID uuid.UUID) error {
	event := telemetry.Event{
		ScooterID: scooterID,
		Type:      telemetry.EventTripStart,
	}

	return c.sendEvent(ctx, event)
}

func (c *Sim) StopRide(ctx context.Context, scooterID uuid.UUID) error {
	event := telemetry.Event{
		ScooterID: scooterID,
		Type:      telemetry.EventTripEnd,
	}

	return c.sendEvent(ctx, event)
}

// UpdateLocation updates the location of the user.
// We assume that the scooter also reports its location
// and the backend eventually matches both locations as a security measure.
func (c *Sim) UpdateLocation(ctx context.Context, scooterID uuid.UUID, lat, lng float64) error {
	c.Lat = lat
	c.Lng = lng
	event := telemetry.Event{
		ScooterID: scooterID,
		Type:      telemetry.EventLocation,
		Lat:       lat,
		Lng:       lng,
	}

	return c.sendEvent(ctx, event)
}

// sendEvent posts an event to the backend /api/v1/events endpoint.
func (c *Sim) sendEvent(ctx context.Context, event telemetry.Event) error {
	url := fmt.Sprintf("http://%s:%d/api/v1/events", APIHost, APIPort)

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	req, err := c.newReq(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("event post failed: %s", resp.Status)
	}

	return nil
}

func (c *Sim) CurrentLocation() (float64, float64) {
	return c.Lat, c.Lng
}

// BoundingBox returns a bounding box (telemetry.Area) around the current Sim
// location for a given distance in meters.
func (c *Sim) BoundingBox(distanceMeters float64) telemetry.Area {
	return BoundingBox(c.Lat, c.Lng, distanceMeters)
}

// MoveAfterFail moves the client a larger random distance from the current
// position, simulating a bigger move in the city.
func (c *Sim) MoveAfterFail() {
	// Move: up to 0.01 degrees in any direction (~1km)
	dLat := (rand.Float64()*2 - 1) * 0.01
	dLng := (rand.Float64()*2 - 1) * 0.01
	c.Lat += dLat
	c.Lng += dLng
	c.logf("stroll to (%.5f, %.5f)", c.Lat, c.Lng)
}

func (c *Sim) Run(ctx context.Context) error {
	c.logf("initial position: (%.5f, %.5f)", c.Lat, c.Lng)
	for {
		scooters, err := c.FindScooters(ctx)
		if err != nil {
			c.logf("error: %v", err)
			time.Sleep(FindScootersRetryDelay)
			continue
		}

		if len(scooters) == 0 {
			c.logf("no scooters found")
			c.MoveAfterFail()
			time.Sleep(NoScootersRestDelay)
			continue
		}

		scooterID, ok := pickNearest(scooters)
		if !ok {
			c.logf("no valid scooter IDs found")
			time.Sleep(NoValidIDsRestDelay)
			continue
		}

		time.Sleep(PreRideDelay)

		if err := c.StartRide(ctx, scooterID); err != nil {
			c.logf("error starting ride: %v", err)
			continue
		}

		c.logf("start ride")

		tripDuration := time.Duration(TripMinDuration+rand.Intn(TripDurationJitter)) * time.Second
		tripStart := time.Now()
		for time.Since(tripStart) < tripDuration {
			lat := c.Lat + (rand.Float64()*2-1)*LatJitter
			lng := c.Lng + (rand.Float64()*2-1)*LngJitter

			_ = c.UpdateLocation(ctx, scooterID, lat, lng)

			c.logf("update location (%.5f, %.5f)", lat, lng)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(UpdateLocationInterval):
			}
		}

		if err := c.StopRide(ctx, scooterID); err != nil {
			c.logf("error stopping ride: %v", err)
		} else {
			c.logf("stop ride")
		}

		restDuration := time.Duration(RestMinDuration+rand.Intn(RestDurationJitter)) * time.Second
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(restDuration):
		}
	}
}

// TagID returns a more friendly ID for the Sim instance.
func (c *Sim) TagID() string {
	idStr := strings.ToLower(c.ID.String())
	if len(idStr) >= 8 {
		idStr = idStr[len(idStr)-8:]
	}

	return fmt.Sprintf("%s-%s", strings.ToLower(c.Tag), idStr)
}

func (c *Sim) newReq(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	if c.APIKey != "" {
		req.Header.Set("X-API-Key", c.APIKey)
	}

	req.Header.Set("X-Client-ID", c.TagID())
	return req, nil
}

func (c *Sim) logf(format string, args ...interface{}) {
	log.Printf("[%s] %s", c.TagID(), fmt.Sprintf(format, args...))
}

// pickNearest returns the first scooter ID from the slice, assuming the slice
// is ordered by distance (nearest first).
func pickNearest(scooters []telemetry.Scooter) (uuid.UUID, bool) {
	if len(scooters) == 0 {
		return uuid.Nil, false
	}

	return scooters[0].ID, true
}
