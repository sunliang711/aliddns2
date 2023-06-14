## usage
### with all params
curl -X POST -H "Content-Type: application/json" -d '{"new_ip":"1.2.3.4","rr":"test1","domain":"gitez.cc","type":"A","ttl":"600","access_key":"","access_secret":""}' http://localhost:8080/updateRecord


### with required params
curl -X POST -H "Content-Type: application/json" -d '{"rr":"test1","domain":"gitez.cc"}' http://localhost:8080/updateRecord