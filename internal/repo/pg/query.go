package pg

const (
	getScooterQueryKey         = "GetScooter"
	updateScooterQueryKey      = "UpdateScooter"
	findScootersInAreaQueryKey = "FindScootersInArea"
	storeEventQueryKey         = "StoreEvent"
)

var query = map[string]string{
	getScooterQueryKey:    `SELECT * FROM scooters WHERE id = $1`,
	updateScooterQueryKey: `UPDATE scooters SET status = :status, lat = :lat, lng = :lng, updated_at = :updated_at WHERE id = :id`,
	findScootersInAreaQueryKey: `
SELECT id, status, lat, lng, updated_at
FROM scooters
WHERE status = :status
  AND ST_Within(
    ST_SetSRID(ST_MakePoint(lng, lat), 4326),
    ST_MakeEnvelope(:min_lng, :min_lat, :max_lng, :max_lat, 4326)
  )
`,
	storeEventQueryKey: `INSERT INTO events (id, scooter_id, type, timestamp, lat, lng) VALUES (:id, :scooter_id, :type, :timestamp, :lat, :lng)`,
}
