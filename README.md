
# fake-response

## About

Fake response for rest api.

http schema only.

## Requirements

- go 

#Installation

```bash
$ go get github.com/narita-takeru/fake-response/cmd/fake-response
```

#Usage

```bash
$ fake-response fake.tml
```

#Sample spec file

```
ports:
  src: 127.0.0.1:12345
  dst: 127.0.0.1:12346
extract: ^GET /1.0/(.*)\?
endpoints:
  users: {users: [{id: 1, name: "a"}, {id: 2, name: "b"}}
```


