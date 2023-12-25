package enums

const (
	daily = "daily"
	none  = "none"
)

type RetentionType string

func (e RetentionType) String() string {
	return string(e)
}

type retentionTypes struct{}

func (retentionTypes) Daily() RetentionType { return daily }
func (retentionTypes) None() RetentionType  { return none }

var RetentionTypes retentionTypes
