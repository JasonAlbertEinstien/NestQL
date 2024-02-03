# NestQL

## Introduction
NestQL is a relationship database that is written in golang.

## Folder Structure
Current folder structure 
```
.
└── src/
    ├── cmd/
    │   └── test_files
    ├── internal/
    │   ├── config
    │   ├── ds
    │   ├── execution
    │   ├── meta
    │   ├── prepare
    │   ├── query
    │   ├── server
    │   ├── storage
    │   ├── table
    │   └── utils
    └── main.go
```

## Contribution

## To-Do
Implementation of B+ tree. 
Create Customize Table
## Execute
```
go build -o NestQL *.go
./NestQL nestdb.nestQL
```
Then open another terminal
```
cd cmd
python3 client_test.py
Enter the command: 
```
For listing out all data
```
select
```
For inserting data
```
insert (id) (username) (email)
```



