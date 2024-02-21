package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {

	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	number, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(number), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {

	p := Parcel{}

	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number", sql.Named("number", number))
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

	return p, err
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {

	var res []Parcel

	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {

	_, err := s.db.Exec("UPDATE parcel SET status = :status where number = :number",
		sql.Named("status", status),
		sql.Named("number", number))

	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {

	var status string
	row := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number", sql.Named("number", number))
	row.Scan(&status)
	if status != ParcelStatusRegistered {
		return fmt.Errorf("unable to change address for the parcel with number: %d. status: %s. allowed status to change address is: %s", number, status, ParcelStatusRegistered)
	}
	_, err := s.db.Exec("UPDATE parcel SET address = :address where number = :number",
		sql.Named("address", address),
		sql.Named("number", number))

	return err
}

func (s ParcelStore) Delete(number int) error {

	var status string
	row := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number", sql.Named("number", number))
	err := row.Scan(&status)
	if err != nil {
		return err
	}
	if status != ParcelStatusRegistered {
		return fmt.Errorf("unable to delete parcel with number: %d. status: %s. allowed status to delete: %s", number, status, ParcelStatusRegistered)
	}
	_, err = s.db.Exec("DELETE FROM parcel where number = :number",
		sql.Named("number", number))

	return err
}
