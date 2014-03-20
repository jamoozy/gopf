<?php
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

include("mysql.php");
include("list.php");

if (!($ip_error = ip_is_ok($_SERVER["REMOTE_ADDR"]))) {
  $playlist = false;
  $media = false;
  if ($_GET) {
    if (array_key_exists('p', $_GET)) {
      $playlist = urldecode($_GET['p']);
      if (array_key_exists('m', $_GET)) {
        $media = urldecode($_GET['m']);

        // TODO
        //if (array_key_exists('t', $_GET)) {
        //  $time = urldecode($_GET['t']);
        //}
      }
    }
  }
}
?>

<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8"/>
    <title id="page-title">Zelda_64_Pachelbels_Ganon_OC_ReMix</title>
    <script src="jquery-1.11.0.min.js" type="text/javascript"></script>
    <!--script src="loc.js" type="text/javascript"></script-->
    <script src="media.js" type="text/javascript"></script>
    <script src="playlist.js" type="text/javascript"></script>
    <script src="input.js" type="text/javascript"></script>
    <link rel="stylesheet" type="text/css" href="http://fonts.googleapis.com/css?family=Poiret+One|Tinos|Headland+One">
    <link rel="stylesheet" type="text/css" href="style.css">
    <link id="prefetch" rel="prefetch" href="">
  </head>

  <body>
    <div id="player-container" class="player-container">
      <div id="title-header">
        <h1 class="header">now playing:</h1>
        <div class="title" style="font-family:'Poiret One'">(nothing loaded)</div>
      </div>
      <div id="notification"> </div>
      <audio id="player" src="" seek="true" controls>
        Hey, man, get an HTML5-compatible browser, okay?
      </audio>
      <div id="controls" class="controls-container">
        <span id="prev" class="controls">&lt;&lt;</span>
        <input id="loop" type="checkbox">
          <span id="loop_label" class="controls"><span class="mnemonic">l</span>oop</span>
        </input>
        <input id="shuf" type="checkbox">
          <span id="shuf_label" class="controls"><span class="mnemonic">s</span>huffle</span>
        </input>
        <span id="next" class="controls">&gt;&gt;</span>
      </div>
    </div>

    <nav id="navigator">
      <div id="playlist-container" class="playlist-container">
        <ul id="playlists" class="playlists"><?=generate_playlists($playlist)?></ul>
        <h1 id="playlist-header" class="header">Playlists</h1>
      </div>

      <div id="media-container" class="media-container">
        <ul id="media" class="media">
        <?if ($media) {?>
          <?=generate_media($playlist, $media)?>
        <?}else{?>
          <?=generate_media($playlist)?>
        <?}?>
        </ul>
        <h1 id="media-header" class="header">Songs</h1>
      </div>
    </nav>

    <div id="help-dialog" class="help-dialog"></div>

    <footer id="footer">
      <div class="name">
        <a href="http://github.com/jamoozy/gopf" target="_blank">
          GOPF: The GNU Online Player Framework
        </a>
      </div>
      <div class="name">
        Written by <author>Andrew "Jamoozy" Correa</author>
      </div>
    </footer>
  </body>
</html>

<!--? if ($playlist && $media) { ?>
<script type="text/javascript">
// Handle Clicked media
var clickPassedMedia = function(req) {
  console.log('clickPassedMedia');
  var media = document.getElementsByClassName("media");
  var not = document.getElementById("notification");
  for (var i = 0; i < media.length; i++) {
    if (media[i].innerHTML == "<?=$media?>") {
      media[i].onclick(media[0]);
      return;
    }
  }
  // Not found!  Log/print error?
}

// Loads the playlist with class="selected".
function loadSelectedPlaylist(e) {
  console.log('loadSelectedPlaylist');
  var selected = document.getElementsByClassName("selected");
  if (selected.length > 0) {
    elem = selected[0];
    elem.setAttribute("class", "unselected");
    playlist.onclick(elem, false, clickPassedMedia);
  } else {
    // Report error?
  }
}

// Listen for the page to load, so that this can load the playlist.
window.addEventListener("load", loadSelectedPlaylist, true);
</script>
<!--? } ?>

<!--?php
//} else {

// The following is all done inline (as opposed to having separate CSS
// and JS files) to decrease complexity.
?--!>
<!--html>
  <head>
    <style type="text/css">
      body {
        text-align: center;
        font-family: sans-serif;
      }
      div {
        border: solid 1px black;
        margin: 40px 20px 5px 20px;
        padding: 20px;
        color: white;
        background-color: red;

        -webkit-border-radius: 1em;
        -moz-border-radius: 1em;
        border-radius: 1em;
      }
      h1 {
        font-size: 20pt;
        margin: 0;
        margin-bottom: 5px;
        padding: 0;
      }
      p {
        font-size: 12pt;
        margin: 0px;
        padding: 0px;
      }
    </style>
  </head>
  <body>
    <div>
      <h1>NO ME GUSTA!</h1>
      <p> You're not in the registered users! </p>
      <p> Contact someone to fix this. </p>
      <hr>
      <p> This website has logged:</p><?=$ip_error?>
    </div>
  </body>
</html-->
<?php
//  die(1);
//}
?>
