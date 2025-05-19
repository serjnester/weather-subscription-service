package enums

type Frequency string

const (
	FrequencyDaily  Frequency = "daily"
	FrequencyHourly Frequency = "hourly"
)

func (f Frequency) String() string {
	return string(f)
}
