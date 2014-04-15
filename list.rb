#!/usr/bin/ruby -w

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


require 'cgi'


if __FILE__ == $0
  $playlist_dir = "data/playlists/";
  $cgi = CGI.new

  # Contents of return page are just text.  Playlist must be in utf-8.
  if $cgi.params.has_key?('playlist')
    playlist = File.join($playlist_dir, $cgi.params['playlist'])
    if File.exists?("#{playlist}.json")
      $cgi.out("text/json"){File.readlines("#{playlist}.json").join}
    elsif File.exists?(playlist)
      $cgi.out("text/plain"){File.readlines(playlist).join}
    else
      $cgi.out("text/json"){"{error:'No such file: \"#{playlist}\"'"}
    end
  end
end
