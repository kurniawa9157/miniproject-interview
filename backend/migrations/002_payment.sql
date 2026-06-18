-- Migration: 002_payment.sql
-- Tambah kolom payment untuk integrasi Midtrans (bonus).
-- Jalankan: psql -U postgres -d miniproject -f migrations/002_payment.sql

ALTER TABLE orders
  ADD COLUMN IF NOT EXISTS amount          INTEGER NOT NULL DEFAULT 150000,
  ADD COLUMN IF NOT EXISTS payment_status  TEXT NOT NULL DEFAULT 'UNPAID',
  ADD COLUMN IF NOT EXISTS payment_token   TEXT;

-- payment_status: UNPAID | PENDING | PAID | FAILED
