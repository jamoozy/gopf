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


$data_dir = "data"
$playlist_dir = "#$data_dir/playlists/";


# Generates playlist from playlist files.  Playlist files are simple text
# files with each line containing the relative path to a song.
#
#   playlist: String name of playlist to start selected.  Default is none
#       selected (false).
#
# Returns the HTML for the playlist list.
def generate_playlists(playlist=false)
  fnames = Dir[$playlist_dir].map do |fname|
    raise "#{fname} DNE?" unless File.exists?(fname) # sanity check
    next if fname.start_with?('.') or File.executable?(fname) or fname[-1..-1] == "~"
  end.reject {|x| x == nil }.sort

  if playlist
    fnames.inject("") do |prev,fname|
      if playlist == fname
        prev + "<li class=\"unselected\">#{fname}</li>"
      else
        prev + "<li class=\"unselected selected\" id=\"selected\">#{fname}</li>"
      end
    end
  else
    fnames.inject('') do |prev,fname|
      prev + "<li class=\"unselected\">#{fname}</li>";
    end
  end
end


# Generates the list of media in the media list.  This list could contain
# songs or videos.
#
#   playlist: The name of the playlist.
#   media: The name of the selected media (if any).  Default is none selected
#       (false).
#
# Returns the HTML for the contents of the playlist with the passed name.
def generate_media(playlist, media=false)
  paths = split("\n", file_get_contents("#$playlist_dir/#{playlist}"));

  paths.inject('') do |prev, path|
    if path.strip.len <= 0
      prev
    else
      name = path[path.rindex('/')+1..-1];
      prev + "<li class=\"media#{name == media ? ' playing' : ''}\" path=\"#{path.sub(/\.\./, $data_dir)}\" onclick=\"media.onclick(this)\">#{name}</li>\n";
    end
  end
end


# Entry point for when this is called directly from a JS file on someone's
# client.  Additionally (thanks to the way the CGI module works), this is also
# the entry point for CL tests.
if __FILE__ == $0
  cgi = CGI.new

  # Contents of return page are just text.  Playlist must be in utf-8.
  if cgi.params.has_key?('op')
    case cgi.params['op']
    when  'ls'
      if cgi.params.has_key?('dir')
        cgi.out('text/plain') { Dir[cgi.params['dir']] }
      else
        cgi.out("text/json") { '{error:"Expected \"dir\" key."}' }
      end
    when 'playlist'
      generate_playlists(cgi.params.has_key?('playlist') ? cgi.params['playlist'] : false)
    end
  elsif cgi.params.has_key?('playlist')
    playlist = File.join($playlist_dir, cgi.params['playlist'])
    if File.exists?("#{playlist}.json")
      cgi.out("text/json") { File.readlines("#{playlist}.json").join }
    elsif File.exists?(playlist)
      cgi.out("text/plain") { File.readlines(playlist).join }
    else
      cgi.out("text/json") { "{error:'No such file: \"#{playlist}\"'" }
    end
  end
end
