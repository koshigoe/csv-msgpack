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
$ ruby -v benchmark-parse.rb data.csv
ruby 2.5.3p105 (2018-10-18 revision 65156) [x86_64-linux]
                   user     system      total        real
csv-msgpack    0.133354   0.015763   0.326335 (  0.177472)
CSV.foreach   26.207001   0.028033  26.235034 ( 26.236198)
```
