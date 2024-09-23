# Counting IP

Small tool to read uniq IP from text file.
We have three options: regular bitmap vs roaring bitmap using library vs .


## Run couning
```bash
ARGS="-filename <path_to_file>" make run
```

## Run tests
```bash
GENERATED_IPS=500000 make test
```

## Memory consuming
(Apple M1)
Regular bitmap:
```
Showing nodes accounting for 512MB, 100% of 512MB total
      flat  flat%   sum%        cum   cum%
     512MB   100%   100%      512MB   100%  main.NewBitmapIPCounter (inline)
         0     0%   100%      512MB   100%  main.main
         0     0%   100%      512MB   100%  runtime.main
```
Execution time is ~10m19s

Roaring bitmap:
```
Showing nodes accounting for 17024.55kB, 100% of 17024.55kB total
Showing top 10 nodes out of 17
      flat  flat%   sum%        cum   cum%
15480.31kB 90.93% 90.93% 15480.31kB 90.93%  github.com/RoaringBitmap/roaring/v2.newBitmapContainer (inline)
 1032.02kB  6.06% 96.99% 16512.33kB 96.99%  github.com/RoaringBitmap/roaring/v2.(*arrayContainer).iaddReturnMinimized
  512.22kB  3.01%   100%   512.22kB  3.01%  runtime.malg
         0     0%   100% 16512.33kB 96.99%  github.com/RoaringBitmap/roaring/v2.(*Bitmap).Add
         0     0%   100% 15480.31kB 90.93%  github.com/RoaringBitmap/roaring/v2.(*arrayContainer).toBitmapContainer
         0     0%   100% 16512.33kB 96.99%  main.(*IPCounterRoaring).Add
         0     0%   100% 16512.33kB 96.99%  main.ProcessFile
         0     0%   100% 16512.33kB 96.99%  main.main
         0     0%   100%   512.22kB  3.01%  runtime.allocm
         0     0%   100% 16512.33kB 96.99%  runtime.main
```
Execution time is ~14m10s


