# Latest benchmark results

Units: [ns/op (ratio from ORM)]

| Test | ORM | Raw | GORM |
| --- | --- | --- | --- |
| Insert | 20985 (100%) | 15152 (72%) | 59076 (281%) |
| Query | 2860340 (100%) | 2590434 (90%) | 8488300 (296%) |
| QueryLargeStruct | 16167653 (100%) | 17622053 (108%) | 75073092 (464%) |


![graph](./benchmark.png)



Benchmark time: 2018-01-27


#### Compared packages:

- [x] ORM: posener/orm (this package)
- [x] RAW: Direct SQL commands
- [x] GORM: jinzhu/gorm

#### Operations:

- [x] Insert: INSERT operations
- [X] Query: SELECT operations on an object with 2 fields
- [X] QueryLargeStruct: SELECT on an object of ~35 different fields

