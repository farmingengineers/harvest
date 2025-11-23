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

class Stats
  def initialize(data)
    @data = data
  end

  attr_reader :data

  def stddev
    Math.sqrt(variance)
  end

  def variance
    @data.sum { |x| (x - mean) * (x - mean) }
  end

  def mean
    @mean ||= @data.sum.to_f / @data.size
  end

  def max
    @max ||= @data.max
  end

  def min
    @min ||= @data.min
  end
end

def covariance(sx, sy)
  if sx.data.size != sy.data.size
    raise "Number of data points for crop 1 (#{sx.data.size}) must be the same as the number for crop 2 (#{sy.data.size})."
  end
  sx.data.size.times.sum { |i| (sx.data[i] - sx.mean) * (sy.data[i] - sy.mean) }
end

sx = Stats.new(x)
sy = Stats.new(y)

printf "%-20s  range = %6.2f-%6.2f  mean = %6.2f  stddev = %6.2f\n",
  crop1, sx.min, sx.max, sx.mean, sx.stddev
printf "%-20s  range = %6.2f-%6.2f  mean = %6.2f  stddev = %6.2f\n",
  crop2, sy.min, sy.max, sy.mean, sy.stddev

correlation = covariance(sx, sy) / (sx.stddev * sy.stddev)
printf "correlation = %5.3f\n",
  correlation
