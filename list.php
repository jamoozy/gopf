<?
// Copyright 2012 Andrew "Jamoozy" Correa
//
// This file is part of GOPF.
//
// GOPF is free software: you can redistribute it and/or modify it under
// the terms of the GNU General Public as published by the Free Software
// Foundation, either version 3 of the License, or (at your option) any
// later version.
//
// Foobar is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License
// for more details.
//
// You should have received a copy of the GNU General Public License
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
function generate_playlists() {
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
    foreach ($fnames as $fname) {
      $rtn.="<li class=\"unselected\" onclick=\"playlist.onclick(this)\">$fname</li>";
    }
  } else {
    $rtn.="<div class=\"error\">An internal error occurred, and your request could not be completed.</div>";
  }
  return $rtn;
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
