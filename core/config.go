package core

import "flag"

type Config struct {
	ListenPort    string //ListenPort
	MySQLPoolSize int
	MySQLURI      string
	MySQLUser     string
	MySQLPswd     string
}

func ParseArgCmd() (*Config, error) {
	var cfg Config
	flag.StringVar(&cfg.ListenPort, "listen-port", ":3308", "mysql listen port")
	flag.StringVar(&cfg.MySQLUser, "mysql-usr", "root", "usr")
	flag.StringVar(&cfg.MySQLPswd, "mysql-psw", "rootkit", "psw")
	flag.StringVar(&cfg.MySQLURI, "mysql-uri", "127.0.0.1:3306", "pwd")
	flag.IntVar(&cfg.MySQLPoolSize, "pool-size", 32, ",poolsize")
	flag.Parse()

	return &cfg, nil
}
