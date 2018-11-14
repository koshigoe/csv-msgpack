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

  x.report('csv-msgpack - Golang') do
    Open3.popen3("../../csv-msgpack encode -i #{path}") do |stdin, stdout, stderr, wait_thr|
      stdin.close

      u = MessagePack::Unpacker.new(stdout)
      u.each do |_|
      end
    end
  end

  x.report('csv-msgpack - Rust') do
    Open3.popen3("../../target/release/csv-msgpack < #{path}") do |stdin, stdout, stderr, wait_thr|
      stdin.close

      u = MessagePack::Unpacker.new(stdout)
      u.each do |_|
      end
    end
  end
end

__END__
$ ruby -v benchmark-parse.rb 100chr-100col-10000row.csv
ruby 2.5.3p105 (2018-10-18 revision 65156) [x86_64-darwin18]
                           user     system      total        real
CSV.foreach          251.196207   1.892531 253.088738 (273.155054)
csv-msgpack encode     0.761726   0.084683   5.311952 (  4.500843)

                           user     system      total        real
CSV.foreach          225.262108   1.418877 226.680985 (238.651944)
csv-msgpack - Golang   0.565393   0.093311   4.628336 (  3.944491)
csv-msgpack - Rust     0.544291   0.075604   1.544383 (  0.965783)
