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

$playlist_dir = "data/playlists/";

# Generates playlist from playlist files.  Playlist files are simple text
# files with each line containing the relative path to a song.
#
# playlist: String name of playlist to start selected.
def generate_playlists(playlist)
  fnames = Dir[File.join(File.dirname(__FILE__), $playlist_dir, '*')].reject{ |f| f[0] == '.' or File.executable?(f) or f[-1] == "~"}.sort

  rtn = '';
  if playlist
    fnames.each do |fname|
      if playlist == fname
        rtn << "<li class=\"unselected\">#{fname}</li>"
      else
        rtn << "<li class=\"unselected selected\" id=\"selected\">#{fname}</li>"
      end
    end
  else
    fnames.each {|fname| rtn << "<li class=\"unselected\">#{fname}</li>"}
  end
  return $rtn;
end

def generate_media(playlist, media=false)
  paths = File.readlines(File.join($playlist_dir, playlist)).map{|f| f.strip}
  html = '';
  paths.reject{|p| p.size == 0}.each do |path|
    name = path[path.rindex('/')...-4]
    html << "<li class=\"media#{name == media ? ' playing' : ''}\" path=\"#{path.sub('..', 'data')}\" onclick=\"media.onclick(this)\">#{name}</li>\n";
  end

  return html;
end

if __FILE__ == $0
  $cgi = CGI.new

  if $cgi.params.has_key?('op')
    if $cgi.params['op'] == 'ls'
      if $cgi.params.has_key?('dir')
        #ls($_GET['dir']);
        $cgi.out("text/json"){'{error:"ls not implemented."}'}
      else
        $cgi.out("text/json"){'{error:"Expected \"dir\" key."}'}
      end
    end
  end

  # Contents of return page are just text.  Playlist must be in utf-8.
  if $cgi.params.has_key? 'playlist'
    $cgi.out("text/plain") do
      File.readlines(File.join($playlist_dir, $cgi.params['playlist'])).join
    end
  end
end
