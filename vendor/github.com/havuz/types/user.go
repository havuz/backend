package types

// User is the holder of user details used and exchanged
// between frontend and backend.
type User struct {
	ID string `csv:"UID"`

	Status,
	StatusReason string

	AllowedIPs      string `csv:"IPs"`
	SimultaneityCap int    `csv:"Slots"`
	Bandwidth       int

	CreatedAt,
	UpdatedAt,
	ExpiredAt string
}
