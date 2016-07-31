package nsone

import "fmt"

// ZoneSecondaryServer wraps elements of a Zone's "primary.secondary" attribute
type ZoneSecondaryServer struct {
	Ip     string `json:"ip"`
	Port   int    `json:"port,omitempty"`
	Notify bool   `json:"notify"`
}

// ZonePrimary wraps a Zone's "primary" attribute
type ZonePrimary struct {
	Enabled     bool                  `json:"enabled"`
	Secondaries []ZoneSecondaryServer `json:"secondaries"`
}

// ZoneSecondary wraps a Zone's "secondary" attribute
type ZoneSecondary struct {
	Status       string `json:"status,omitempty"`
	Last_xfr     int    `json:"last_xfr,omitempty"`
	Primary_ip   string `json:"primary_ip,omitempty"`
	Primary_port int    `json:"primary_port,omitempty"`
	Enabled      bool   `json:"enabled"`
	Expired      bool   `json:"expired,omitempty"`
}

// ZoneRecord wraps Zone's "records" attribute
type ZoneRecord struct {
	Domain   string   `json:"Domain,omitempty"`
	Id       string   `json:"id,omitempty"`
	Link     string   `json:"link,omitempty"`
	ShortAns []string `json:"short_answers,omitempty"`
	Tier     int      `json:"tier,omitempty"`
	Ttl      int      `json:"ttl,omitempty"`
	Type     string   `json:"type,omitempty"`
}

// Zone wraps an NS1 /zone resource
type Zone struct {
	Id            string            `json:"id,omitempty"`
	Ttl           int               `json:"ttl,omitempty"`
	Nx_ttl        int               `json:"nx_ttl,omitempty"`
	Retry         int               `json:"retry,omitempty"`
	Zone          string            `json:"zone,omitempty"`
	Refresh       int               `json:"refresh,omitempty"`
	Expiry        int               `json:"expiry,omitempty"`
	Primary       *ZonePrimary      `json:"primary,omitempty"`
	Dns_servers   []string          `json:"dns_servers,omitempty"`
	Networks      []int             `json:"networks,omitempty"`
	Network_pools []string          `json:"network_pools,omitempty"`
	Hostmaster    string            `json:"hostmaster,omitempty"`
	Pool          string            `json:"pool,omitempty"`
	Meta          map[string]string `json:"meta,omitempty"`
	Secondary     *ZoneSecondary    `json:"secondary,omitempty"`
	Link          string            `json:"link,omitempty"`
	Records       []ZoneRecord      `json:"records,omitempty"`
	Serial        int               `json:"serial,omitempty"`
}

// NewZone takes a zone domain name and creates a new primary *Zone
func NewZone(zone string) *Zone {
	z := Zone{
		Zone: zone,
	}
	z.MakePrimary()
	return &z
}

// MakePrimary enables Primary, disables Secondary, and sets primary's Secondaries to all provided ZoneSecondaryServers
func (z *Zone) MakePrimary(secondaries ...ZoneSecondaryServer) {
	z.Secondary = nil
	z.Primary = &ZonePrimary{
		Enabled:     true,
		Secondaries: secondaries,
	}
	if z.Primary.Secondaries == nil {
		z.Primary.Secondaries = make([]ZoneSecondaryServer, 0)
	}
}

// MakeSecondary enables Secondary, disables Primary, and sets secondary's Primary_ip to provided ip
func (z *Zone) MakeSecondary(ip string) {
	z.Secondary = &ZoneSecondary{
		Enabled:      true,
		Primary_ip:   ip,
		Primary_port: 53,
	}
	z.Primary = &ZonePrimary{
		Enabled:     false,
		Secondaries: make([]ZoneSecondaryServer, 0),
	}
}

// LinkTo sets Link to a target zone domain name and unsets all other configuration properties
func (z *Zone) LinkTo(to string) {
	z.Meta = nil
	z.Ttl = 0
	z.Nx_ttl = 0
	z.Retry = 0
	z.Refresh = 0
	z.Expiry = 0
	z.Primary = nil
	z.Dns_servers = nil
	z.Networks = nil
	z.Network_pools = nil
	z.Hostmaster = ""
	z.Pool = ""
	z.Secondary = nil
	z.Link = to
}

// GetZones returns all active zones and basic zone configuration details for each
func (c APIClient) GetZones() ([]Zone, error) {
	var zl []Zone
	_, err := c.doHTTPUnmarshal("GET", "https://api.nsone.net/v1/zones", nil, &zl)
	return zl, err
}

// GetZone takes a zone and returns a single active zone and its basic configuration details
func (c APIClient) GetZone(zone string) (*Zone, error) {
	z := NewZone(zone)
	status, err := c.doHTTPUnmarshal("GET", fmt.Sprintf("https://api.nsone.net/v1/zones/%s", z.Zone), nil, z)
	if status == 404 {
		z.Id = ""
		z.Zone = ""
		return z, nil
	}
	return z, err
}

// DeleteZone takes a zone and destroys an existing DNS zone and all records in the zone
func (c APIClient) DeleteZone(zone string) error {
	return c.doHTTPDelete(fmt.Sprintf("https://api.nsone.net/v1/zones/%s", zone))
}

// CreateZone takes a *Zone and creates a new DNS zone
func (c APIClient) CreateZone(z *Zone) error {
	return c.doHTTPBoth("PUT", fmt.Sprintf("https://api.nsone.net/v1/zones/%s", z.Zone), z)
}

// UpdateZone takes a *Zone and modifies basic details of a DNS zone
func (c APIClient) UpdateZone(z *Zone) error {
	return c.doHTTPBoth("POST", fmt.Sprintf("https://api.nsone.net/v1/zones/%s", z.Zone), z)
}
