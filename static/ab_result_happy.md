```sh
ab -c 53 -t 180 -n 2000000 -m GET k8s-backend-hash-e9d8d347e8-226d5937ed46c589.elb.us-east-1.amazonaws.com/b
This is ApacheBench, Version 2.3 <$Revision: 1903618 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking k8s-backend-hash-e9d8d347e8-226d5937ed46c589.elb.us-east-1.amazonaws.com (be patient)
Completed 200000 requests
Completed 400000 requests
Completed 600000 requests
Completed 800000 requests
Completed 1000000 requests
Completed 1200000 requests
Completed 1400000 requests
Completed 1600000 requests
Finished 1687588 requests


Server Software:        
Server Hostname:        k8s-backend-hash-e9d8d347e8-226d5937ed46c589.elb.us-east-1.amazonaws.com
Server Port:            80

Document Path:          /b
Document Length:        58 bytes

Concurrency Level:      53
Time taken for tests:   180.000 seconds
Complete requests:      1687588
Failed requests:        0
Non-2xx responses:      1687589
Total transferred:      378019936 bytes
HTML transferred:       97880162 bytes
Requests per second:    9375.49 [#/sec] (mean)
Time per request:       5.653 [ms] (mean)
Time per request:       0.107 [ms] (mean, across all concurrent requests)
Transfer rate:          2050.89 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        1    2   3.3      2    1064
Processing:     2    4   1.3      3     245
Waiting:        2    3   1.2      3     245
Total:          3    6   3.5      5    1068

Percentage of the requests served within a certain time (ms)
  50%      5
  66%      6
  75%      6
  80%      6
  90%      7
  95%      8
  98%      9
  99%      9
 100%   1068 (longest request)
 ```