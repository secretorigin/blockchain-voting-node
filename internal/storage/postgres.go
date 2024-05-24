package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`

	id  string       `json:"-"`
	mut sync.RWMutex `json:"-"`
}

func (cfg *PostgresConfig) SetId(id string) {
	cfg.mut.Lock()
	defer cfg.mut.Unlock()
	cfg.id = id
}

func (cfg *PostgresConfig) GetId() string {
	return cfg.id
}

func (cfg *PostgresConfig) SetJson(raw []byte) error {
	var newCfg PostgresConfig

	err := json.Unmarshal(raw, &newCfg)
	if err != nil {
		return err
	}

	cfg.Set(&newCfg)

	return nil
}

func (cfg *PostgresConfig) Set(newCfg *PostgresConfig) {
	cfg.mut.Lock()
	defer cfg.mut.Unlock()
	cfg.Host = newCfg.GetHost()
	cfg.Port = newCfg.GetPort()
	cfg.User = newCfg.GetUser()
	cfg.Password = newCfg.GetPassword()
	cfg.Database = newCfg.GetDatabase()
}

func (cfg *PostgresConfig) ToConnectionString() string {
	cfg.mut.RLock()
	defer cfg.mut.RUnlock()
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)
}

func (cfg *PostgresConfig) GetHost() string {
	cfg.mut.RLock()
	defer cfg.mut.RUnlock()
	return cfg.Host
}

func (cfg *PostgresConfig) GetPort() int {
	cfg.mut.RLock()
	defer cfg.mut.RUnlock()
	return cfg.Port
}

func (cfg *PostgresConfig) GetUser() string {
	cfg.mut.RLock()
	defer cfg.mut.RUnlock()
	return cfg.User
}

func (cfg *PostgresConfig) GetPassword() string {
	cfg.mut.RLock()
	defer cfg.mut.RUnlock()
	return cfg.Password
}

func (cfg *PostgresConfig) GetDatabase() string {
	cfg.mut.RLock()
	defer cfg.mut.RUnlock()
	return cfg.Database
}

type Postgres struct {
	requests map[string]string
	path     string
	cfg      *PostgresConfig
	mut      sync.RWMutex
}

func New(cfg *PostgresConfig) (*Postgres, error) {
	postgres := &Postgres{
		requests: make(map[string]string),
		cfg:      cfg,
	}

	// check postgres
	db, err := postgres.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return postgres, nil
}

func NewWithPath(path string, cfg *PostgresConfig) (*Postgres, error) {
	postgres := &Postgres{
		requests: make(map[string]string),
		path:     path,
		cfg:      cfg,
	}

	// check postgres
	db, err := postgres.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return postgres, nil
}

func (db *Postgres) Prepare(request string) error {
	db.mut.Lock()
	defer db.mut.Unlock()

	dat, err := db.getFromFile(db.path, request)
	if err != nil {
		return err
	}
	db.requests[request] = string(dat)

	return nil
}

func (db *Postgres) PrepareMultiple(requests []string) error {
	db.mut.Lock()
	defer db.mut.Unlock()

	for _, request := range requests {
		dat, err := db.getFromFile(db.path, request)
		if err != nil {
			return err
		}
		db.requests[request] = string(dat)
	}

	return nil
}

func (db *Postgres) PrepareMultipleFromPath(path string, requests []string) error {
	db.mut.Lock()
	defer db.mut.Unlock()

	for _, request := range requests {
		dat, err := db.getFromFile(path, request)
		if err != nil {
			return err
		}
		db.requests[request] = string(dat)
	}

	return nil
}

func (db *Postgres) getFromFile(path string, request string) ([]byte, error) {
	return os.ReadFile(path + request + ".sql")
}

func (db *Postgres) Get(request string) string {
	db.mut.RLock()
	defer db.mut.RUnlock()

	val, ok := db.requests[request]
	if !ok {
		return ""
	}

	return val
}

func (db *Postgres) Open() (*sql.DB, error) {
	dbconn, err := sql.Open("postgres", db.cfg.ToConnectionString())
	return dbconn, err
}

// NullByteSlice represents a []byte that may be null.
type NullByteSlice struct {
	ByteSlice []byte
	Valid     bool
}

// Scan implements the Scanner interface for NullByteSlice.
func (n *NullByteSlice) Scan(value interface{}) error {
	if value == nil {
		n.ByteSlice, n.Valid = nil, false
		return nil
	}
	n.Valid = true
	n.ByteSlice = value.([]byte)
	return nil
}
