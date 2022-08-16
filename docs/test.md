```text 
create database test; 
```

```text 
create table testdata(id UInt64, email String, data1 Int64, data2 Int64) engine=Memory

create table user(id UInt64) engine=Memory
```

```text
insert into testdata values ( 1, "CC11001100@qq.com", 1, 2)

insert into user(id) values (1) 

insert into test.user values set id = 1; 

insert into test.user values; 

insert into user FORMAT TabSeparated values (1)
```



