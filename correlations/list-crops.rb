require "csv"
require "set"

crops = Set.new
Dir["harvest/*.csv"].each do |f|
  has_header = false
  CSV.new(File.open(f)).each do |r|
    label = r.first
    next if label.nil? || label == ""
    if label =~ /^Week/
      has_header = true
      next
    end
    unless has_header
      $stderr.puts "#{f}: no header found before this row:", "  #{r.to_csv}"
      break
    end
    crops << label
  end
end

puts crops.to_a.sort
