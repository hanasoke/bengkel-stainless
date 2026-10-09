package repository

import (
	"context"
	"database/sql"
	"errors"

	"bengkel-stainless/config"
	"bengkel-stainless/models"
)

var ErrNotFound = errors.New("produk tidak ditemukan")

type ProductRepo struct{}

func NewProductRepo() *ProductRepo {
	return &ProductRepo{}
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

// READ - semua (dengan pagination & search)
func (r *ProductRepo) GetAll(ctx context.Context, search string, limit, offset int) ([]models.Product, int, error) {
	// Hitung total
	countQuery := `SELECT COUNT(*) FROM products 
		WHERE deleted_at IS NULL AND (nama LIKE ? OR kode LIKE ?)`
	var total int
	like := "%" + search + "%"
	if err := config.DB.QueryRowContext(ctx, countQuery, like, like).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Ambil data
	query := `SELECT id, kode, nama, category_id, material_id, ketebalan, panjang, lebar,
	          satuan, harga_beli, harga_jual, stok, stok_minimal, created_at, updated_at
	          FROM products
	          WHERE deleted_at IS NULL AND (nama LIKE ? OR kode LIKE ?)
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

// READ - by ID
func (r *ProductRepo) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	query := `SELECT id, kode, nama, category_id, material_id, ketebalan, panjang, lebar,
	          satuan, harga_beli, harga_jual, stok, stok_minimal, created_at, updated_at
	          FROM products WHERE id = ? AND deleted_at IS NULL`

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

// UPDATE
func (r *ProductRepo) Update(ctx context.Context, id int64, u *models.ProductUpdate) error {
	// Bangun query dinamis
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
		WHERE id = ? AND deleted_at IS NULL`

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

// DELETE (soft delete)
func (r *ProductRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE products SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	res, err := config.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
