require 'benchmark'
require 'csv'
require 'msgpack'
require 'open3'

path = ARGV[0]

Benchmark.bm(12) do |x|
  x.report('csv-msgpack') do
    Open3.popen3("../target/release/csv-msgpack < #{path}") do |stdin, stdout, stderr, wait_thr|
      stdin.close

      u = MessagePack::Unpacker.new(stdout)
      u.each do |_|
      end
    end
  end

  x.report('CSV.foreach') do
    CSV.foreach(path) do |_|
    end
  end
end
