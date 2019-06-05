package hotel

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

type regionRepositoryInt interface {
	update(Regions) error
	get(dest string) (Region, error)
}

type regionRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) regionRepository {
	return regionRepository{
		db: db,
	}
}

func (repository regionRepository) update(regions Regions) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return errors.New(fmt.Sprintf("tx begin err %v", err))
	}
	_, err = tx.Exec(`delete from regions`)
	if err != nil {
		return errors.New(fmt.Sprintf("error deleting %v", err))
	}
	query := `insert into regions (id, name, data) values ($1, $2, $3)`

	for _, value := range regions {
		data, err := json.Marshal(value)
		_, err = tx.Exec(query, value.Id, value.Name, data)
		if err != nil {
			return errors.New(fmt.Sprintf(" insert error %v", err))
		}
	}
	err = tx.Commit()
	if err != nil {
		return errors.New(fmt.Sprintf("commit error %v", err))
	}
	return nil
}

func (repository regionRepository) get(dest string) (Region, error) {
	tx, err := repository.db.Begin()
	if err != nil {
		fmt.Println("tx begin error", err)
		return Region{}, err
	}
	var b []byte
	query := `select data from regions where name=$1`
	row := tx.QueryRow(query, dest)
	err = row.Scan(&b)
	if err != nil {
		fmt.Println("tx scan error ", err)
		return Region{}, err
	}
	var region Region
	err = json.Unmarshal(b, &region)
	if err != nil {
		fmt.Println("json unmarshal error ", err)
		return Region{}, err
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println("tx commit error", err)
		return Region{}, err
	}
	return region, nil
}
