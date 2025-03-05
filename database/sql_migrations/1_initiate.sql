-- +migrate Up


CREATE TYPE user_role AS ENUM ('doctor', 'patient', 'pharmacist');

CREATE TYPE prescription_status AS ENUM ('pending', 'approved', 'rejected');

CREATE TYPE transaction_status AS ENUM ('pending', 'completed', 'cancelled');

-- Tabel Users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role user_role NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel Categories
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

-- Tabel Medicines
CREATE TABLE medicines (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category_id INT NOT NULL,
    stock INT DEFAULT 0,
    price DECIMAL(10,2) NOT NULL,
    expiration_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

-- Tabel Prescriptions
CREATE TABLE prescriptions (
    id SERIAL PRIMARY KEY,
    doctor_id INT NOT NULL,
    patient_id INT NOT NULL,
    medicine_id INT NOT NULL,
    dosage VARCHAR(50) NOT NULL,
    instructions TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW,
    status prescription_status DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (doctor_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (patient_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (medicine_id) REFERENCES medicines(id) ON DELETE CASCADE
);

-- Tabel Transactions
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    patient_id INT NOT NULL,
    prescription_id INT NOT NULL,
    medicine_id INT NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    status transaction_status DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (patient_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (prescription_id) REFERENCES prescriptions(id) ON DELETE CASCADE,
    FOREIGN KEY (medicine_id) REFERENCES medicines(id) ON DELETE CASCADE
);


-- +migrate Down

DROP TABLE prescriptions;
DROP TABLE categories;
DROP TABLE medicines;
DROP TABLE users;
DROP TABLE transactions;
