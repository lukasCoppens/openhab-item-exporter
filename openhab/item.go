package openhab

var (
	itemMapping = map[string]int{
		"OFF":    0,
		"CLOSED": 0,
		"ON":     1,
		"OPEN":   1,
	}
)

type Item struct {
	State      string
	Type       string
	Name       string
	Tags       []string
	GroupNames []string
}

func (it Item) GetIntState() int {
	if value, ok := itemMapping[it.State]; ok {
		return value
	}
	return 99
}
