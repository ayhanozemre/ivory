# Ivory

**Ivory** is a port scanner for ip blocks.

**How to build**
> go build -o scanner .

## Parameters

parameter | required | type | default | description
--- | --- | --- | --- |---
port | true | integer | 9200 | scan port
label | true | string | ElasticSearch | label
storage | false | string | csv | storage type (notification url or csv)
concurrent-count | false | integer | 50 | concurrent count
first-block | false | integer | 0 | ip address first block [0-255]
second-block | false | integer | 0 | ip address second block [0-255]
third-block | false | integer | 0 | ip address third block [0-255]

**Example**
port scan for ElasticSearch
> ./scanner  --port=9200 --label=elasticSearch --storage=http://localhost:8080/portscanner/create --first-block=192 --second-block=168 --third-block=1

> ./scanner  --port=9200 --label=ElasticSearch --first-block=192 --second-block=168  --third-block=1 --storage=csv

 - The block to look at in this example is 192.168.1.[0-255].
 - The fewer blocks defined, the longer the scanning process.
 - If storage is empty writes to result.csv file, but if you give a url, it sends the request to the corresponding url with the get method.The query string contains ip, port and label parameters.
