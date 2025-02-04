package repository

import (
	"database/sql"
	
	"github.com/VK-Container-Pinger/backend/models"
)

type DockerHostRepository struct {
	DB    *sql.DB
	stmts *dockerHostStatements
}

type dockerHostStatements struct {
	createTable *sql.Stmt
	insertHost  *sql.Stmt
	getHosts    *sql.Stmt
	deleteHost  *sql.Stmt
}

func NewDockerHostRepository(db *sql.DB) (*DockerHostRepository, error) {
	stmts, err := newDockerHostStatements(db)
	if err != nil {
		return nil, err
	}
	
	return &DockerHostRepository{
		DB:    db,
		stmts: stmts,
	}, nil
}

func newDockerHostStatements(db *sql.DB) (*dockerHostStatements, error) {
	var err error
	stmts := &dockerHostStatements{}

	stmts.createTable, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS docker_hosts (
			id SERIAL PRIMARY KEY,
			name TEXT,
			ip_address TEXT NOT NULL UNIQUE
		);
	`)
	if err != nil {
		return nil, err
	}

	stmts.insertHost, err = db.Prepare(`
		INSERT INTO docker_hosts (name, ip_address)
		VALUES ($1, $2) RETURNING id
	`)
	if err != nil {
		return nil, err
	}

	stmts.getHosts, err = db.Prepare(`
		SELECT id, name, ip_address FROM docker_hosts ORDER BY id
	`)
	if err != nil {
		return nil, err
	}

	stmts.deleteHost, err = db.Prepare(`
		DELETE FROM docker_hosts WHERE id = $1
	`)
	if err != nil {
		return nil, err
	}

	return stmts, nil
}

func (r *DockerHostRepository) InitTable() error {
	_, err := r.stmts.createTable.Exec()
	return err
}

func (r *DockerHostRepository) InsertHost(name, ip string) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	
	stmt := tx.Stmt(r.stmts.insertHost)
	var id int
	err = stmt.QueryRow(name, ip).Scan(&id)
	if err != nil {
		return 0, err
	}
	
	err = tx.Commit()
	return id, err
}


func (r *DockerHostRepository) GetHosts() ([]models.DockerHost, error) {
	rows, err := r.stmts.getHosts.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hosts []models.DockerHost
	for rows.Next() {
		var host models.DockerHost
		if err := rows.Scan(&host.ID, &host.Name, &host.IP); err != nil {
			continue
		}
		hosts = append(hosts, host)
	}
	return hosts, nil
}

func (r *DockerHostRepository) DeleteHost(id int) error {
	_, err := r.stmts.deleteHost.Exec(id)
	return err
}