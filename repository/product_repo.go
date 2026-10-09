package repository

import (
	"context"
	"database/sql"
	"errors"

	"bengkel-stainless/config"
	"bengkel-stainless/models"

	"github.com/go-sql-driver/mysql"
)

var ErrNotFound = errors.New("produk tidak ditemukan")

type ProductRepo struct{}

func NewProductRepo() *ProductRepo {
	return &ProductRepo{}
}

func (r *ProductRepo) KodeExists(ctx context.Context, kode string) (bool, error) {
	var count int
	err := config.DB.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM products WHERE kode = ?`, kode).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CREATE
func (r *ProductRepo) Create(ctx context.Context, p *models.Product) (int64, error) {
	query := `INSERT INTO products 
		(kode, nama, category_id, material_id, ketebalan, panjang, lebar, 
		 satuan, harga_beli, harga_jual, stok, stok_minimal)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`

	res, err := config.DB.ExecContext(ctx, query,
		p.Kode, p.Nama, p.CategoryID, p.MaterialID,
		p.Ketebalan, p.Panjang, p.Lebar, p.Satuan,
		p.HargaBeli, p.HargaJual, p.Stok, p.StokMinimal,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// READ - semua (TANPA deleted_at)
func (r *ProductRepo) GetAll(ctx context.Context, search string, limit, offset int) ([]models.Product, int, error) {
	like := "%" + search + "%"

	countQuery := `SELECT COUNT(*) FROM products WHERE nama LIKE ? OR kode LIKE ?`
	var total int
	if err := config.DB.QueryRowContext(ctx, countQuery, like, like).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `SELECT id, kode, nama, category_id, material_id, ketebalan, panjang, lebar,
	          satuan, harga_beli, harga_jual, stok, stok_minimal, created_at, updated_at
	          FROM products
	          WHERE nama LIKE ? OR kode LIKE ?
	          ORDER BY id DESC LIMIT ? OFFSET ?`

	rows, err := config.DB.QueryContext(ctx, query, like, like, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.ID, &p.Kode, &p.Nama, &p.CategoryID, &p.MaterialID,
			&p.Ketebalan, &p.Panjang, &p.Lebar, &p.Satuan,
			&p.HargaBeli, &p.HargaJual, &p.Stok, &p.StokMinimal,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		products = append(products, p)
	}
	return products, total, nil
}

// READ - by ID (TANPA deleted_at)
func (r *ProductRepo) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	query := `SELECT id, kode, nama, category_id, material_id, ketebalan, panjang, lebar,
	          satuan, harga_beli, harga_jual, stok, stok_minimal, created_at, updated_at
	          FROM products WHERE id = ?`

	var p models.Product
	err := config.DB.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Kode, &p.Nama, &p.CategoryID, &p.MaterialID,
		&p.Ketebalan, &p.Panjang, &p.Lebar, &p.Satuan,
		&p.HargaBeli, &p.HargaJual, &p.Stok, &p.StokMinimal,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// UPDATE (TANPA deleted_at)
func (r *ProductRepo) Update(ctx context.Context, id int64, u *models.ProductUpdate) error {
	query := `UPDATE products SET 
		kode = COALESCE(?, kode),
		nama = COALESCE(?, nama),
		category_id = COALESCE(?, category_id),
		material_id = COALESCE(?, material_id),
		ketebalan = COALESCE(?, ketebalan),
		panjang = COALESCE(?, panjang),
		lebar = COALESCE(?, lebar),
		satuan = COALESCE(?, satuan),
		harga_beli = COALESCE(?, harga_beli),
		harga_jual = COALESCE(?, harga_jual),
		stok = COALESCE(?, stok),
		stok_minimal = COALESCE(?, stok_minimal)
		WHERE id = ?`

	res, err := config.DB.ExecContext(ctx, query,
		u.Kode, u.Nama, u.CategoryID, u.MaterialID,
		u.Ketebalan, u.Panjang, u.Lebar, u.Satuan,
		u.HargaBeli, u.HargaJual, u.Stok, u.StokMinimal, id,
	)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// DELETE (hard delete + handle FK constraint)
func (r *ProductRepo) Delete(ctx context.Context, id int64) error {
	res, err := config.DB.ExecContext(ctx, `DELETE FROM products WHERE id = ?`, id)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1451 {
			return errors.New("produk tidak bisa dihapus karena masih dipakai di order")
		}
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
