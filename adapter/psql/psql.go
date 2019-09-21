package psql

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-pg/pg"
	"github.com/mayur-tolexo/aero/conf"
	"github.com/mayur-tolexo/flaw"
)

var (
	//StartLogging is set to log query
	StartLogging     bool
	isDebuggerActive bool
	debuggerStatus   map[*pg.DB]bool
	dbPostgresWrite  map[string]*pg.DB
	dbPostgresRead   map[string][]*pg.DB
	masterContainer  = "database.master"
	slaveContainer   = "database.slaves"
)

func init() {
	debuggerStatus = make(map[*pg.DB]bool)
	dbPostgresWrite = make(map[string]*pg.DB)
	dbPostgresRead = make(map[string][]*pg.DB)
}

//init master connection
func initMaster() (err error) {
	if conf.Exists(masterContainer) {

		if dbPostgresWrite[masterContainer] == nil {
			var postgresWriteOption pg.Options
			if postgresWriteOption, err = getPostgresOptions(masterContainer); err == nil {
				dbPostgresWrite[masterContainer] = pg.Connect(&postgresWriteOption)
			}
		}
	} else {
		err = errors.New("Master config does not exists")
	}
	return
}

//init slave connections
func initSlaves() (err error) {
	if conf.Exists(slaveContainer) {
		slaves := conf.StringSlice([]string{}, slaveContainer)

		if dbPostgresRead[slaveContainer] == nil {
			dbPostgresRead[slaveContainer] = make([]*pg.DB, len(slaves))
		}
		for i, container := range slaves {
			if dbPostgresRead[slaveContainer][i] == nil {
				var postgresReadOption pg.Options
				if postgresReadOption, err = getPostgresOptions(container); err != nil {
					break
				}
				dbPostgresRead[slaveContainer][i] = pg.Connect(&postgresReadOption)
			}
		}
	} else {
		err = errors.New("Slaves config does not exists")
	}
	return
}

//CreateMaster will create new master connection
func CreateMaster() (err error) {
	if conf.Exists(masterContainer) {
		var postgresWriteOption pg.Options
		if postgresWriteOption, err = getPostgresOptions(masterContainer); err == nil {
			dbPostgresWrite[masterContainer] = pg.Connect(&postgresWriteOption)
		}
	} else {
		err = errors.New("Master config does not exists")
	}
	return
}

//CreateSlave will create new slave connections
func CreateSlave() (err error) {
	if conf.Exists(slaveContainer) {
		slaves := conf.StringSlice([]string{}, slaveContainer)

		if dbPostgresRead[slaveContainer] == nil {
			dbPostgresRead[slaveContainer] = make([]*pg.DB, len(slaves))
		}
		for i, container := range slaves {

			var postgresReadOption pg.Options
			if postgresReadOption, err = getPostgresOptions(container); err != nil {
				break
			}
			dbPostgresRead[slaveContainer][i] = pg.Connect(&postgresReadOption)
		}
	} else {
		err = errors.New("Slaves config does not exists")
	}
	return
}

//set postgres connection options from conf
func getPostgresOptions(container string) (pgOption pg.Options, err error) {
	if !conf.Exists(container) {
		err = errors.New("Container for postgres configuration not found")
		return
	}
	host := conf.String("", container+".host")
	port := conf.String("", container+".port")
	addr := ""
	if host != "" && port != "" {
		addr = host + ":" + port
	}
	pgOption.Addr = addr
	pgOption.User = conf.String("", container+".username")
	pgOption.Password = conf.String("", container+".password")
	pgOption.Database = conf.String("", container+".db")
	pgOption.MaxRetries = conf.Int(3, container+".maxRetries")
	pgOption.RetryStatementTimeout = conf.Bool(false, container+".retryStmTimeout")
	return
}

//Conn will return postgres connection
func Conn(writable bool) (dbConn *pg.DB, err error) {
	rand.Seed(time.Now().UnixNano())
	if writable {
		if err = initMaster(); err == nil {
			dbConn = dbPostgresWrite[masterContainer]
		}
	} else {
		if err = initSlaves(); err == nil {
			if len(dbPostgresRead) == 0 {
				if err = initMaster(); err == nil {
					dbConn = dbPostgresWrite[masterContainer]
				}
			} else {
				dbConn = dbPostgresRead[slaveContainer][rand.Intn(len(dbPostgresRead[slaveContainer]))]
			}
		}
	}
	if err == nil {
		if debuggerStatus[dbConn] == false {
			debuggerStatus[dbConn] = true
			logQuery(dbConn)
		}
	} else {
		err = flaw.ConnError(err)
	}
	return
}

//Tx will return postgres transaction
func Tx() (tx *pg.Tx, err error) {
	var conn *pg.DB
	if conn, err = Conn(true); err == nil {
		if tx, err = conn.Begin(); err != nil {
			err = flaw.TxError(err)
		}
	}
	return
}

//ConnByContainer will return postgres connection by container
func ConnByContainer(container string) (*pg.DB, error) {
	if strings.HasSuffix(container, "master") == true {
		oldContainer := masterContainer
		masterContainer = container
		conn, err := Conn(true)
		masterContainer = oldContainer
		return conn, err
	} else if strings.HasSuffix(container, "slaves") == true {
		oldContainer := slaveContainer
		slaveContainer = container
		conn, err := Conn(false)
		slaveContainer = oldContainer
		return conn, err
	}
	return nil, errors.New("No master or slaves container found in: " + container)
}

//logQuery : Print postgresql query on terminal
func logQuery(conn *pg.DB) {
	conn.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
		if StartLogging == true {
			if query, err := event.FormattedQuery(); err == nil {
				var queryError string
				if event.Error != nil {
					queryError = "\nQUERY ERROR: " + event.Error.Error()
				}
				fmt.Println("----DEBUGGER----")
				fmt.Printf("\nFile: %v : %v\nFunction: %v\nQuery Execution Taken: %s\n%s%s\n\n",
					event.File, event.Line, event.Func, time.Since(event.StartTime), query, queryError)
			} else {
				fmt.Println("Debugger Error: " + err.Error())
			}
		}
	})
}

//Debug : Print postgresql query on terminal
func Debug(conn *pg.DB) {
	if isDebuggerActive == false {
		isDebuggerActive = true
		conn.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
			if query, err := event.FormattedQuery(); err == nil {
				var queryError string
				if event.Error != nil {
					queryError = "\nQUERY ERROR: " + event.Error.Error()
				}
				fmt.Println("----DEBUGGER----")
				fmt.Printf("\nFile: %v : %v\nFunction: %v\nQuery Execution Taken: %s\n%s%s\n\n",
					event.File, event.Line, event.Func, time.Since(event.StartTime), query, queryError)
			} else {
				fmt.Println("Debugger Error: " + err.Error())
			}
		})
	}
}
