#!/usr/bin/env ruby
#
# Copyright (C) 2011-2015 Andrew "Jamoozy" Sabisch
#
# This file is part of GOPF.
#
# GOPF is free software: you can redistribute it and/or modify it under the
# terms of the GNU Affero General Public as published by the Free Software
# Foundation, either version 3 of the License, or (at your option) any later
# version.
#
# GOPF is distributed in the hope that it will be useful, but WITHOUT ANY
# WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
# A PARTICULAR PURPOSE. See the GNU Affero General Public License for more
# details.
#
# You should have received a copy of the GNU Affero General Public License
# along with GOPF. If not, see http://www.gnu.org/licenses/.

require 'cgi'
require 'erb'

$cgi = CGI.new

bind = binding

$cgi.out('text/html') do
  File.open('index.html.erb', "r") do |f|
    ERB.new(f.read).result(bind)
  end
end
