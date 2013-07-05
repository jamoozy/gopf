#!/usr/bin/ruby -w
#
# Copyright 2012-2013 Andrew "Jamoozy" Correa S.
#
# This file is part of GOPF.
#
# GOPF is free software: you can redistribute it and/or modify it under
# the terms of the GNU Affero General Public as published by the Free Software
# Foundation, either version 3 of the License, or (at your option) any
# later version.
#
# GOPF is distributed in the hope that it will be useful, but WITHOUT
# ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
# FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License
# for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with GOPF. If not, see http://www.gnu.org/licenses/.

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
