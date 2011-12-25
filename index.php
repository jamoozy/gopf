<?php
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

include("mysql.php");
include("list.php");

if (!($ip_error = ip_is_ok($_SERVER["REMOTE_ADDR"]))) {
?>

<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <title id="page-title">Zelda_64_Pachelbels_Ganon_OC_ReMix</title>
    <script src="loc.js" type="text/javascript"></script>
    <script src="media.js" type="text/javascript"></script>
    <script src="playlist.js" type="text/javascript"></script>
    <script src="input.js" type="text/javascript"></script>
    <link rel="stylesheet" type="text/css" href="style.css">
    <link id="prefetch" rel="prefetch" href="">
  </head>

  <body>
    <div id="player-container" class="player-container">
      <div id="title-header">
        <h1 class="header">Now Playing:</h1>
        <div class="title"></div>
      </div>
      <div id="notification"></div>
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
        <ul id="playlists" class="playlists"><?=generate_playlists();?></ul>
        <h1 id="playlist-header" class="header">Playlists</h1>
      </div>

      <div id="media-container" class="media-container">
        <ul id="media" class="media">
          <li class="dummy">(nothing loaded)</li>
        </ul>
        <h1 id="media-header" class="header">Songs</h1>
      </div>
    </nav>

    <footer id="footer">
      <div class="name">
        <a href="http://code.google.com/p/gopf" target="_blank">
          GOPF: The GNU Online Player Framework
        </a>
      </div>
      <div class="name">
        Written by <author>Andrew "jamoozy" Correa</author>
      </div>
    </footer>
  </body>
</html>

<?php
} else {

// The following is all done inline (as opposed to having separate CSS
// and JS files) to decrease complexity.
?>
<html>
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
</html>
<?php
  die(1);
}
?>
