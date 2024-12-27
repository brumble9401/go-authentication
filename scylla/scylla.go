package scylla

import (
	"github.com/brumble9401/golang-authentication/config"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
)

// Manager manages establishing connections to ScyllaDB.
type Manager struct {
	cfg config.ScyllaConfig
}

// NewManager creates a new Manager.
func NewManager(cfg config.ScyllaConfig) *Manager {
	return &Manager{
		cfg: cfg,
	}
}

// Connect connects to scylla and returns a session.
func (m *Manager) Connect() (gocqlx.Session, error) {
	return m.connect(m.cfg.ScyllaKeyspace, m.cfg.ScyllaHosts)
}

// CreateKeyspace creates a keyspace.
func (m *Manager) CreateKeyspace(keyspace string) error {
	session, err := m.connect("system", m.cfg.ScyllaHosts)
	if err != nil {
		return err
	}
	defer session.Close()

	stmt := `CREATE KEYSPACE IF NOT EXISTS authentication WITH replication = {'class': 'NetworkTopologyStrategy', 'replication_factor': 3}`
	return session.ExecStmt(stmt)
}

func (m *Manager) connect(keyspace string, hosts []string) (gocqlx.Session, error) {
	c := gocql.NewCluster(hosts...)
	c.Keyspace = keyspace
	c.Authenticator = gocql.PasswordAuthenticator{Username: "cassandra", Password: "cassandra"}
	return gocqlx.WrapSession(c.CreateSession())
}