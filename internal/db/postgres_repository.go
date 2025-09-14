package db

import (
	"context"
	"log"
	"time"
	"wb-tech-L0/domain/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, dbURL string) (*PostgresRepository, error) {
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Printf("Unable to parse databaseURL: %v\n", err)
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Printf("Unable to create connection pool: %v\n", err)
		return nil, err
	}

	return &PostgresRepository{pool}, nil
}

func (r *PostgresRepository) SaveOrder(ctx context.Context, order *model.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

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
		time.Unix(order.Payment.PaymentDT, 0),
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

func (r *PostgresRepository) GetOrder(ctx context.Context, orderUID uuid.UUID) (*model.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var order model.Order

	row := tx.QueryRow(ctx,
		`SELECT orders.order_uid, orders.track_number, orders.entry, orders.locale, orders.internal_signature,
    	orders.customer_id, orders.delivery_service, orders.shardkey, orders.sm_id, orders.date_created, orders.oof_shard,
    	deliveries.name, deliveries.phone, deliveries.zip, deliveries.city, deliveries.address, deliveries.region,
    	deliveries.email,
    	payments.transaction, payments.request_id, payments.currency, payments.provider, payments.amount,
    	payments.payment_dt, payments.bank, payments.delivery_cost, payments.goods_total, payments.custom_fee
		FROM orders JOIN deliveries ON orders.order_uid = deliveries.order_uid
    	JOIN payments ON orders.order_uid = payments.order_uid
		WHERE orders.order_uid = $1`,
		orderUID,
	)

	var paymentDTTimestamp time.Time

	err = row.Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID,
		&order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard,
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address,
		&order.Delivery.Region, &order.Delivery.Email,
		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider,
		&order.Payment.Amount, &paymentDTTimestamp, &order.Payment.Bank, &order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal, &order.Payment.CustomFee,
	)

	if err != nil {
		return nil, err
	}

	order.Payment.PaymentDT = paymentDTTimestamp.Unix()

	var rows pgx.Rows
	rows, err = tx.Query(ctx,
		`SELECT items.id, items.chrt_id, items.track_number, items.price, items.rid, items.name, items.sale,
		items.size, items.total_price, items.nm_id, items.brand, items.status
		FROM items WHERE items.order_uid = $1`, orderUID,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var item model.Item
		err = rows.Scan(
			&item.ID, &item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *PostgresRepository) GetTodayOrdersUIDs(ctx context.Context) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0)
	rows, err := r.pool.Query(ctx,
		`SELECT order_uid
			FROM orders
			WHERE date_created::date = CURRENT_DATE`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *PostgresRepository) GetDataForCache(ctx context.Context) ([]*model.Order, error) {
	ids, err := r.GetTodayOrdersUIDs(ctx)
	if err != nil {
		return nil, err
	}

	orders := make([]*model.Order, 0)
	for _, id := range ids {
		var order *model.Order
		order, err = r.GetOrder(ctx, id)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *PostgresRepository) Close() {
	if r == nil || r.pool == nil {
		return
	}
	r.pool.Close()
}
