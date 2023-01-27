package database

import (
	"database/sql"
	"fmt"

	"github.com/leandrobraga/goexpert-desafio-cleanarch/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) List(page, limit int, sort string) ([]entity.Order, error) {
	if sort != "" && sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	queryString := "SELECT * FROM orders"
	offset := (limit * page) - limit

	if page != 0 && limit != 0 {
		queryString = fmt.Sprintf("SELECT * FROM orders ORDER BY final_price %s LIMIT %d OFFSET %d ", sort, limit, offset)
	}
	res, err := r.Db.Query(queryString)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var orders []entity.Order
	for res.Next() {
		var order entity.Order
		err := res.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)

	}
	return orders, nil

}
