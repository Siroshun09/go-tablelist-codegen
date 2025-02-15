package internal

import "flag"

var Flag = struct {
	PackageName string
	DBConnInfo  struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
		SSLMode  string
	}
	Output string
	Debug  bool
}{}

func ParseFlags() {
	flag.StringVar(&Flag.PackageName, "package-name", "", "package name")
	flag.StringVar(&Flag.DBConnInfo.Host, "host", "", "host")
	flag.IntVar(&Flag.DBConnInfo.Port, "port", 0, "port")
	flag.StringVar(&Flag.DBConnInfo.User, "user", "", "user")
	flag.StringVar(&Flag.DBConnInfo.Password, "password", "", "password")
	flag.StringVar(&Flag.DBConnInfo.Database, "database", "", "database")
	flag.StringVar(&Flag.DBConnInfo.SSLMode, "sslmode", "disable", "sslmode")
	flag.StringVar(&Flag.Output, "output", "", "output")
	flag.BoolVar(&Flag.Debug, "debug", false, "debug")

	flag.Parse()
}
