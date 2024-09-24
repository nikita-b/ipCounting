# Counting IP

A small tool to read unique IPs from a text file.
We have three options: regular bitmap vs Roaring bitmap using a library vs a concurrent solution based on Roaring bitmap.
The concurrent solution isn't good for this task because it's an I/O-bound task, and we need to minimize memory usage.
But just for fun, I've implemented it.

## Run counting
```bash
ARGS="-filename <path_to_file>" make run
```
Or a more complicated example:
```bash
ARGS="-profile -concurrency <amount of workers> --filename <path_to_file>  -algo <algoritm>" make run

-profile : - enable profiling
-concurrency : - amount of workers for concurrent solution. Makes sense only for concurrent solution
-algo : 0 - bitmap 1 - roaring bitmap 2 - concurrent

````

## Run tests
```bash
GENERATED_IPS=500000 make test
```

## Memory Consumption
(Apple M1) 100 GB File

### Regular bitmap:
```
Showing nodes accounting for 512MB, 100% of 512MB total
      flat  flat%   sum%        cum   cum%
     512MB   100%   100%      512MB   100%  main.NewBitmapIPCounter (inline)
         0     0%   100%      512MB   100%  main.main
         0     0%   100%      512MB   100%  runtime.main
```
Execution time is ~10m01s

### Roaring bitmap:
```
Showing nodes accounting for 509.45MB, 99.64% of 511.29MB total
Dropped 1 node (cum <= 2.56MB)
      flat  flat%   sum%        cum   cum%
  509.45MB 99.64% 99.64%   509.45MB 99.64%  github.com/RoaringBitmap/roaring/v2.newBitmapContainer (inline)
         0     0% 99.64%   511.29MB   100%  github.com/RoaringBitmap/roaring/v2.(*Bitmap).Add
         0     0% 99.64%   509.45MB 99.64%  github.com/RoaringBitmap/roaring/v2.(*arrayContainer).iaddReturnMinimized
         0     0% 99.64%   509.45MB 99.64%  github.com/RoaringBitmap/roaring/v2.(*arrayContainer).toBitmapContainer
         0     0% 99.64%   511.29MB   100%  github.com/RoaringBitmap/roaring/v2/roaring64.(*Bitmap).Add
         0     0% 99.64%   511.29MB   100%  main.(*IPCounterRoaring).Add
         0     0% 99.64%   511.29MB   100%  main.ProcessFile
         0     0% 99.64%   511.29MB   100%  main.main
         0     0% 99.64%   511.29MB   100%  runtime.main
```
Execution time is ~14m10s


### Concurrent solution:
```
(pprof) top
Showing nodes accounting for 2071.03MB, 99.69% of 2077.38MB total
Dropped 8 nodes (cum <= 10.39MB)
      flat  flat%   sum%        cum   cum%
 2071.03MB 99.69% 99.69%  2071.03MB 99.69%  github.com/RoaringBitmap/roaring/v2.newBitmapContainer (inline)
         0     0% 99.69%  2074.84MB 99.88%  github.com/RoaringBitmap/roaring/v2.(*Bitmap).AddMany
         0     0% 99.69%   888.15MB 42.75%  github.com/RoaringBitmap/roaring/v2.(*Bitmap).addwithptr
         0     0% 99.69%  2071.03MB 99.69%  github.com/RoaringBitmap/roaring/v2.(*arrayContainer).iaddReturnMinimized
         0     0% 99.69%  2071.03MB 99.69%  github.com/RoaringBitmap/roaring/v2.(*arrayContainer).toBitmapContainer
         0     0% 99.69%  2074.84MB 99.88%  github.com/RoaringBitmap/roaring/v2/roaring64.(*Bitmap).AddMany
         0     0% 99.69%  2074.84MB 99.88%  main.(*IPCounterConcurrent).AddConcurrent
         0     0% 99.69%  2074.84MB 99.88%  main.ProcessFileConcurrency.func2
```
Execution time is ~21m10s

Solution is slower but at least uses more memory! :(((

It's much better with a smaller file. Need to investigate more.