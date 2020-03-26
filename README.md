![ex_screenshot](./img/MariaDB-Tablespec-logo.png)
# mariadb-tablespec
This is Table Specification for MariaDB.


## Version
### 1.0
+ Build Up & Test.

### 1.1
+ Add Update Rule and Delete Rule for FK.
+ Add Glide (Dependency Management)

## Use Package 
```
go get github.com/360EntSecGroup-Skylar/excelize
go get github.com/go-sql-driver/mysql
```

## Arg
```
go run .\MariaDB-TableSpec.go -h
Usage of 
  -database string
        Export for Database.
  -file string
        Export File Path and Name.
  -host string
        Target Database Host.
  -password string
        Database Connect Password.
  -port string
        Target Database Port. (default "3306")
  -user string
        Database Connect User. (default "root")
```

## Example Usage
```
MariaDB-TableSpec.exe -host="127.0.0.1" -port="3306" -user="root" -password="root" -database="mysql" -file="mysql_table.xlsx"
```
or
```
go run MariaDB-TableSpec.go -host="127.0.0.1" -port="3306" -user="root" -password="root" -database="mysql" -file="mysql_table.xlsx"
```
> If you ' go run ' or build , you have to do dependency Package Install

## Output File
![ex_output](./img/ex_Output.png)