-- Attendance System Database Migration
-- Updated: 03-10-2025

-- Create database
CREATE DATABASE IF NOT EXISTS attendance_system CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE attendance_system;

-- Departments table
CREATE TABLE IF NOT EXISTS departments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    max_clock_in TIME NOT NULL,
    max_clock_out TIME NOT NULL,
    late_tolerance INT DEFAULT 15 COMMENT 'Tolerance in minutes',
    early_leave_penalty INT DEFAULT 30 COMMENT 'Penalty threshold in minutes',
    status ENUM('active', 'inactive') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_department_status (status),
    INDEX idx_department_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Employees table
CREATE TABLE IF NOT EXISTS employees (
    id INT AUTO_INCREMENT PRIMARY KEY,
    employee_id VARCHAR(50) NOT NULL UNIQUE,
    department_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    address TEXT NOT NULL,
    position VARCHAR(100),
    status ENUM('active', 'inactive', 'suspended') DEFAULT 'active',
    join_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    INDEX idx_employee_department (department_id),
    INDEX idx_employee_status (status),
    INDEX idx_employee_id (employee_id),
    INDEX idx_employee_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Attendance table
CREATE TABLE IF NOT EXISTS attendances (
    id INT AUTO_INCREMENT PRIMARY KEY,
    attendance_id VARCHAR(100) NOT NULL UNIQUE,
    employee_id VARCHAR(50) NOT NULL,
    clock_in TIMESTAMP NOT NULL,
    clock_in_date DATE,
    clock_out TIMESTAMP NULL,
    work_hours DECIMAL(4,2) NULL COMMENT 'Work hours in decimal',
    status ENUM('present', 'late', 'half-day', 'absent') DEFAULT 'present',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (employee_id) REFERENCES employees(employee_id) ON DELETE CASCADE ON UPDATE CASCADE,
    INDEX idx_attendance_employee (employee_id),
    INDEX idx_attendance_clock_in (clock_in),
    INDEX idx_attendance_clock_out (clock_out),
    INDEX idx_attendance_date (clock_in_date),
    INDEX idx_attendance_status (status),
    INDEX idx_attendance_id (attendance_id),
    UNIQUE KEY unique_employee_clock_in (employee_id, clock_in_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Attendance history table 
CREATE TABLE IF NOT EXISTS attendance_histories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    employee_id VARCHAR(50) NOT NULL,
    attendance_id VARCHAR(100) NOT NULL,
    date_attendance TIMESTAMP NOT NULL,
    attendance_type TINYINT NOT NULL COMMENT '1: Clock In, 2: Clock Out, 3: Adjustment, 4: Correction',
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (employee_id) REFERENCES employees(employee_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (attendance_id) REFERENCES attendances(attendance_id) ON DELETE CASCADE ON UPDATE CASCADE,
    INDEX idx_history_employee (employee_id),
    INDEX idx_history_attendance (attendance_id),
    INDEX idx_history_date (date_attendance),
    INDEX idx_history_type (attendance_type),
    INDEX idx_history_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Updated Users table with setup token fields
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NULL, -- NULL initially until account setup
    role VARCHAR(50) DEFAULT 'employee',
    employee_id VARCHAR(50) NULL,
    is_active BOOLEAN DEFAULT FALSE, -- False until account setup complete
    setup_token VARCHAR(64) NULL UNIQUE, -- Secure token for account setup
    token_expires TIMESTAMP NULL, -- Token expiration time
    last_login TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (employee_id) REFERENCES employees(employee_id) ON DELETE SET NULL ON UPDATE CASCADE,
    INDEX idx_user_username (username),
    INDEX idx_user_email (email),
    INDEX idx_user_role (role),
    INDEX idx_user_active (is_active),
    INDEX idx_user_employee (employee_id),
    INDEX idx_user_setup_token (setup_token),
    INDEX idx_user_token_expires (token_expires)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert sample departments
INSERT INTO departments (name, description, max_clock_in, max_clock_out, late_tolerance, early_leave_penalty) VALUES
('IT Department', 'Information Technology Department responsible for software development and infrastructure', '08:30:00', '17:00:00', 15, 30),
('HR Department', 'Human Resources Department handling recruitment and employee relations', '09:00:00', '17:30:00', 10, 15),
('Finance Department', 'Finance and Accounting Department managing company finances', '08:00:00', '16:30:00', 5, 10),
('Marketing Department', 'Marketing and Sales Department handling promotions and sales', '08:30:00', '17:30:00', 20, 30),
('Operations Department', 'Operations Department managing daily business operations', '08:00:00', '17:00:00', 10, 20);

-- Insert sample employees
INSERT INTO employees (employee_id, department_id, name, phone, address, position, status, join_date) VALUES
('EMP001', 1, 'John Doe', '+1234567890', '123 Main Street, City A, State X', 'Senior Software Engineer', 'active', '2023-01-15'),
('EMP002', 2, 'Jane Smith', '+1234567891', '456 Oak Avenue, City B, State Y', 'HR Manager', 'active', '2023-02-01'),
('EMP003', 3, 'Bob Johnson', '+1234567892', '789 Pine Road, City C, State Z', 'Finance Analyst', 'active', '2023-01-20'),
('EMP004', 1, 'Alice Brown', '+1234567893', '321 Elm Street, City D, State W', 'DevOps Engineer', 'active', '2023-03-10'),
('EMP005', 4, 'Charlie Wilson', '+1234567894', '654 Maple Drive, City E, State V', 'Marketing Specialist', 'active', '2023-02-15'),
('EMP006', 1, 'David Lee', '+1234567895', '987 Cedar Lane, City F, State U', 'Frontend Developer', 'active', '2023-04-01'),
('EMP007', 2, 'Sarah Chen', '+1234567896', '654 Birch Street, City G, State T', 'Recruitment Specialist', 'active', '2023-03-15'),
('EMP008', 3, 'Mike Garcia', '+1234567897', '321 Spruce Avenue, City H, State S', 'Senior Accountant', 'active', '2023-02-28');

-- Insert sample attendance records
INSERT INTO attendances (attendance_id, employee_id, clock_in, clock_in_date, clock_out, work_hours, status, notes) VALUES
('ATT001', 'EMP001', DATE_SUB(NOW(), INTERVAL 8 HOUR), CURDATE(), DATE_SUB(NOW(), INTERVAL 1 HOUR), 7.0, 'present', 'Regular work day'),
('ATT002', 'EMP002', DATE_SUB(NOW(), INTERVAL 8 HOUR), CURDATE(), DATE_SUB(NOW(), INTERVAL 1 HOUR), 7.0, 'present', 'Regular work day'),
('ATT003', 'EMP003', DATE_SUB(NOW(), INTERVAL 8 HOUR), CURDATE(), DATE_SUB(NOW(), INTERVAL 1 HOUR), 7.0, 'present', 'Regular work day'),
('ATT004', 'EMP004', DATE_SUB(NOW(), INTERVAL 9 HOUR), CURDATE(), DATE_SUB(NOW(), INTERVAL 2 HOUR), 7.0, 'late', 'Late due to traffic'),
('ATT005', 'EMP005', DATE_SUB(NOW(), INTERVAL 8 HOUR), CURDATE(), NULL, NULL, 'present', 'Still working');

-- Insert sample attendance history records
INSERT INTO attendance_histories (employee_id, attendance_id, date_attendance, attendance_type, description) VALUES
('EMP001', 'ATT001', DATE_SUB(NOW(), INTERVAL 8 HOUR), 1, 'Clock In recorded'),
('EMP001', 'ATT001', DATE_SUB(NOW(), INTERVAL 1 HOUR), 2, 'Clock Out recorded'),
('EMP004', 'ATT004', DATE_SUB(NOW(), INTERVAL 9 HOUR), 1, 'Late Clock In');

-- Insert default admin user (password: admin123)
INSERT INTO users (username, email, password, role, is_active) VALUES
('admin', 'admin@company.com', '$2a$10$lr9BGyvFP2VjwICQX10mQuHk5FKyf14nXYpLRCqLz5Xq7qaK4Uh4G', 'admin', TRUE);

-- Insert user accounts for employees with setup tokens (NO PASSWORD initially)
INSERT INTO users (username, email, password, role, employee_id, is_active, setup_token, token_expires) VALUES
('john.doe', 'john.doe@company.com', NULL, 'employee', 'EMP001', FALSE, 'setup_token_emp001_abc123def456', DATE_ADD(NOW(), INTERVAL 7 DAY)),
('jane.smith', 'jane.smith@company.com', '$2a$10$nt3gEDP3zJAyVPfXOYRG2OQJoeNyGKjqinn33plGEajuha/bWgqf6', 'manager', 'EMP002', TRUE, NULL, NULL),
('bob.johnson', 'bob.johnson@company.com', NULL, 'employee', 'EMP003', FALSE, 'setup_token_emp003_ghi789jkl012', DATE_ADD(NOW(), INTERVAL 7 DAY)),
('alice.brown', 'alice.brown@company.com', '$2a$10$nt3gEDP3zJAyVPfXOYRG2OQJoeNyGKjqinn33plGEajuha/bWgqf6', 'employee', 'EMP004', TRUE, NULL, NULL),
('charlie.wilson', 'charlie.wilson@company.com', NULL, 'employee', 'EMP005', FALSE, 'setup_token_emp005_mno345pqr678', DATE_ADD(NOW(), INTERVAL 7 DAY)),
('david.lee', 'david.lee@company.com', NULL, 'employee', 'EMP006', FALSE, 'setup_token_emp006_stu901vwx234', DATE_ADD(NOW(), INTERVAL 7 DAY)),
('sarah.chen', 'sarah.chen@company.com', NULL, 'employee', 'EMP007', FALSE, 'setup_token_emp007_yza567bcd890', DATE_ADD(NOW(), INTERVAL 7 DAY)),
('mike.garcia', 'mike.garcia@company.com', NULL, 'employee', 'EMP008', FALSE, 'setup_token_emp008_efg123hij456', DATE_ADD(NOW(), INTERVAL 7 DAY));

-- Create views for common queries
CREATE OR REPLACE VIEW employee_attendance_summary AS
SELECT 
    e.employee_id,
    e.name,
    d.name as department_name,
    COUNT(a.id) as total_attendance,
    SUM(CASE WHEN a.status = 'late' THEN 1 ELSE 0 END) as late_count,
    AVG(a.work_hours) as avg_work_hours
FROM employees e
LEFT JOIN departments d ON e.department_id = d.id
LEFT JOIN attendances a ON e.employee_id = a.employee_id AND a.clock_in_date = CURDATE()
GROUP BY e.id, e.employee_id, e.name, d.name;

-- Create view for pending account setups
CREATE OR REPLACE VIEW pending_account_setups AS
SELECT 
    u.id as user_id,
    u.username,
    u.email,
    u.setup_token,
    u.token_expires,
    e.employee_id,
    e.name as employee_name,
    d.name as department_name
FROM users u
INNER JOIN employees e ON u.employee_id = e.employee_id
INNER JOIN departments d ON e.department_id = d.id
WHERE u.is_active = FALSE 
AND u.setup_token IS NOT NULL 
AND u.token_expires > NOW();

-- Display success message
SELECT 'Database migration completed successfully!' as message;