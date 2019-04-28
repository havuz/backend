package types

// Tunnel is an edge tunnel returned by backend.
type Tunnel struct {
	AppID string `json:"app_id"`

	SSHHost string `json:"tunnel_host"`
	SSHUser string `json:"client_user"`
	SSHKey  string `json:"client_key"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	Dyno string `json:"dyno_name"`
	User string `json:"dyno_user"`
	IP   string `json:"dyno_ip"`

	Region string `json:"region"`
}

// ReferenceID returns a logical ID that refers
// to the tunnel by concatenating AppID and Dyno.
func (t Tunnel) ReferenceID() string {
	return t.AppID + "-" + t.Dyno
}
