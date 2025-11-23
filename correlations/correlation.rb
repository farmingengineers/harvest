require "csv"

if ARGV.size != 2
  puts "Usage: ruby correlation.rb CROP1 CROP2"
  puts "See the available crops with 'ruby list-crops.rb'."
  exit 1
end

crop1, crop2 = ARGV
x = []
y = []
Dir["harvest/*.csv"].sort.each do |f|
  c1 = nil
  c2 = nil
  CSV.new(File.open(f)).each do |r|
    c1 = r if r.first == crop1
    c2 = r if r.first == crop2
  end
  if c1 && c2 && c1.size == c2.size
    x += c1.map(&:to_f)
    y += c2.map(&:to_f)
  end
end

idx = (0..x.size-1).to_a
skips = idx.select { |i| x[i] == 0.0 && y[i] == 0.0 }
keep = idx - skips
x = keep.map { |i| x[i] }
y = keep.map { |i| y[i] }

idx = (0..x.size-1).to_a
mean_x = x.sum.to_f / x.size
mean_y = y.sum.to_f / y.size

dx = x.map { |v| v - mean_x }
dy = y.map { |v| v - mean_y }
rxy = (dx.zip(dy).sum { |v1, v2| v1 * v2 }) /
  (
    Math.sqrt(dx.sum { |v| v * v }) *
    Math.sqrt(dy.sum { |v| v * v })
  )

printf "%-20s  range %6.2f .. %6.2f  mean = %6.2f\n",
  crop1, x.min, x.max, mean_x
printf "%-20s  range %6.2f .. %6.2f  mean = %6.2f\n",
  crop2, y.min, y.max, mean_y

printf "correlation = %5.3f\n", rxy
