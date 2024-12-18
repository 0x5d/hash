```sh
ab -c 53 -t 180 -n 2000000 -m GET k8s-backend-hash-e9d8d347e8-226d5937ed46c589.elb.us-east-1.amazonaws.com/miss
This is ApacheBench, Version 2.3 <$Revision: 1903618 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking k8s-backend-hash-e9d8d347e8-226d5937ed46c589.elb.us-east-1.amazonaws.com (be patient)
Completed 200000 requests
Completed 400000 requests
Completed 600000 requests
Finished 705983 requests


Server Software:        
Server Hostname:        k8s-backend-hash-e9d8d347e8-226d5937ed46c589.elb.us-east-1.amazonaws.com
Server Port:            80

Document Path:          /miss
Document Length:        38 bytes

Concurrency Level:      53
Time taken for tests:   180.000 seconds
Complete requests:      705983
Failed requests:        0
Non-2xx responses:      705983
Total transferred:      146138481 bytes
HTML transferred:       26827354 bytes
Requests per second:    3922.12 [#/sec] (mean)
Time per request:       13.513 [ms] (mean)
Time per request:       0.255 [ms] (mean, across all concurrent requests)
Transfer rate:          792.85 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        2    2   5.2      2    1060
Processing:     3   11   5.1     10     255
Waiting:        3   11   5.1     10     254
Total:          5   13   7.3     13    1071

Percentage of the requests served within a certain time (ms)
  50%     13
  66%     13
  75%     14
  80%     14
  90%     16
  95%     19
  98%     23
  99%     27
 100%   1071 (longest request)
 ```