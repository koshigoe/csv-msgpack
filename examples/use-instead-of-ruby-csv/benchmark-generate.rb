require 'benchmark'
require 'csv'
require 'msgpack'
require 'open3'

row = Array.new(100) do |i|
  chr = ('A'.ord + (i % 26)).chr
  Array.new(100, chr).join("\n")
end

L = 10_000

Benchmark.bm(20) do |x|
  x.report('CSV.generate_line') do
    L.times { CSV.generate_line(row) }
  end

  x.report('csv-msgpack decode') do
    Open3.popen3('../../csv-msgpack decode -o benchmark-generate.csv') do |stdin, stdout, stderr, wait_thr|
      L.times do
        stdin.write(MessagePack.pack(row))
      end
      stdin.close
    end
  end
end

__END__
$ ruby -v benchmark-generate.rb
ruby 2.5.3p105 (2018-10-18 revision 65156) [x86_64-darwin18]
                           user     system      total        real
CSV.generate_line      3.100018   0.011189   3.111207 (  3.120140)
csv-msgpack decode     0.253302   0.149663   4.513202 (  4.269278)
