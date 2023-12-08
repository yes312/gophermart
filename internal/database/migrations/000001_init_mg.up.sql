CREATE TABLE IF NOT EXISTS users (		
    user_id VARCHAR PRIMARY KEY,
	hash VARCHAR NOT NULL
);
	
CREATE TABLE IF NOT EXISTS orders (
	number VARCHAR PRIMARY KEY,
	user_id VARCHAR NOT NULL,
	uploaded_at timestamp NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(user_id)
);
	
CREATE TABLE IF NOT EXISTS billing (
	order_number VARCHAR NOT NULL,
	status VARCHAR NOT NULL,
	accrual int, 
	uploaded_at timestamp NOT NULL,
	time timestamp NOT NULL,
	FOREIGN KEY (order_number) REFERENCES orders(number),
	CONSTRAINT unique_order_number_status UNIQUE (order_number, status)
);	