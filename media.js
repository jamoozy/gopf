// Copyright (C) 2011-2015 Andrew "Jamoozy" Sabisch
//
// This file is part of GOPF.
//
// GOPF is free software: you can redistribute it and/or modify it under the
// terms of the GNU Affero General Public as published by the Free Software
// Foundation, either version 3 of the License, or (at your option) any later
// version.
//
// GOPF is distributed in the hope that it will be useful, but WITHOUT ANY
// WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
// FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more
// details.
//
// You should have received a copy of the GNU Affero General Public License
// along with GOPF. If not, see http://www.gnu.org/licenses/.

var media = (function() {
  // Adds a string to the notification thing---meant for notifications
  // to the user.  To print debug statements, use window.console.log().
  function notify(str) {
    $("#notification").html(str);
  }

  function inspect(obj) {
    var out = '';
    for (var p in obj) {
      out += p + ': ' + obj[p] + '\n';
    }
    window.console.log(out);
  }

  function get_ith_media(i) {
    var meds = $("#media").children();
    if (meds.length > i) {
      return meds[i];
    } else {
      return null;
    }
  }

  function play(m) {
    var playing = $(".playing");
    if (playing.length > 0) {
      playing.first().attr("class", "media");
    }

    var media = $(m),
        path = media.attr("path"),
        player = $("#player");

    $("#page-title").html(media.html());
    $("#title-header").html('<h1 class="header">Now Playing:</h1>\n' +
        '<div class="title"><span id="url-link" class="url-link" href="#">☃</span> ' + media.html() + "</div>");
    $("#url-link").click(function(e) {
      window.console.log("toggling " + $("#url"));
      $("#url").html(document.location.origin + document.location.pathname +
        "?p=" + encodeURIComponent($(".selected").html()) +
        "&m=" + encodeURIComponent($(".playing").html()) +
        "&t=" + Math.round($("#player")[0].currentTime));
      $("#url").toggle();
    });

    media.attr("class", "media playing");
    player.prop("src", path);
    if (!shouldShuffle()) {
      var i = nextID();
      if (i >= 0) {
        var path = $($("#media").children()[i]).attr("path");
        window.console.log("Trying path: " + path);
        window.console.log("... which is media #i: " + i);
        $("#preload").prop("href", path);
      } else {
        window.console.log("i was " + i);
      }
    }
    player[0].play();
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
    return $("#loop").prop('checked');
  }

  function shouldShuffle() {
    return $("#shuf").prop('checked');
  }

  function nextID() {
    var meds = $("#media").children(),
        i = media.i + 1;
    window.console.log("i: " + media.i + "  meds.length: " + meds.length);
    if (meds.length > i) {
      return i;
    } else if (shouldLoop()) {
      return 0;
    }
    return -1;
  }

  return {
    i : 0,  // index of playing media

    init : function(e) {
      var player = $("#player");
      player.on("play", function(e) {
        media.onplay(player);
      });
      player.on("ended", media.onended);
      player.on("timeupdate", media.onprogress);

      $("#prev").click(media.prev);
      $("#next").click(media.next);
      $("#loop_label").click(function(e) {
        $("#loop").trigger("click");
      });
      $("#shuf_label").click(function(e) {
        $("#shuf").trigger("click");
      });

      // On any kind of player error (during playback?), just go to the next
      // song.
      player.on("error", media.onerror);

      // Check for a get request that requests we play something right away.
      var playing = $('.playing');
      if (playing.size() > 0) {
        media.i = $(".media").index(playing) - 1;
        play(playing[0]);

        // TODO something with parsing t=\d+ from URL and setting seconds
        var match = /\Wt=(\d+)/.exec(document.location.href);
        if (match) {
          var quickplay = function() {
            this.currentTime = parseInt(match[1]);
            $(this).off('playing', null, quickplay);
            this.play();
          }
          $("#player").on('playing', quickplay);
        }
      }
    },

    onclick : function(med) {
      var meds = $("#media").children();

      // Find media's place in the playlist
      for (media.i = 0; media.i < meds.size(); media.i++) {
        if (meds[media.i] === med) {
          break;
        }
      }

      // Set media to #selected?
      $("#selected").removeAttr('id');
      $(med).attr('id', 'selected');
      play(med);
    },

    onended : function(e) {
      window.console.log("Media ended.  Calling next.");
      media.next(e);
    },

    onerror : function(e) {
      window.console.log("Warning!  Got a playback error:");
      window.console.log(e);
      media.next(e);
    },

    load : function(med) {
      // TODO make sure it's the right type of object.
      var player = $("#player");
      player.setAttribute("src", med.getAttribute("path"));
    },

    next : function(e) {
      var meds = $("#media").children();
      if (meds.length > 0) {
        if (shouldShuffle()) {
          playRandomSong(meds);
        } else {
          var i = nextID();
          if (i != -1) {
            media.i = i
            play(meds[i])
          }
        }
      }
    },

    prev : function(e) {
      var meds = $("#media").children();
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
      var selected = $(".playing");
      if (selected.length === 0) {
        play($("#media").children().first());
      }
    }
  };
})();

$(window).load(media.init);
