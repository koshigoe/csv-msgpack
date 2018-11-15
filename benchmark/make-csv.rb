require 'csv'
require 'securerandom'

header = Array.new(99) { |i| i }

L = 10_000

File.open('data.csv', 'wb') do |f|
  f.puts(CSV.generate_line((['id'] + header)))

  L.times do |i|
    r = [format('%020d', i)]
    r += Array.new(99) { SecureRandom.hex(10).chars.join("\n") }
    f.puts(CSV.generate_line(r))
  end
end
