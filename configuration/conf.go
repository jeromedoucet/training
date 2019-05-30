package configuration

import "fmt"

// GlobalConf store db connection informations
type GlobalConf struct {
	DbName   string
	User     string
	Password string
	Host     string
	Port     uint
}

// DbStringConnection return the string connection that will be used
// by postgres driver
func (g GlobalConf) DbStringConnection() string {
	var strConn string

	if len(g.DbName) > 0 {
		strConn = fmt.Sprintf("dbname=%s", g.DbName)
	}

	if len(g.User) > 0 {
		strConn = fmt.Sprintf("%s user=%s", strConn, g.User)
	}

	if len(g.Password) > 0 {
		strConn = fmt.Sprintf("%s password=%s", strConn, g.Password)
	}

	if len(g.Host) > 0 {
		strConn = fmt.Sprintf("%s host=%s", strConn, g.Host)
	}

	if g.Port != 0 {
		strConn = fmt.Sprintf("%s port=%d", strConn, g.Port)
	}

	return strConn
}
