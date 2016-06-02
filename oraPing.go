package ora

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-oci8"
)

type OraConnParams struct {
	OraConnID string
	OraUser   string
	OraPasswd string
}

func PingOra(oraConn OraConnParams, interval time.Duration) error {

	connURl := fmt.Sprintf("%s:%s@%s", oraConn.OraUser, oraConn.OraPasswd, oraConn.OraConnID)
	oraDB, err := sql.Open("oci8", connURl)
	defer oraDB.Close()
	log.Printf("Connecting as %s", connURl)

	if err != nil {
		log.Fatal(err)
		return err
	}

	c := time.Tick(interval)
	for now := range c {
		if inst_name, host_name, version, err := getOracleInstanceInfor(oraDB); err != nil {
			return nil
		} else {
			fmt.Printf("%v:Oracle instance name: %s hostname:%s version:%s)\n", now, inst_name, host_name, version)

		}
	}
	return nil
}

func getOracleInstanceInfor(db *sql.DB) (string, string, string, error) {
	var inst_name, host_name, version string
	err := db.QueryRow("select instance_name,host_name,version from v$instance  ").Scan(&inst_name, &host_name, &version)
	switch {
	case err == sql.ErrNoRows:
		log.Fatal("Panic:'select instance_name from v$instance' returns no row")
		return "", "", "", nil
	case err != nil:
		log.Fatal(err)
		return "", "", "", nil
	default:
	}
	return inst_name, host_name, version, nil
}
