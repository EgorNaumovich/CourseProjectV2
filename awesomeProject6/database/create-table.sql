CREATE TABLE Deliveries (
	Id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	Container_Count INT NOT NULL,
	Weight INT NOT NULL,
	Description VARCHAR(1024),
	Transport VARCHAR(256),
	Customer VARCHAR(256),
	User VARCHAR(256),
	Origin VARCHAR(256)
)