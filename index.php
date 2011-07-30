<!DOCTYPE html>
<?php
include("list.php");
include("mysql.php");

handle($_SERVER["REMOTE_ADDR"]);
?>

<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <title id="page-title">Zelda_64_Pachelbels_Ganon_OC_ReMix</title>
    <script src="audio.js" type="text/javascript"></script>
    <script src="playlist.js" type="text/javascript"></script>
    <script src="input.js" type="text/javascript"></script>
    <link rel="stylesheet" type="text/css" href="style.css">
    <link id="prefetch" rel="prefetch" href="">
  </head>
  <body>

    <div id="media-container" class="media-container">
      <div id="title-header">
        <h1 class="header">Now Playing:</h1>
        <div class="title">Zelda_64_Pachelbels_Ganon_OC_ReMix</div>
      </div>
      <div id="notification"></div>
      <audio id="player" src="data/Zelda_64_Pachelbels_Ganon_OC_ReMix.mp3"
             seek="true" controls> Hey, man, get an HTML5-compatible browser, okay?
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

      <div id="song-container" class="song-container">
        <ul id="songs" class="songs">
          <li class="dummy">(nothing loaded)</li>
        </ul>
        <h1 id="song-header" class="header">Songs</h1>
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
