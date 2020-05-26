package postgres

import (
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
)

type MyTid struct {
	*pgtype.TID
}

func (dst *MyTid) Set(src interface{}) error {
	return errors.Errorf("cannot convert %v to TID", src)
}

func (dst MyTid) Get() interface{} {
	switch dst.Status {
	case pgtype.Present:
		return dst
	case pgtype.Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *MyTid) AssignTo(dst interface{}) error {
	if src.Status == pgtype.Present {
		switch v := dst.(type) {
		case *string:
			*v = fmt.Sprintf(`(%d,%d)`, src.BlockNumber, src.OffsetNumber)
			return nil
		default:
			if nextDst, retry := pgtype.GetAssignToDstType(dst); retry {
				return src.AssignTo(nextDst)
			}
			return errors.Errorf("unable to assign to %T", dst)
		}
	}

	return errors.Errorf("cannot assign %v to %T", src, dst)
}
