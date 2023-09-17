# File-Based Database System in Go
This Go program implements a simple file-based database system for storing and managing JSON records. It provides CRUD (Create, Read, Update, Delete) operations on JSON data, ensuring data integrity and safety. This README will guide you through the code structure, features, and how to use the database system.
# Features
1. CRUD Operations: The database supports Create, Read, Update, and Delete operations for JSON records.
2. Concurrent Access: Safely handles concurrent access using mutex locks.
3. Error Handling: Provides robust error handling and meaningful error messages.
4. Logging: Utilizes a flexible logging system to record events and errors.
5. JSON Serialization: Serializes and deserializes data to/from JSON format.
6. Version Control: Maintains a version number to track changes.
# Prerequisites
Before using this database system, make sure you have the following prerequisites installed:
1. Go (Golang): [Installation Guide](https://go.dev/doc/install)
# Installation
1. Clone the repository to your local machine:
git clone https://github.com/bakhtybayevn/simple-database .git
2. Navigate to the project directory:
cd file-based-db-go
3. Run the Go program:
go run main.go
# Usage
# Initialize the Database
dir := "./"
db, err := New(dir, nil)
if err != nil {
    fmt.Println("Error: ", err)
}

1. dir is the path to the directory where the database will be stored.
2. You can provide a custom logger by passing an instance of the Logger interface to the New function.
# Write a Record
employee := User{
    Name:    "John",
    Age:     "30",
    Contact: "213",
    Company: "ABC",
    Address: Address{"Street 1", "City 1", "Country 1", "123456"},
}

err := db.Write("users", employee.Name, &employee)
if err != nil {
    fmt.Println("Error: ", err)
}
1. Use the Write method to add a new record to the specified collection.
# Read a Record
var user User
err := db.Read("users", "John", &user)
if err != nil {
    fmt.Println("Error: ", err)
}

1. Use the Read method to retrieve a specific record from the collection.
# Read All Records
records, err := db.ReadAll("users")
if err != nil {
    fmt.Println("Error: ", err)
}
1. Use the ReadAll method to fetch all records from a collection.
# Delete a Record
err := db.Delete("users", "John")
if err != nil {
    fmt.Println("Error: ", err)
}

1. Use the Delete method to remove a specific record from the collection.
# Acknowledgments
1. Thanks to the Go community for providing helpful resources and libraries.
