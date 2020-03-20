package purchase

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Service struct {
}

type Purchase struct {
	ID          int64  `json:"id"`
	User_id     int64  `json:"user_id"`
	Product_id  int64  `json:"product_id"`
	Price       int    `json:"price"`
	Kol         int 	`json:"kol"`
	Date 		string 	`json:"date"`
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) AddNewProduct(prod Purchase, pool *pgxpool.Pool) (err error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), `INSERT INTO products(name, description, price, pic)
VALUES ($1, $2, $3, $4);`, prod.Name, prod.Description, prod.Price, prod.Pic)
	if err != nil {
		return
	}
	return nil
}

func (s *Service) ProductList(pool *pgxpool.Pool) (list []Purchase, err error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(),
		`select id, name, description, price, pic from products where removed=false;`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := Purchase{}
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.Pic)
		if err != nil {
			return nil, errors.New("can't scan row from rows")
		}
		list = append(list, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.New("rows error!")
	}
	return
}

func (s *Service) RemoveByID(id int64, pool *pgxpool.Pool) (err error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return errors.New("can't connect to database!")
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), `update products set removed = true where id = $1`, id)
	if err != nil {
		return errors.New(fmt.Sprintf("can't remove from database purchase (id: %d)!", id))
	}
	return nil
}

func (s *Service) ProductByID(id int64, pool *pgxpool.Pool) (prod Purchase, err error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return Purchase{}, errors.New("can't connect to database!")
	}
	defer conn.Release()
	err = conn.QueryRow(context.Background(), `select id, name, description, price, pic from products where id=$1`,
		id).Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.Pic)
	if err != nil {
		return Purchase{}, errors.New(fmt.Sprintf("can't remove from database burger (id: %d)!", id))
	}
	return
}
