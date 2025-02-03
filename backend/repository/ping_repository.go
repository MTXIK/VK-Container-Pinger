package repository

import (
	"database/sql"
	"time"
	"github.com/VK-Container-Pinger/backend/models"
)

type preparedStatements struct {
	initTable      *sql.Stmt
	insertPing     *sql.Stmt
	getPingResults *sql.Stmt
	deleteOldPingResults *sql.Stmt
	deleteOldRecordsForIps *sql.Stmt
}

type PingRepository struct {
	DB *sql.DB
	stmts *preparedStatements
}

func NewPingRepository(db *sql.DB) (*PingRepository, error) {
	stmts, err := newPreparedStatements(db)
	if err != nil {
		return nil, err
	}
	
	return &PingRepository{
		DB: db,
		stmts: stmts,
	}, nil
}

func newPreparedStatements(db *sql.DB) (*preparedStatements, error) {
	var err error
	stmts := &preparedStatements{}

	stmts.initTable, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS pings (
			id SERIAL PRIMARY KEY,
			ip_address TEXT,
			ping_time INTEGER,
			last_success TIMESTAMP
		);
	`)
	if err != nil {
		return nil, err
	}

	stmts.insertPing, err = db.Prepare("INSERT INTO pings (ip_address, ping_time, last_success) VALUES ($1, $2, $3)")
	if err != nil {
		return nil, err
	}

	stmts.getPingResults, err = db.Prepare("SELECT ip_address, ping_time, last_success FROM pings ORDER BY id DESC LIMIT $1")
	if err != nil {
		return nil, err
	}
	
	stmts.deleteOldPingResults, err = db.Prepare("DELETE FROM pings WHERE last_success < $1")
	if err != nil {
		return nil, err
	}
	
	stmts.deleteOldRecordsForIps, err = db.Prepare(`
        DELETE FROM pings
        WHERE last_success < NOW() - INTERVAL '24 hours'
          AND ip_address IN (
            SELECT DISTINCT ip_address
            FROM pings
            WHERE last_success >= NOW() - INTERVAL '24 hours'
          )
    `)
    if err != nil {
        return nil, err
    }

	return stmts, nil
}

func (r *PingRepository) InitTable() error {
	_, err := r.stmts.initTable.Exec()
	return err
}

func (r *PingRepository) InsertPingResult(pr models.PingResult) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	stmt := tx.Stmt(r.stmts.insertPing)
	_, err = stmt.Exec(pr.IPAddress, pr.PingTime, pr.LastSuccess)
	if err != nil {
		return err
	}
	
	err = tx.Commit()
	return err
}

func (r *PingRepository) GetPingResults(limit int) ([]models.PingResult, error) {
	rows, err := r.stmts.getPingResults.Query(limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	results := make([]models.PingResult, 0)
	for rows.Next() {
		var pr models.PingResult
		if err := rows.Scan(&pr.IPAddress, &pr.PingTime, &pr.LastSuccess); err != nil {
			continue
		}
		results = append(results, pr)
	}
	
	return results, nil
}

func (r *PingRepository) DeleteOldPingResults(before time.Time) error {
	_, err := r.stmts.deleteOldPingResults.Exec(before)
	return err
}

func (r *PingRepository) DeleteOldRecordsForIps() error {
    _, err := r.stmts.deleteOldRecordsForIps.Exec()
    return err
}