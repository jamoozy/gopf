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

var media = (function() {
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

  function get_ith_media(i) {
    var meds = document.getElementById("media").childNodes;
    if (meds.length > i) {
      return meds[i];
    } else {
      return null;
    }
  }

  function play(elem) {
    var playing = document.getElementsByClassName("playing");
    if (playing.length > 0) {
      playing[0].setAttribute("class", "media");
    }

    var path = elem.getAttribute("path");
    var player = document.getElementById("player");

    var title = document.getElementById("page-title");
    var th = document.getElementById("title-header");
    title.innerHTML = elem.innerHTML;
    th.innerHTML = '<h1 class="header">now playing:</h1>\n' +
                   '<div class="title">' + elem.innerHTML + "</div>" +
                   '<a download="" class="download" href="' + path + '">(download)</a>';

    elem.setAttribute("class", "media playing");
    player.setAttribute("src", path);
    player.play();
  }

  function playRandomSong(meds) {
    if (!!meds && meds.length > 0) {
      var j = media.i;
      while ((j = Math.floor(Math.random() * meds.length)) == media.i);
      media.i = j;
      play(meds[media.i]);
    }
  }

  function shouldLoop() {
    return document.getElementById("loop").checked;
  }

  function shouldShuffle() {
    return document.getElementById("shuf").checked;
  }

  return {
    i : 0,  // index of playing media

    init : function(e) {
      var player = document.getElementById("player");
      player.addEventListener("play", function(e) {
        media.onplay(player);
      }, true);
      player.addEventListener("ended", media.onended, true);
      player.addEventListener("timeupdate", media.onprogress, true);

      document.getElementById("prev").onclick = media.prev;
      document.getElementById("next").onclick = media.next;
      document.getElementById("loop_label").onclick = function(e) {
        document.getElementById("loop").click();
      };
      document.getElementById("shuf_label").onclick = function(e) {
        document.getElementById("shuf").click();
      };
    },

    onclick : function(med) {
      var meds = document.getElementById("media").childNodes;
      // Find media's place in the playlist
      for (media.i = 0; media.i < meds.length; media.i++) {
        if (meds[media.i] === med) {
          break;
        }
      }
      // Set media to #selected?
      window.console.log("Deleting selected ID");
      document.getElementById("selected").removeAttribute('id');
      med.setAttribute('id', 'selected');
      play(med);
    },

    onended : function(e) {
      media.next(e);
    },

    load : function(med) {
      // TODO make sure it's the right type of object.
      var player = document.getElementById("player");
      player.setAttribute("src", med.getAttribute("path"));
    },

    next : function(e) {
      var meds = document.getElementById("media").childNodes;
      if (meds.length > 0) {
        if (shouldShuffle()) {
          playRandomSong(meds);
        } else {
          media.i += 1;
          if (meds.length > media.i) {
            play(meds[media.i]);
          } else if (shouldLoop()) {
            media.i = 0;
            play(meds[0]);
          } else {
            media.i -= 1;
          }
        }
      }
    },

    prev : function(e) {
      var meds = document.getElementById("media").childNodes;
      if (shouldShuffle()) {
        playRandomSong(meds);
      } else {
        if (media.i > 0) {
          media.i -= 1;
        } else if (shouldLoop()) {
          media.i = meds.length - 1;  // loop 'round
        }
        play(meds[media.i]);
      }
    },

    onprogress : function(e) {
                   //inspect(e);
    },

    onplay : function(player) {
      var selected = document.getElementsByClassName("playing");
      if (selected.length === 0) {
        play(document.getElementById("media").firstChild);
      }
    }
  };
})();

window.addEventListener("load", media.init, true);
