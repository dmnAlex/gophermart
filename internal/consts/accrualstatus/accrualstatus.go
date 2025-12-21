package accrualstatus

type Type string

const (
	New        Type = "REGISTERED"
	Processing Type = "PROCESSING"
	Invalid    Type = "INVALID"
	Processed  Type = "PROCESSED"
)

func (t Type) String() string {
	return string(t)
}
