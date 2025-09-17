package ses

var (
	serviceFile = ""
	name        = "space_envoy"
)

func SetServiceFile(n string) {
	serviceFile = n
}

func SetName(n string) {
	name = n
}
