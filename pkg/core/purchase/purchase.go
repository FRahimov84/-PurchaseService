package purchase

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"time"
)

type Service struct {
}

type Purchase struct {
	ID         int64  `json:"id"`
	User_id    int64  `json:"user_id"`
	Product_id int64  `json:"product_id"`
	Price      int    `json:"price"`
	Kol        int    `json:"kol"`
	Date       time.Time `json:"date"`
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) AddNewPurchase(prod Purchase, pool *pgxpool.Pool) (err error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), `INSERT INTO purchase(user_id, product_id, price, kol)
VALUES ($1, $2, $3, $4);`, prod.User_id, prod.Product_id, prod.Price, prod.Kol)
	if err != nil {
		return
	}
	return nil
}

func (s *Service) PurchaseList(pool *pgxpool.Pool) (list []Purchase, err error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(),
		`select id, user_id, product_id, price, kol, purchase_date from purchase where removed=false;`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := Purchase{}
		err := rows.Scan(&item.ID, &item.User_id, &item.Product_id, &item.Price, &item.Kol, &item.Date)
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
	_, err = conn.Exec(context.Background(), `update purchase set removed = true where id = $1`, id)
	if err != nil {
		return errors.New(fmt.Sprintf("can't remove from database purchase (id: %d)!", id))
	}
	return nil
}

func (s *Service) PurchaseByID(id int64, pool *pgxpool.Pool) (prod Purchase, err error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return Purchase{}, errors.New("can't connect to database!")
	}
	defer conn.Release()
	err = conn.QueryRow(context.Background(),
		`select id, user_id, product_id, price, kol, purchase_date from purchase where id=$1`,
		id).Scan(&prod.ID, &prod.User_id, &prod.Product_id, &prod.Price, &prod.Kol, &prod.Date)
	if err != nil {
		return Purchase{}, errors.New(fmt.Sprintf("can't remove from database burger (id: %d)!", id))
	}
	return
}

func (s *Service) PurchaseByUserID(id int64, pool *pgxpool.Pool) (list []Purchase, err error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		err =errors.New("can't connect to database!")
		return
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(),
		`select id, user_id, product_id, price, kol, purchase_date from purchase where user_id=$1`,
		id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := Purchase{}
		err := rows.Scan(&item.ID, &item.User_id, &item.Product_id, &item.Price, &item.Kol, &item.Date)
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
