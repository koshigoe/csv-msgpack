csv-msgpack
====


Example
----

```ruby
require 'csv'

header = Array.new(99) { |i| i }
cols = header.map do |i|
  chr = ('A'.ord + (i % 26)).chr
  Array.new(100, chr).join("\n")
end

L = 10_000

File.open("100chr-100col-#{L}row.csv", 'wb') do |f|
  f.puts(CSV.generate_line((['id'] + header)))

  L.times do |i|
    r = [format('%0100d', i)] + cols
    f.puts(CSV.generate_line(r))
  end
end
```

```ruby
# benchmark-parse.rb
require 'benchmark'
require 'csv'
require 'msgpack'
require 'open3'

path = ARGV[0]

Benchmark.bm(20) do |x|
  x.report('CSV.foreach') do
    CSV.foreach(path) do |_|
    end
  end

  x.report('csv-msgpack encode') do
    Open3.popen3("./csv-msgpack encode -i #{path}") do |stdin, stdout, stderr, wait_thr|
      stdin.close

      u = MessagePack::Unpacker.new(stdout)
      u.each do |_|
      end
    end
  end
end

__END__

$ ruby benchmark-parse.rb 100chr-100col-10000row.csv
                           user     system      total        real
CSV.foreach          217.847741   1.345944 219.193685 (226.247598)
csv-msgpack encode     0.579895   0.099204   4.679996 (  3.981893)
```
