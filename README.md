##Tin

```bash
$ go test -v ./... -race
=== RUN   TestCache
--- PASS: TestCache (0.00s)
=== RUN   TestGetOrderedKeySlice
--- PASS: TestGetOrderedKeySlice (0.00s)
=== RUN   TestGetOrderedValuesSlice
--- PASS: TestGetOrderedValuesSlice (0.00s)
=== RUN   TestCacheConcurrent
--- PASS: TestCacheConcurrent (0.01s)
PASS
ok      github.com/serboox/tests        1.024s
```