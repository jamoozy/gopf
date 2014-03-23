<?
// Copyright 2012-2013 Andrew "Jamoozy" Correa S.
//
// This file is part of GOPF.
//
// GOPF is free software: you can redistribute it and/or modify it under
// the terms of the GNU Affero General Public as published by the Free Software
// Foundation, either version 3 of the License, or (at your option) any
// later version.
//
// GOPF is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with GOPF. If not, see http://www.gnu.org/licenses/.

$playlist_dir = "data/playlists/";

// Gets the contents of a directory in JSON format.
function ls($dir) {
  $json = "{ls:[";
  $handle = opendir($dir);
  while (($entry = readdir($handle)) !== false) {
    $json .= "'$entry',";
  }
  $json .= "]}";
  return $json;
}

// Generates playlist from playlist files.  Playlist files are simple text
// files with each line containing the relative path to a song.
//
// $playlist: String name of playlist to start selected.
function generate_playlists($playlist) {
  global $playlist_dir;

  $dir = dirname(__FILE__)."/$playlist_dir";
  $iter = new DirectoryIterator($dir);

  $fnames = array();
  foreach ($iter as $fi) {
    $fname = $fi->getFilename();

    if ($fi->isDot() or $fi->isExecutable()) { continue; }
    if (strcmp(substr($fname, 0, 1), ".") == 0) { continue; }
    if (strcmp(substr($fname, strlen($fname) - 1), "~") == 0) { continue; }

    $fnames[] = $fname;
  }

  $rtn = '';
  if (sort($fnames)) {
    if ($playlist) {
      foreach ($fnames as $fname) {
        if (strcmp($playlist, $fname)) {
          $rtn.='<li class="unselected">'.$fname.'</li>';
        } else {
          $rtn.='<li class="unselected selected" id="selected">'.$fname.'</li>';
        }
      }
    } else {
      foreach ($fnames as $fname) {
        $rtn.="<li class=\"unselected\">$fname</li>";
      }
    }
  } else {
    $rtn.="<div class=\"error\">\n";
    $rtn.="  An internal error occurred, and your request could not be";
    $rtn.="  completed.\n";
    $rtn.="</div>";
  }
  return $rtn;
}

function generate_media($playlist, $media=false) {
  global $playlist_dir;

  $paths = split("\n", file_get_contents($playlist_dir.$playlist));
  $html = '';
  foreach ($paths as $path) {
    if (strlen(trim($path)) <= 0) {
      continue;
    }
    $name = substr($path, strrpos($path, '/') + 1, -4);
    $html .= '<li class="media'.(strcmp($name, $media) == 0 ? ' playing' : '').'" path="'.str_replace('..', 'data', $path)."\" onclick=\"media.onclick(this)\">$name</li>\n";
  }

  return $html;
}

if ($_GET) {
  if (array_key_exists('op', $_GET)) {
    if ($_GET['op'] == "ls") {
      echo ls($_GET['dir']);
    }
  }

  // Contents of return page are just text.  Playlist must be in utf-8.
  if (array_key_exists('playlist', $_GET)) {
    echo file_get_contents($playlist_dir.$_GET['playlist']);
  }
}
?>
