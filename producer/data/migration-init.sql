CREATE TABLE IF NOT EXISTS companies (
  id BINARY(16) NOT NULL COMMENT 'uuid binary' ,
  name VARCHAR(15) NOT NULL,
  description VARCHAR(3000) NULL,
  amt_employees INT NOT NULL ,
  registered BOOLEAN NOT NULL ,
  `type` ENUM('Corporations','NonProfit','Cooperative','Sole Proprietorship') NOT NULL,
  PRIMARY KEY (id),
  UNIQUE (name)
) ENGINE = InnoDB;
