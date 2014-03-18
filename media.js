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

  function play(elem) {
    var playing = $(".playing");
    if (playing.length > 0) {
      playing.first().attr("class", "media");
    }

    elem = $(elem);
    var path = elem.attr("path");
    var player = $("#player");

    $("#page-title").html(elem.html());
    $("#title-header").html('<h1 class="header">Now Playing:</h1>\n' +
        '<div class="title">' + elem.html() + "</div>");

    elem.attr("class", "media playing");
    player.attr("src", path);
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
      window.console.log("Deleting selected ID");
      $("#selected").removeAttr('id');
      $(med).attr('id', 'selected');
      play(med);
    },

    onended : function(e) {
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
