# Ivory

Ivory is a tool that can do port scanning between ip blocks dynamically.
You can view the results as csv or send the results to your notification address.

#### Build
> go build -o scanner .

#### Parameters
parameter | required | type | default | description
--- | --- | --- | --- |---
ip | true | string | 192.168.* . * | ip address third block [0-255]
port | true | integer | 9200 | scan port
label | true | string | ElasticSearch | label
storage | false | string | csv | storage type (notification url or csv)
concurrent-count | false | integer | 50 | concurrent count

#### Example

> ./scanner --ip=192.168.1 .* --port=9200 --label=ElasticSearch --storage=http://notification.address

> ./scanner --ip=192.168.1 .* --port=9200 --label=ElasticSearch --storage=csv

 - The octet to look at in this example is 192.168.1.[0-255].
 - The fewer octets defined, the longer the scanning process.
 - If storage is empty writes to result.csv file, but if you give a url, it sends the request to the corresponding url with the get method.The query string contains ip, port and label parameters.
