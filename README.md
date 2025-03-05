# MediTrack – Sistem Manajemen Resep & Stok Obat untuk Praktik Pribadi

## 📌 Deskripsi
MediTrack adalah solusi digital untuk dokter yang menjalankan praktik pribadi dalam mengelola resep elektronik, riwayat pemeriksaan pasien, serta stok opname obat dengan lebih efisien. Dengan MediTrack, dokter dapat dengan mudah mencatat, mengakses, dan memperbarui riwayat resep pasien, memastikan stok obat selalu tersedia.


## 🛠️ Fitur Utama
✅ **Autentikasi & Authorization**
- Dokter, Pasien, memiliki role yang berbeda
- Login & Register untuk semua pengguna
- JWT Authentication untuk akses API

✅ **Manajemen Resep Elektronik**
- Dokter bisa membuat resep elektronik untuk pasien
- Pasien bisa melihat daftar resep mereka
- Dokter/Apoteker bisa memverifikasi & memproses resep yang masuk

✅ **Manajemen Stok Obat**
- Dokter bisa menambahkan & mengedit stok obat
- Kategori obat (Antibiotik, Vitamin, dll.) untuk pengelolaan yang lebih rapi

✅ **History**
- Riwayat resep pasien
- Laporan stok obat

---

## 📊 Database Schema (ERD)
Schema database terdiri dari 6 tabel utama:
1. **Users** (id, name, email, password, role) → Dokter, Pasien, Apoteker
2. **Prescriptions** (id, doctor_id, patient_id, medicine_id, dosage, instructions, status)
3. **Medicines** (id, name, category_id, stock, price, expiration_date)
4. **Categories** (id, name) → Misal: Antibiotik, Vitamin
5. **Transactions** (id, patient_id, prescription_id, medicine_id, total_price, status)
6. **Payments** (id, transaction_id, amount, payment_status)

---

## 📌 API Documentation
Semua route berada di dalam grup `/api` dan dilindungi oleh **JWT Middleware**.

### 🔑 Authentication
| Method | Endpoint  | Deskripsi |
|--------|----------|-----------|
| POST   | `/api/register` | Register user |
| POST   | `/api/login` | Login user |

### 🏥 Manajemen Resep
| Method | Endpoint  | Deskripsi |
|--------|----------|-----------|
| POST   | `/api/prescription` | Dokter membuat resep |
| GET    | `/api/prescription/:patient_id` | Pasien melihat resep mereka |
| PUT    | `/api/prescription/:prescription_id/status` | Dokter mengupdate status resep |

### 🏷️ Manajemen Kategori Obat
| Method | Endpoint  | Deskripsi |
|--------|----------|-----------|
| POST   | `/api/categories` | Menambahkan kategori obat |
| GET    | `/api/categories` | Melihat semua kategori |

### 💊 Manajemen Obat
| Method | Endpoint  | Deskripsi |
|--------|----------|-----------|
| POST   | `/api/medicines` | Apotek menambahkan obat |
| GET    | `/api/medicines` | Melihat semua obat |
| PUT    | `/api/medicines/:id` | Mengupdate stok obat |
| DELETE | `/api/medicines/:id` | Menghapus obat |

### 🛒 Transaksi & Pembayaran
| Method | Endpoint  | Deskripsi |
|--------|----------|-----------|
| POST   | `/api/transactions` | Pasien membeli obat dari resep |
| POST   | `/api/payments` | Pasien melakukan pembayaran |
| GET    | `/api/payments/:id` | Mengecek status pembayaran |

---



