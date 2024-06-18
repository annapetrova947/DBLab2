
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

func main() {
    // Database connection details
    dsn := "user:password@tcp(127.0.0.1:3306)/"
    dbName := "Observatory"

    // Connect to the database
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create the database
    _, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
    if err != nil {
        log.Fatal(err)
    }

    // Use the created database
    _, err = db.Exec(fmt.Sprintf("USE %s", dbName))
    if err != nil {
        log.Fatal(err)
    }

    // Create tables
    createTables := []string{
        `CREATE TABLE IF NOT EXISTS Sector (
            id INT AUTO_INCREMENT PRIMARY KEY,
            coordinates VARCHAR(255),
            light_intensity FLOAT,
            foreign_objects INT,
            star_objects_count INT,
            unknown_objects_count INT,
            defined_objects_count INT,
            notes TEXT
        )`,
        `CREATE TABLE IF NOT EXISTS Objects (
            id INT AUTO_INCREMENT PRIMARY KEY,
            type VARCHAR(255),
            accuracy FLOAT,
            quantity INT,
            time TIME,
            date DATE,
            notes TEXT
        )`,
        `CREATE TABLE IF NOT EXISTS NaturalObjects (
            id INT AUTO_INCREMENT PRIMARY KEY,
            type VARCHAR(255),
            galaxy VARCHAR(255),
            accuracy FLOAT,
            light_flow FLOAT,
            associated_objects VARCHAR(255),
            notes TEXT
        )`,
        `CREATE TABLE IF NOT EXISTS Position (
            id INT AUTO_INCREMENT PRIMARY KEY,
            earth_position VARCHAR(255),
            sun_position VARCHAR(255),
            moon_position VARCHAR(255)
        )`,
        `CREATE TABLE IF NOT EXISTS Observation (
            id INT AUTO_INCREMENT PRIMARY KEY,
            sector_id INT,
            object_id INT,
            natural_object_id INT,
            position_id INT,
            FOREIGN KEY (sector_id) REFERENCES Sector(id),
            FOREIGN KEY (object_id) REFERENCES Objects(id),
            FOREIGN KEY (natural_object_id) REFERENCES NaturalObjects(id),
            FOREIGN KEY (position_id) REFERENCES Position(id)
        )`,
    }

    for _, query := range createTables {
        _, err = db.Exec(query)
        if err != nil {
            log.Fatal(err)
        }
    }

    // Create trigger
    createTrigger := `
    DELIMITER //
    CREATE TRIGGER UpdateObjects
    AFTER UPDATE ON Objects
    FOR EACH ROW
    BEGIN
        DECLARE column_exists INT DEFAULT 0;

        -- Check for the existence of the column date_update
        SELECT COUNT(*) INTO column_exists
        FROM INFORMATION_SCHEMA.COLUMNS
        WHERE TABLE_NAME = 'Objects' AND COLUMN_NAME = 'date_update';

        -- Add the column if it does not exist
        IF column_exists = 0 THEN
            SET @alter_sql = 'ALTER TABLE Objects ADD COLUMN date_update TIMESTAMP';
            PREPARE stmt FROM @alter_sql;
            EXECUTE stmt;
            DEALLOCATE PREPARE stmt;
        END IF;

        -- Update the date_update column with the current date and time
        SET NEW.date_update = NOW();
    END //
    DELIMITER ;
    `
    _, err = db.Exec(createTrigger)
    if err != nil {
        log.Fatal(err)
    }

    // Create procedure
    createProcedure := `
    DELIMITER //
    CREATE PROCEDURE JoinTables(IN table1 VARCHAR(255), IN table2 VARCHAR(255))
    BEGIN
        SET @sql = CONCAT('SELECT * FROM ', table1, ' t1 JOIN ', table2, ' t2 ON t1.id = t2.id');
        PREPARE stmt FROM @sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
    END //
    DELIMITER ;
    `
    _, err = db.Exec(createProcedure)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Database schema initialized successfully.")
}
