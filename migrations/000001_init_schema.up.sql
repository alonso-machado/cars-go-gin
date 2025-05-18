-- Create extension for UUID if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create cars table
CREATE TABLE IF NOT EXISTS cars (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    brand VARCHAR(100) NOT NULL,
    manufacturing_value DECIMAL(15, 2) NOT NULL CHECK (manufacturing_value > 0 AND manufacturing_value < 15000000),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_cars_name ON cars(name) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_cars_brand ON cars(brand) WHERE deleted_at IS NULL;

-- Create function to update updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to update updated_at column
CREATE TRIGGER update_cars_updated_at
BEFORE UPDATE ON cars
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Insert sample data
INSERT INTO cars (name, brand, manufacturing_value, description) VALUES
('Model S', 'Tesla', 79990.00, 'Luxury electric sedan'),
('Model 3', 'Tesla', 46990.00, 'Compact electric sedan'),
('Model X', 'Tesla', 99990.00, 'Luxury electric SUV'),
('Model Y', 'Tesla', 53990.00, 'Compact electric SUV'),
('Mustang Mach-E', 'Ford', 43995.00, 'Electric SUV'),
('F-150 Lightning', 'Ford', 39974.00, 'Electric pickup truck'),
('Ioniq 5', 'Hyundai', 39450.00, 'Electric crossover'),
('EV6', 'Kia', 40990.00, 'Electric crossover'),
('ID.4', 'Volkswagen', 41190.00, 'Electric SUV'),
('iX', 'BMW', 84900.00, 'Luxury electric SUV')
ON CONFLICT (name) DO NOTHING;
