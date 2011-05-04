var audio = (function() {
  // Adds a string to the notification thing.
  function notify(str) {
    document.getElementById("notification").innerHTML = str;
  }

  function inspect(obj) {
    var out = '';
    for (var p in obj) {
      out += p + ': ' + obj[p] + '\n';
    }
    notify(out);
  }

  function get_ith_song(i) {
    var songs = document.getElementById("songs").childNodes;
    if (songs.length > i) {
      return songs[i];
  } else {
      return null;
    }
  }

  function play(elem) {
    var playing = document.getElementsByClassName("playing");
    if (playing.length > 0) {
      playing[0].setAttribute("class", "song");
    }

    var path = elem.getAttribute("path");
    var player = document.getElementById("player");

    var th = document.getElementById("title-header");
    th.innerHTML = '<h1 class="header">Now Playing:</h1>\n' +
                      '<div class="title">' + elem.innerHTML; + "</div>";

    elem.setAttribute("class", "song playing");
    player.setAttribute("src", path);
    player.play();
  }

  function playRandomSong(songs) {
    if (!!songs && songs.length > 0) {
      audio.i = Math.floor(Math.random() * songs.length);
      play(songs[audio.i]);
    }
  }

  function shouldLoop() {
    return document.getElementById("loop").checked;
  }

  function shouldShuffle() {
    return document.getElementById("shuf").checked;
  }

  return {
    i : 0,  // index of playing song

    onclick : function(song) {
      var songs = document.getElementById("songs").childNodes;
      // find song's place in the playlist
      for (audio.i = 0; audio.i < songs.length; audio.i++) {
        if (songs[audio.i] === song) {
          break;
        }
      }
      // set song to #selected?
      play(song);
    },

    onended : function(player) {
      audio.next();
    },

    load : function(song) {
      // TODO make sure it's the right type of object.
      var player = document.getElementById("player");
      player.setAttribute("src", song.getAttribute("path"));
    },

    next : function() {
      var songs = document.getElementById("songs").childNodes;
      if (songs.length > 0) {
        if (shouldShuffle()) {
          playRandomSong(songs);
        } else {
          audio.i += 1;
          if (songs.length > audio.i) {
            play(songs[audio.i]);
          } else if (shouldLoop()) {
            audio.i = 0;
            play(songs[0]);
          } else {
            audio.i -= 1;
          }
        }
      }
    },

    prev : function() {
      var songs = document.getElementById("songs").childNodes;
      if (shouldShuffle()) {
        playRandomSong(songs);
      } else {
        if (audio.i > 0) {
          audio.i -= 1;
        } else if (shouldLoop()) {
          audio.i = songs.length - 1;  // loop 'round
        }
        play(songs[audio.i]);
      }
    },

    onprogress : function(e) {
                   //inspect(e);
    },

    onplay : function(player) {
      var selected = document.getElementsByClassName("playing");
      if (selected.length === 0) {
        play(document.getElementById("songs").firstChild);
      }
    }
  };
})();
