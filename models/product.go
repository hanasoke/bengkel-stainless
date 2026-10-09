package models

import "time"

type Product struct {
	ID          int64      `json:"id"`
	Kode        string     `json:"kode" validate:"required,min=2,max=50"`
	Nama        string     `json:"nama" validate:"required,min=2,max=150"`
	CategoryID  *int64     `json:"category_id"`
	MaterialID  *int64     `json:"material_id"`
	Ketebalan   *float64   `json:"ketebalan"`
	Panjang     *float64   `json:"panjang"`
	Lebar       *float64   `json:"lebar"`
	Satuan      string     `json:"satuan" validate:"required,oneof=pcs batang meter kg lembar"`
	HargaBeli   float64    `json:"harga_beli" validate:"gte=0"`
	HargaJual   float64    `json:"harga_jual" validate:"gte=0"`
	Stok        int        `json:"stok" validate:"gte=0"`
	StokMinimal int        `json:"stok_minimal" validate:"gte=0"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// Untuk update parsial (opsional)
type ProductUpdate struct {
	Kode        *string  `json:"kode" validate:"omitempty,min=2,max=50"`
	Nama        *string  `json:"nama" validate:"omitempty,min=2,max=150"`
	CategoryID  *int64   `json:"category_id"`
	MaterialID  *int64   `json:"material_id"`
	Ketebalan   *float64 `json:"ketebalan"`
	Panjang     *float64 `json:"panjang"`
	Lebar       *float64 `json:"lebar"`
	Satuan      *string  `json:"satuan" validate:"omitempty,oneof=pcs batang meter kg lembar"`
	HargaBeli   *float64 `json:"harga_beli" validate:"omitempty,gte=0"`
	HargaJual   *float64 `json:"harga_jual" validate:"omitempty,gte=0"`
	Stok        *int     `json:"stok" validate:"omitempty,gte=0"`
	StokMinimal *int     `json:"stok_minimal" validate:"omitempty,gte=0"`
}
