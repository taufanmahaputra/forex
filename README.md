# Forex
[![Build Status](https://travis-ci.com/taufanmahaputra/forex.svg?branch=master)](https://travis-ci.com/taufanmahaputra/forex)  [![codecov](https://codecov.io/gh/taufanmahaputra/forex/branch/master/graph/badge.svg)](https://codecov.io/gh/taufanmahaputra/forex)  [![Go Report Card](https://goreportcard.com/badge/github.com/taufanmahaputra/forex)](https://goreportcard.com/report/github.com/taufanmahaputra/forex)  [![Maintainability](https://api.codeclimate.com/v1/badges/6c8ef9c28fe335f3c9d0/maintainability)](https://codeclimate.com/github/taufanmahaputra/forex/maintainability)

----
Forex is a simple API foreign exchange rate for daily basis written in Go with PostgreSQL and Docker.

| Documentation | Link |
| ------ | ------ |
| API | [http://bit.ly/forex-api-doc][api] |
| Database | [http://bit.ly/forex-db-doc][db] |

### Usage
----
#####  Run
```sh
$ git clone https://github.com/taufanmahaputra/forex.git
$ cd forex
$ docker-compose up
```
##### Test
```sh
$ go test -v ./test/...
```

License
----

MIT

   [api]: <http://bit.ly/forex-api-doc>
   [db]: <http://bit.ly/forex-db-doc>
   
