var audio = (function() {
  // Adds a string to the notification thing---meant for notifications
  // to the user.  To print debug statements, use window.console.log().
  function notify(str) {
    document.getElementById("notification").innerHTML = str;
  }

  function inspect(obj) {
    var out = '';
    for (var p in obj) {
      out += p + ': ' + obj[p] + '\n';
    }
    window.console.log(out);
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

    var title = document.getElementById("page-title");
    var th = document.getElementById("title-header");
    title.innerHTML = elem.innerHTML;
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

    init : function(e) {
      var player = document.getElementById("player");
      player.onplay = function(e) { audio.onplay(player); };
      player.onended = audio.onended;
      player.ontimeupdate = audio.onprogress;

      document.getElementById("prev").onclick = audio.prev;
      document.getElementById("next").onclick = audio.next;
      document.getElementById("loop_label").onclick = function(e) {
        document.getElementById("loop").click();
      };
      document.getElementById("shuf_label").onclick = function(e) {
        document.getElementById("shuf").click();
      };
    },

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

    onended : function(e) {
      audio.next();
    },

    load : function(song) {
      // TODO make sure it's the right type of object.
      var player = document.getElementById("player");
      player.setAttribute("src", song.getAttribute("path"));
    },

    next : function(e) {
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

    prev : function(e) {
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

window.addEventListener("load", audio.init, true);
