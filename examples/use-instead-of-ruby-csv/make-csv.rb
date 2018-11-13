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
