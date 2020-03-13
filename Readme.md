Rabbitcli
=========

Simple cmd client to read and write from rabbit message server.

To read from rabbit use the -dial and -channel parameters as above:

```
rabbitreader -dial amqp://guest:guest@localhost:5672/ -channel test
```

To write to rabbit use the -dial and -channel again:

```
cat data | rabbitwriter -dial amqp://guest:guest@localhost:5672/ -channel test
```

