csv-msgpack
====


Description
----

NOTE: This is an experimental.

The `csv-msgpack` parse CSV and serialize each record as MessagePack.


Example
----

```
$ cargo build --release
$ echo a,b,c \
  | ./target/release/csv-msgpack  \
  | ruby -r msgpack -e 'MessagePack::Unpacker.new(STDIN).each { |x| p x }'
["a", "b", "c"]
```


Benchmark
----

### Compare with the Ruby CSV

```
$ cd benchmark
$ gem install msgpack
$ ruby make-csv.rb
$ ruby -v benchmark-parse.rb 100chr-100col-10000row.csv
ruby 2.5.3p105 (2018-10-18 revision 65156) [x86_64-linux]
                   user     system      total        real
csv-msgpack    0.186039   0.058284   0.797959 (  0.553839)
CSV.foreach  128.064341   0.081668 128.146009 (128.161941)
```
