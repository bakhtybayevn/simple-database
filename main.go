package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

const Version = "1.0.0"

type Logger interface {
	Fatal(string, ...interface{})
	Error(string, ...interface{})
	Warn(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
}

type Driver struct {
	mu     sync.Mutex
	caches map[string]*sync.Mutex
	dir    string
	log    Logger
}

type Options struct {
	Logger
}

func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)
	opts := Options{}
	if options != nil {
		opts = *options
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger(lumber.INFO)
	}

	driver := &Driver{
		dir:    dir,
		log:    opts.Logger,
		caches: make(map[string]*sync.Mutex),
	}

	if _, err := os.Stat(dir); err != nil {
		opts.Logger.Debug("Using '%s' (database already exists)\n", dir)
		return driver, nil
	}

	opts.Logger.Debug("Creating the database at '%s'...\n", dir)
	return driver, os.MkdirAll(dir, 0755)
}

func (d *Driver) Write(collection, resource string, v interface{}) error {
	if collection == "" || resource == "" {
		return fmt.Errorf("collection and resource must not be empty when writing a record")
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, resource+".json")
	tmpPath := fnlPath + ".tmp"

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))

	if err := ioutil.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, fnlPath)
}

func (d *Driver) Read(collection, resource string, v interface{}) error {
	if collection == "" || resource == "" {
		return fmt.Errorf("collection and resource must not be empty when reading a record")
	}

	recordPath := filepath.Join(d.dir, collection, resource)
	if _, err := Stat(recordPath); err != nil {
		return err
	}

	b, err := ioutil.ReadFile(recordPath + ".json")
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &v)
}

func (d *Driver) ReadAll(collection string) ([]string, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must not be empty when reading all records")
	}

	dir := filepath.Join(d.dir, collection)
	if _, err := Stat(dir); err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	records := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		records = append(records, string(b))
	}

	return records, nil
}

func (d *Driver) Delete(collection, resource string) error {
	if collection == "" || resource == "" {
		return fmt.Errorf("collection and resource must not be empty when deleting a record")
	}

	path := filepath.Join(collection, resource)
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)

	switch fi, err := Stat(dir); {
	case fi == nil, err != nil:
		return fmt.Errorf("unable to delete '%s' because it does not exist", path)
	case fi.Mode().IsDir():
		return os.RemoveAll(dir)
	case fi.Mode().IsRegular():
		return os.RemoveAll(dir + ".json")
	}

	return nil
}

func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {
	d.mu.Lock()
	defer d.mu.Unlock()
	m, ok := d.caches[collection]
	if !ok {
		m = &sync.Mutex{}
		d.caches[collection] = m
	}
	return m
}

func Stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return
}

type Address struct {
	Street  string      `json:"street"`
	City    string      `json:"city"`
	Country string      `json:"country"`
	Pincode json.Number `json:"pincode"`
}

type User struct {
	Name    string      `json:"name"`
	Age     json.Number `json:"age"`
	Contact string      `json:"contact"`
	Company string      `json:"company"`
	Address Address     `json:"address"`
}

func main() {
	dir := "./"
	db, err := New(dir, nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	employees := []User{
		{"John", "30", "213", "ABC", Address{"Street 1", "City 1", "Country 1", "123456"}},
		{"Paul", "27", "213", "Facebook", Address{"Street 2", "City 2", "Country 2", "123456"}},
		{"Jessica", "22", "213", "Google", Address{"Street 3", "City 3", "Country 3", "123456"}},
		{"Akhil", "34", "213", "Meta", Address{"Street 4", "City 4", "Country 4", "123456"}},
		{"Alba", "42", "213", "Amazon", Address{"Street 5", "City 5", "Country 5", "123456"}},
		{"Stipe", "45", "213", "Yandex", Address{"Street 6", "City 6", "Country 6", "123456"}},
	}

	for _, employee := range employees {
		err := db.Write("users", employee.Name, &employee)
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	records, err := db.ReadAll("users")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("Records: ", records)

	allUsers := make([]User, 0)
	for _, record := range records {
		var user User
		err := json.Unmarshal([]byte(record), &user)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		allUsers = append(allUsers, user)
	}
	fmt.Println("All Users: ", allUsers)

	if err := db.Delete("users", "John"); err != nil {
		fmt.Println("Error: ", err)
	}
}
