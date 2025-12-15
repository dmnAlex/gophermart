package orderstatus

type Type string

const (
	New        Type = "NEW"
	Processing Type = "PROCESSING"
	Invalid    Type = "INVALID"
	Processed  Type = "PROCESSED"
)

func (t Type) String() string {
	return string(t)
}
