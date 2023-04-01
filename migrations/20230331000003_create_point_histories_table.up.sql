CREATE TABLE point_histories (
    id INT IDENTITY(1,1) PRIMARY KEY,
    user_id INT NOT NULL,
    point INT NOT NULL,
    action NVARCHAR(255) NOT NULL,
    details NVARCHAR(255) NOT NULL,
    created_at DATETIME2 NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
