CREATE TABLE orders (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  symbol VARCHAR(20),
  side ENUM('buy', 'sell'),
  type ENUM('limit', 'market'),
  price DECIMAL(10,2),
  quantity INT,
  remaining_quantity INT,
  status ENUM('open', 'filled', 'canceled'),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE trades (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  buy_order_id BIGINT,
  sell_order_id BIGINT,
  symbol VARCHAR(20),
  price DECIMAL(10,2),
  quantity INT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (buy_order_id) REFERENCES orders(id),
  FOREIGN KEY (sell_order_id) REFERENCES orders(id)
);
