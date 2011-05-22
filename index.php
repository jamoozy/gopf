<!DOCTYPE html>
<?php include("list.php"); ?>
<?php include("mysql.php");
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
  <body onkeydown="input.onkey(event)" onload="input.init()">

    <center class="media-container">
      <div id="title-header">
        <h1 class="header">Now Playing:</h1>
        <div class="title">Zelda_64_Pachelbels_Ganon_OC_ReMix</div>
      </div>
      <div id="notification"></div>
      <audio id="player" src="data/Zelda_64_Pachelbels_Ganon_OC_ReMix.mp3"
             onplay="audio.onplay(this)"
             onended="audio.onended(this)"
             ontimeupdate="audio.onprogress(event)"
             seek="true" controls> Hey, man, get an HTML5-compatible browser, okay?
      </audio>
      <div id="controls" class="controls-container">
        <span id="prev" class="controls" onclick="audio.prev()">&lt;&lt;</span>
        <input id="loop" type="checkbox">
          <span class="controls" onclick="document.getElementById('loop').click()"><span class="mnemonic">l</span>oop</span>
        </input>
        <input id="shuf" type="checkbox">
          <span class="controls" onclick="document.getElementById('shuf').click()"><span class="mnemonic">s</span>huffle</span>
        </input>
        <span id="next" class="controls" onclick="audio.next()">&gt;&gt;</span>
      </div>
    </center>

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

  </body>
</html>
