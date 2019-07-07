OctoSQL
=======
OctoSQL is a data query tool, allowing you to join, analyze and transform data from multiple data sources and file formats using SQL.

## Table of Contents
- [What is OctoSQL?](#what-is-octosql)
- [Quickstart](#quickstart)
- [Installation](#installation)
- [Configuration](#configuration)
- [Documentation](#documentation)
- Supported Databases (with an optimal functionality matrix)
- [Architecture](#architecture)
- Some examples of query diagrams
- [Roadmap](#roadmap)

## What is OctoSQL?
OctoSQL is a SQL query engine which allows you to write standard SQL queries on data in multiple SQL databases, NoSQL databases and files in various formats trying to push down as much of the work as possible to the source databases, not transferring unnecessary data. 

OctoSQL does that by creating an internal representation of your query and later translating parts of it into the query languages or APIs of the source databases. Whenever a datasource doesn't support a given operation, OctoSQL will execute it in memory, so you don't have to worry about the specifics of the underlying datasources. 

With OctoSQL you don't need O(n) client tools or a large data analysis system deployment. Everything's contained in a single binary.

## Quickstart
Let's say we have a csv file with cats, and a redis database with people (potential cat owners). Now we want to get a list of cats with the cities their owners live in.

First, create a configuration file ([Configuration Syntax](#configuration))
For example:
```yaml
dataSources:
  - name: cats
    type: csv
    config:
      path: "~/Documents/cats.csv"
  - name: people
    type: redis
    config:
      address: "localhost:6379"
      password: ""
      databaseIndex: 0
      databaseKeyName: "id"
```

Then, set the **OCTOSQL_CONFIG** environment variable to point to the configuration file.
```bash
export OCTOSQL_CONFIG=~/octosql.yaml
```

Finally, query to your hearts desire:
```bash
octosql "SELECT c.name, c.livesleft, p.city
FROM cats c JOIN people p ON c.ownerid = p.id"
```
Example output:
```
+----------+-------------+----------------+
|  c.name  | c.livesleft |     p.city     |
+----------+-------------+----------------+
| Buster   |           6 | Ivanhoe        |
| Tiger    |           4 | Brambleton     |
| Lucy     |           1 | Dunlo          |
| Pepper   |           3 | Alleghenyville |
| Tiger    |           2 | Katonah        |
| Molly    |           6 | Babb           |
| Precious |           8 | Holcombe       |
+----------+-------------+----------------+
```
You can choose between table, tabbed, json and csv output formats.

## Installation
Either download the binary for your operating system (Linux, OS X and Windows are supported) from the [Releases page](https://github.com/cube2222/octosql/releases), or install using the go command line tool:
```bash
go get github.com/cube2222/octosql/cmd/octosql
```
## Configuration
The configuration file has the form
```yaml
dataSources:
  - name: <table_name_in_octosql>
    type: <datasource_type>
    config:
      <datasource_specific_key>: <datasource_specific_value>
      <datasource_specific_key>: <datasource_specific_value>
      ...
  - name: <table_name_in_octosql>
    type: <datasource_type>
    config:
      <datasource_specific_key>: <datasource_specific_value>
      <datasource_specific_key>: <datasource_specific_value>
      ...
    ...
```
### Supported Datasources
#### JSON
JSON file in one of the following forms:
- one record per line, no commas
- JSON list of records
##### options:
- path - path to file containing the data, required
- arrayFormat - if the JSON list of records format should be used, defaults to false

---
#### CSV
CSV file seperated using commas.
##### options:
- path - path to file containing the data, required

---
#### PostgreSQL
Single PostgreSQL database table.
##### options:
- address - address including port number, defaults to localhost:5432
- user - required
- password - required
- databaseName - required
- tableName - required

---
#### MySQL
Single MySQL database table.
##### options:
- address - address including port number, defaults to localhost:3306
- user - required
- password - required
- databaseName - required
- tableName - required

---
#### Redis
Redis database with the given index. Currently only hashes are supported.
##### options:
- address - address including port number, defaults to localhost:6379
- password - defaults to ""
- databaseIndex - index number of Redis database, defaults to 0
- databaseKeyName - column name of Redis key in OctoSQL records, defaults to "key"

## Documentation
The SQL dialect documentation:

Function documentation:

## Roadmap
- Add arithmetic operators.
- Write custom sql parser, so we can use sane function names.
- Push down functions to supported databases.
- Implement an in-memory index to save values of subqueries, so as not to recalculate them each time.
- MapReduce style distributed execution mode.
- Function Tables (RANGE(1, 10) for example)
- Better nested JSON support.
- HAVING clause.

## Architecture
### Project Structure