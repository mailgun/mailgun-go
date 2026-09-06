package mailgun

// Version of current release
/*
TODO(vtopc): automate this:
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range info.Deps {
			if dep.Path == "github.com/mailgun/mailgun-go(/v5)" {
				v = dep.Version
				break
			}
		}
	}
*/
const Version = "5.0.0"
