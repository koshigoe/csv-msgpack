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
    Open3.popen3("../../csv-msgpack encode -i #{path}") do |stdin, stdout, stderr, wait_thr|
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
