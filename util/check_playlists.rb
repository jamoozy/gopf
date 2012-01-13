#!/usr/bin/ruby -w

# Checks the playlists in "data/" to make sure they point to actual data.
# Prints warnings if "broken links" are found.

require 'ftools'

ROOT_DIR = 'data'
verbose = ARGV.size > 0 and ARGV[0] == '-v'

Dir["#{ROOT_DIR}/playlists/*"].each do |playlist|
  errors = 0
  lines = 0
  IO.foreach(playlist) do |fname|
    lines += 1
    fname = fname.sub('..', ROOT_DIR).strip
    unless File.exist?(fname)
      errors += 1
      puts "Error for \"#{fname}\"" if verbose
    end
  end
  puts "\"#{playlist}\" has #{errors}/#{lines} errors" unless errors == 0
end
