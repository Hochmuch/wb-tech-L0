package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
	"wb-tech-L0/internal/model"
)

type Repository struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, pool *pgxpool.Pool) *Repository {
	repository := &Repository{pool}
	if err := repository.init(ctx); err != nil {
		panic("Couldn't init repository")
	}
	return repository
}

func (r *Repository) init(ctx context.Context) error {
	_, err := r.Pool.Exec(ctx,
		`CREATE EXTENSION IF NOT EXISTS "pgcrypto";
		
		CREATE TABLE IF NOT EXISTS orders (
    	order_uid UUID PRIMARY KEY,
    	track_number TEXT NOT NULL,
        entry TEXT,
        locale TEXT,
		internal_signature TEXT,
		customer_id TEXT,
		delivery_service TEXT,
		shardkey TEXT,
		sm_id INT,
		date_created TIMESTAMP,
		oof_shard TEXT
        );

		CREATE TABLE IF NOT EXISTS deliveries (
		order_uid UUID PRIMARY KEY REFERENCES orders(order_uid),
		name TEXT,
		phone TEXT,
		zip TEXT,
		city TEXT,
		address TEXT,
		region TEXT,
		email TEXT
		);

		CREATE TABLE IF NOT EXISTS payments (
		order_uid UUID PRIMARY KEY REFERENCES orders(order_uid),
		transaction TEXT,
		request_id TEXT,
		currency TEXT,
		provider TEXT,
		amount INT,
		payment_dt TIMESTAMP,
		bank TEXT,
		delivery_cost INT,
		goods_total INT,
		custom_fee INT
		);

		CREATE TABLE IF NOT EXISTS items (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		order_uid UUID REFERENCES orders(order_uid),
		chrt_id BIGINT,
		track_number TEXT,
		price INT,
		rid TEXT,
		name TEXT,
		sale INT,
		size TEXT,
		total_price INT,
		nm_id BIGINT,
		brand TEXT,
		status INT
		);
	`)
	return err
}

func (r *Repository) Save(ctx context.Context, order *model.Order) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id,
                    delivery_service, shardkey, sm_id, date_created, oof_shard)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO payments (order_uid, transaction, request_id, currency, provider,
                      amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		time.Unix(int64(order.Payment.PaymentDT), 0),
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err = tx.Exec(ctx,
			`INSERT INTO items (id, order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price,
             nm_id, brand, status)
			 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			order.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status)
		if err != nil {
			return err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetOrder(ctx context.Context, orderUID uuid.UUID) (*model.Order, error) {
	var order model.Order
	row := r.Pool.QueryRow(ctx, `SELECT
	order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service,
	shardkey, sm_id, date_created, oof_shard
	FROM orders WHERE order_uid = $1`, orderUID)

	err := row.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
		&order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated,
		&order.OofShard)

	if err != nil {
		fmt.Println("НЕ НАЙДЕН order", err)
		return nil, err
	}

	var delivery model.Delivery
	row = r.Pool.QueryRow(ctx, `SELECT
	name, phone, zip, city, address, region, email
	FROM deliveries WHERE order_uid = $1`, orderUID)

	err = row.Scan(&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City,
		&delivery.Address, &delivery.Region, &delivery.Email)

	order.Delivery = delivery

	if err != nil {
		fmt.Println("НЕ НАЙДЕН delivery", err)
		return nil, err
	}

	var payment model.Payment
	row = r.Pool.QueryRow(ctx, `SELECT
	transaction, request_id, currency, provider, amount,
	payment_dt, bank, delivery_cost, goods_total, custom_fee
	FROM payments WHERE order_uid = $1`, orderUID)

	var tmpPaymentDT time.Time

	err = row.Scan(&payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider, &payment.Amount,
		&tmpPaymentDT, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee)

	payment.PaymentDT = tmpPaymentDT.Unix()

	if err != nil {
		fmt.Println("НЕ НАЙДЕН payment", err)
		return nil, err
	}

	order.Payment = payment

	items := make([]model.Item, 0)

	rows, err := r.Pool.Query(ctx, `SELECT
	chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
	FROM items
	WHERE order_uid = $1`, order.OrderUID)
	if err != nil {
		fmt.Println("НЕ НАЙДЕН items", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Item

		if err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	order.Items = items
	return &order, nil
}

func (r *Repository) GetDataForCache(ctx context.Context) ([]*model.Order, error) {
	rows, err := r.Pool.Query(ctx, `
	SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service,
	shardkey, sm_id, date_created, oof_shard 
	FROM orders
	ORDER BY date_created DESC
	LIMIT 100;`)

	if err != nil {
		log.Println("Couldn't get cache to restore")
		return nil, err
	}
	defer rows.Close()

	result := make([]*model.Order, 0, 100)

	for rows.Next() {
		var order model.Order

		if err := rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
			&order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated,
			&order.OofShard); err != nil {
			log.Println("Couldn't save restored cache")
			return nil, err
		}

		result = append(result, &order)
	}
	return result, nil
}
