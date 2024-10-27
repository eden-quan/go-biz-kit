package uuid

import (
	"time"

	"github.com/rs/xid"
)

func New() string {
	return xid.New().String()
}

func NewWithTime(time time.Time) string {
	return xid.NewWithTime(time).String()
}

// Sort sorts ids in place
func Sort(ids []string) error {
	xIds := make([]xid.ID, 0, len(ids))
	for _, id := range ids {
		x, err := xid.FromString(id)
		if err != nil {
			return err
		}
		xIds = append(xIds, x)
	}

	xid.Sort(xIds)
	for i, x := range xIds {
		ids[i] = x.String()
	}

	return nil
}
