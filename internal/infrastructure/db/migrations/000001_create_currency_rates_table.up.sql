CREATE TABLE currency_rates (
    id SERIAL PRIMARY KEY,
    pair VARCHAR(10) NOT NULL,
    ask_price DECIMAL(10, 2) NOT NULL,
    bid_price DECIMAL(10, 2) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL
);