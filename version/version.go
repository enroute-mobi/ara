package version

var value string

func Value() string {
	if value == "" {
		value = "20170529-122240" // quelque chose comme "20170529-122240"
	}
	return value
}
