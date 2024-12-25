package order

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Repository interface{
  Close()
  PutOrder(ctx context.Context, o Order) error
  GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type postgresRepository struct{
  db *sql.DB
}

func NewPostgresRepository(url string)(Repository, error){
  db, err := sql.Open("postgres", url)
  if err != nil{
    return nil, err
  }
  err = db.Ping()
  if err != nil{
    return nil, err
  }
  return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close(){
  r.db.Close()
}

func (r *postgresRepository) PutOrder(ctx context.Context, o Order)(err error) {
  tx, err := r.db.BeginTx(ctx, nil)
  if err != nil{
    return err
  }
  defer func(){
    if err != nil{
      tx.Rollback()
      return
    }
    err = tx.Commit()
  }()
  tx.ExecContext(ctx, "INSERT INTO orders(id, created_at, account_id, total_price) VALUES ($1, $2, $3, $4)",
    o.ID, o.CreatedAt, o.AccountID, o.TotalPrice)
  
  if err != nil{
    return 
  }

  stmt, _ := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))
  for _, p := range o.Products{
    _, err = stmt.ExecContext(ctx, o.ID, p.ID, p.Quantity)
    if err != nil{
      return
    } 
  }
  _, err = stmt.ExecContext(ctx)
  if err != nil{
    return
  }
  stmt.Close()
  return
}


func (r *postgresRepository) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
    // Execute query
    rows, err := r.db.QueryContext(ctx, `
        SELECT o.id, o.created_at, o.account_id, o.total_price::money::numeric::float8,
               op.product_id, op.quantity
        FROM orders o
        JOIN order_products op ON o.id = op.order_id
        WHERE o.account_id = $1
        ORDER BY o.id
    `, accountID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Initialize variables
    var orders []Order
    var lastOrder *Order
    var products []OrderedProduct

    for rows.Next() {
        var orderID, accountID string
        var createdAt time.Time
        var totalPrice float64
        var productID string
        var quantity uint32

        // Scan row data
        err := rows.Scan(&orderID, &createdAt, &accountID, &totalPrice, &productID, &quantity)
        if err != nil {
            return nil, err
        }

        // Check for a new order
        if lastOrder != nil && lastOrder.ID != orderID {
            lastOrder.Products = products
            orders = append(orders, *lastOrder)
            products = []OrderedProduct{}
        }

        // Create or update current order
        if lastOrder == nil || lastOrder.ID != orderID {
            lastOrder = &Order{
                ID:         orderID,
                AccountID:  accountID,
                CreatedAt:  createdAt,
                TotalPrice: totalPrice,
            }
        }

        // Append product to the current order
        products = append(products, OrderedProduct{
            ID:       productID,
            Quantity: uint32(quantity),
        })
    }

    // Finalize the last order
    if lastOrder != nil {
        lastOrder.Products = products
        orders = append(orders, *lastOrder)
    }

    // Check for errors during iteration
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return orders, nil
}
